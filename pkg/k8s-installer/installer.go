package k8s_installer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/helm"
	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib/ubuntu"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"go.uber.org/zap"
)

const (
	LOG_INITIAL_SIZE = 2048
)

var (
	grafanaArgs = map[string]string{
		// comma seperated values to set
		"set": "admin.password=admin",
	}
)

type Installer struct {
	r  internal.Repository
	l  *zap.Logger
	hi *helm.HelmInstaller
}

func NewInstaller(l *zap.Logger, r internal.Repository, hi *helm.HelmInstaller) *Installer {
	return &Installer{
		r:  r,
		l:  l,
		hi: hi,
	}
}

func (installer *Installer) installKubeadm() []cl.CommandAndParser {
	commandLib := ubuntu.Ubuntu2004CommandLib{}

	commands := []cl.CommandAndParser{
		commandLib.SudoUpdate(),
		commandLib.SudoFullUpgrade(),
		commandLib.AddCRIORepos(),
		commandLib.ImportGPGKey(),
		commandLib.SudoUpdate(),
		commandLib.InstallCRIO(),
		commandLib.StartCRIO(),
		commandLib.DisableSWAP(),
		commandLib.InstallUtils(),
		commandLib.DownloadGoogleCloudSigningKey(),
		commandLib.AddK8SRepo(),
		commandLib.InstallKubeadm(),
		commandLib.SetModprobe(),
		commandLib.SetIpForward(),
	}
	return commands
}

func (installer *Installer) kubeadmInit() []cl.CommandAndParser {
	commandLib := ubuntu.Ubuntu2004CommandLib{}

	commands := []cl.CommandAndParser{
		commandLib.InitKubeadm(installer.parseKubeadmInit),
		commandLib.AddKubeConfig(),
		commandLib.UntaintControlPlane(),
		commandLib.AddFlannel(),
		commandLib.InstallHelm(),
		commandLib.AddBitnamiRepo(),
		commandLib.InstallPrometheus(),
	}
	return commands
}

func (installer *Installer) kubeadmCreateGrafana() []cl.CommandAndParser {
	commandLib := ubuntu.Ubuntu2004CommandLib{}

	commands := []cl.CommandAndParser{
		commandLib.AddStorageClass(),
		commandLib.AddGrafanaPV(),
		commandLib.AddPostgresPV(),
		commandLib.AddGrafanaIngress(),
		commandLib.CreateFolderForPV(),
	}
	return commands
}

func (installer *Installer) kubeadmReset() []cl.CommandAndParser {
	commandLib := ubuntu.Ubuntu2004CommandLib{}

	commands := []cl.CommandAndParser{
		commandLib.KubeadmReset(),
		commandLib.StopKubelet(),
		commandLib.StopCRIO(),
		commandLib.LinkDownCNI0(),
		commandLib.IpconfigCNI0Down(),
		commandLib.IpconfigFlannelDown(),
		commandLib.BrctlDelbr(),
	}
	return commands
}

func (installer *Installer) kubeadmJoin(token, ip, hash string) []cl.CommandAndParser {
	commandLib := ubuntu.Ubuntu2004CommandLib{}

	return []cl.CommandAndParser{
		commandLib.KubeadmJoin(netip.MustParseAddrPort(ip), token, hash),
	}
}

var matchRe = regexp.MustCompile(`(?P<hostport>[a-z0-9-_:.]*) --token (?P<token>[a-z0-9-_.]*) \\\n\t--discovery-token-ca-cert-hash (?P<hash>[a-z0-9-:]*)`)

func (installer *Installer) parseKubeadmInit(output []byte, extraData interface{}) error {
	outputstr := string(output)
	outputstrs := strings.Split(outputstr, "kubeadm join ")
	matchMap := make(map[string]string, len(matchRe.SubexpNames()))
	match := matchRe.FindStringSubmatch(outputstrs[1])
	if len(match) == 0 {
		return fmt.Errorf("no match for regexp")
	}

	for i, group := range matchRe.SubexpNames() {
		matchMap[group] = match[i]
	}

	return installer.r.AddClusterTokenIPAndHash(context.Background(), 1, matchMap["token"], matchMap["hostport"], matchMap["hash"])
}

func (installer *Installer) InstallK8S(conn client_conn.ClientConn, nodeid int, sendProgress func(percent int, status internal.TaskStatus, log string, err string)) error {

	kubeadmInstallCommands := installer.installKubeadm()

	isClusterExists, err := installer.r.CheckClusterTokenIPAndHash(context.Background(), 1)
	if err != nil {
		return err
	}

	commandNumber := 20
	percent, k := 1, 1
	percentNext := func() int {
		percent = (k*100 - 1) / commandNumber
		k++
		return percent
	}

	if isClusterExists {
		token, ip, hash, err := installer.r.GetClusterTokenIPAndHash(context.Background(), 1)
		if err != nil {
			return err
		}
		installer.l.Info("Adding new worker to cluster")
		kubeadmInstallCommands = append(kubeadmInstallCommands, installer.kubeadmJoin(token, ip, hash)...)
	} else {
		installer.l.Info("Adding new control plane to cluster")
		kubeadmInstallCommands = append(kubeadmInstallCommands, installer.kubeadmInit()...)
		commandNumber = 38
	}

	log := make([]byte, 0, LOG_INITIAL_SIZE)

	for _, command := range kubeadmInstallCommands {
		exec, err := conn.Exec(string(command.Command))
		log = pushToLog(log, []byte(command.Command), exec)
		if err != nil && command.Condition != cl.Anyway {
			sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
			installer.l.Error("exec failed", zap.String("command", string(command.Command)))
			return err
		}

		if command.Parser != nil {
			err = command.Parser(exec, nil)
			if err != nil {
				sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
				return err
			}

		}
		sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
		installer.l.Info("installation percent", zap.Int("percent", percent), zap.String("command", string(command.Command)))
	}

	err = installer.r.SetNodeClusterID(context.Background(), nodeid, 1)
	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		return err
	}

	if isClusterExists {
		sendProgress(100, internal.STATUS_SUCCESS, string(log), "")
		return nil
	}

	config, err := installer.getAdminConf(context.Background(), conn)

	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		return err
	}

	sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")

	err = os.WriteFile("./config", config, 0664)
	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		installer.l.Error("error writing admin.conf to ./config", zap.String("error", err.Error()))
		return err
	}

	time.Sleep(30 * time.Second)
	sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")

	installer.hi.SetNewConfig()
	err = installer.hi.InstallChart("metallb", "metallb", "metallb", nil)
	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		return err
	}

	sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
	time.Sleep(30 * time.Second)

	commandLib := ubuntu.Ubuntu2004CommandLib{}
	command := commandLib.AddMetallbConf()

	exec, err := conn.Exec(string(command.Command))
	log = pushToLog(log, []byte(command.Command), exec)
	installer.l.Info("metallb installed")
	if err != nil && command.Condition != cl.Anyway {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		installer.l.Error("exec failed", zap.String("command", string(command.Command)), zap.String("res", string(exec)))
		return err
	}

	sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
	if command.Parser != nil {
		err = command.Parser(exec, nil)
		if err != nil {
			sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
			return err
		}
	}

	sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
	time.Sleep(5 * time.Second)

	err = installer.hi.InstallChart("nginx-ingress-controller", "bitnami", "nginx-ingress-controller", nil)
	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		return err
	}

	sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
	time.Sleep(30 * time.Second)

	grafanaCommands := installer.kubeadmCreateGrafana()
	for _, command := range grafanaCommands {
		exec, err := conn.Exec(string(command.Command))
		log = pushToLog(log, []byte(command.Command), exec)
		installer.l.Info("grafana installation percent", zap.Int("percent", percent/len(grafanaCommands)))

		if err != nil && command.Condition != cl.Anyway {
			installer.l.Error("exec failed", zap.String("command", string(command.Command)))
			sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
			return err
		}

		if command.Parser != nil {
			err = command.Parser(exec, nil)
			if err != nil {
				sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
				return err
			}
		}
		sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
	}

	time.Sleep(1 * time.Minute)

	sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
	err = installer.hi.InstallChart("grafana", "bitnami", "grafana", grafanaArgs)

	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		return err
	}

	time.Sleep(1 * time.Minute)
	exec, err = conn.Exec("kubectl exec --namespace default -it $(kubectl get pods --namespace default -lapp.kubernetes.io/name=grafana -o jsonpath=\"{.items[0].metadata.name}\") grafana-cli admin reset-admin-password admin")

	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		return err
	}

	err = installer.portForwarding("default", "grafana", "3000", "3000")
	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		installer.l.Error("exec failed", zap.String("command", string(command.Command)))
		return err
	}

	err = installer.portForwarding("default", "prometheus", "9090", "9090")
	if err != nil {
		sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
		installer.l.Error("exec failed", zap.String("command", string(command.Command)))
		return err
	}

	sendProgress(100, internal.STATUS_SUCCESS, string(log), "")

	return nil
}

func (installer *Installer) RemoveK8S(conn client_conn.ClientConn, sendProgress func(percent int, status internal.TaskStatus, log string, err string)) error {
	kubeadmStopCommands := installer.kubeadmReset()

	percent, k := 1, 1
	percentNext := func() int {
		percent = (k*100 - 1) / 8
		k++
		return percent
	}

	log := make([]byte, 0, LOG_INITIAL_SIZE)

	for i, command := range kubeadmStopCommands {
		exec, err := conn.Exec(string(command.Command))
		log = pushToLog(log, []byte(command.Command), exec)
		installer.l.Info("installation percent", zap.Int("percent", (i+1)*100/len(kubeadmStopCommands)))
		if err != nil && command.Condition != cl.Anyway {
			sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
			installer.l.Error("exec failed", zap.String("command", string(command.Command)))
			return err
		} else if err != nil {
			installer.l.Warn("exec failed", zap.String("command", string(command.Command)))
		}

		if command.Parser != nil {
			err = command.Parser(exec, nil)
			if err != nil {
				sendProgress(percentNext(), internal.STATUS_ERROR, string(log), err.Error())
				return err
			}
		}
		sendProgress(percentNext(), internal.STATUS_IN_PROCESS, string(log), "")
	}
	sendProgress(100, internal.STATUS_SUCCESS, string(log), "")
	return nil
}

func (installer *Installer) getAdminConf(ctx context.Context, cc client_conn.ClientConn) ([]byte, error) {
	cl := ubuntu.Ubuntu2004CommandLib{}
	getAdminConfCommand := cl.CatAdminConfFile()
	output, err := cc.Exec(string(getAdminConfCommand.Command))
	if err != nil {
		installer.l.Error("error getting admin.conf", zap.String("error", err.Error()))
		return nil, err
	}
	return output, nil
}

func (installer *Installer) portForwarding(namespace, appName, portLocal, portRemote string) error {
	config, err := clientcmd.BuildConfigFromFlags("", "./config")
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{LabelSelector: `app.kubernetes.io/name=` + appName})
	if len(pods.Items) != 1 {
		return errors.New("error no found or multiple choices pods by labelselector" + `app.kubernetes.io/name=` + appName)
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, pods.Items[0].ObjectMeta.Name)
	hostIP := strings.TrimLeft(config.Host, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}

	roundTripper, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	go func() {
		forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%s:%s", portLocal, portRemote)}, nil, nil, out, errOut)
		if err = forwarder.ForwardPorts(); err != nil { // Locks until stopChan is closed.
			installer.l.Error("error forwarding", zap.Error(err))
		}
	}()

	return nil
}

func pushToLog(log []byte, command []byte, output []byte) []byte {
	command = bytes.ReplaceAll(command, []byte("\n"), []byte("\n$ "))
	return bytes.Join([][]byte{log, append([]byte("$ "), command...), output}, []byte("\n"))
}
