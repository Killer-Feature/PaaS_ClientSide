package k8s_installer

import (
	"context"
	"fmt"
	"net/netip"
	"regexp"
	"strings"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib/ubuntu"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"go.uber.org/zap"
)

type Installer struct {
	r internal.Repository
	l *zap.Logger
}

func NewInstaller(l *zap.Logger, r internal.Repository) *Installer {
	return &Installer{
		r: r,
		l: l,
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
		commandLib.UntaintControlPlane(),
		commandLib.AddKubeConfig(),
		commandLib.AddFlannel(),
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

func (installer *Installer) InstallK8S(conn client_conn.ClientConn) error {
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
