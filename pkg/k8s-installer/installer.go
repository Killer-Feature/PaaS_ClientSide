package k8s_installer

import (
	"context"
	"fmt"
	"net/netip"
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

func (installer *Installer) InstallK8S(conn client_conn.ClientConn, nodeid int) error {
	kubeadmInstallCommands := installer.installKubeadm()

	isClusterExists, err := installer.r.CheckClusterTokenIPAndHash(context.Background(), 1)
	if err != nil {
		return err
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
	}

	for i, command := range kubeadmInstallCommands {
		exec, err := conn.Exec(string(command.Command))
		installer.l.Info("installation percent", zap.Int("percent", (i+1)*100/len(kubeadmInstallCommands)))
		if err != nil && command.Condition != cl.Anyway {
			installer.l.Error("exec failed", zap.String("command", string(command.Command)))
			return err
		}

		if command.Parser != nil {
			err = command.Parser(exec, nil)
			if err != nil {
				return err
			}
		}
	}

	err = installer.r.SetNodeClusterID(context.Background(), nodeid, 1)
	if err != nil {
		return err
	}

	if isClusterExists {
		return nil
	}

	// TODO: это надо делать вообще только при добавлении мастера, для этого надо перенести эту операцию в s.k8sInstaller.InstallK8S(cc)
	config, err := installer.getAdminConf(context.Background(), conn)
	if err == nil {
		err = os.WriteFile("./config", config, 0664)
		if err != nil {
			installer.l.Error("error writing admin.conf to ./config", zap.String("error", err.Error()))
			return err
		}
	}

	time.Sleep(30 * time.Second)

	err = installer.hi.InstallChart("metallb", "metallb", "metallb", nil)
	if err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	commandLib := ubuntu.Ubuntu2004CommandLib{}
	command := commandLib.AddMetallbConf()

	exec, err := conn.Exec(string(command.Command))
	installer.l.Info("metallb installed")
	if err != nil && command.Condition != cl.Anyway {
		installer.l.Error("exec failed", zap.String("command", string(command.Command)))
		return err
	}

	if command.Parser != nil {
		err = command.Parser(exec, nil)
		if err != nil {
			return err
		}
	}

	time.Sleep(5 * time.Second)

	err = installer.hi.InstallChart("nginx-ingress-controller", "bitnami", "nginx-ingress-controller", nil)
	if err != nil {
		return err
	}

	time.Sleep(30 * time.Second)

	grafanaCommands := installer.kubeadmCreateGrafana()
	for i, command := range grafanaCommands {
		exec, err := conn.Exec(string(command.Command))
		installer.l.Info("grafana installation percent", zap.Int("percent", (i+1)*100/len(grafanaCommands)))
		if err != nil && command.Condition != cl.Anyway {
			installer.l.Error("exec failed", zap.String("command", string(command.Command)))
			return err
		}

		if command.Parser != nil {
			err = command.Parser(exec, nil)
			if err != nil {
				return err
			}
		}
	}

	time.Sleep(1 * time.Minute)

	err = installer.hi.InstallChart("grafana", "bitnami", "grafana", nil)
	if err != nil {
		return err
	}

	return nil
}

func (installer *Installer) RemoveK8S(conn client_conn.ClientConn) error {
	kubeadmStopCommands := installer.kubeadmReset()

	for i, command := range kubeadmStopCommands {
		exec, err := conn.Exec(string(command.Command))
		installer.l.Info("installation percent", zap.Int("percent", (i+1)*100/len(kubeadmStopCommands)))
		if err != nil && command.Condition != cl.Anyway {
			installer.l.Error("exec failed", zap.String("command", string(command.Command)))
			return err
		} else if err != nil {
			installer.l.Warn("exec failed", zap.String("command", string(command.Command)))
		}

		if command.Parser != nil {
			err = command.Parser(exec, nil)
			if err != nil {
				return err
			}
		}
	}
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
