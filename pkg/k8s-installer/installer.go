package k8s_installer

import (
	"fmt"
	"net/netip"

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
		commandLib.AddKubeConfig(),
		commandLib.AddFlannel(),
	}
	return commands
}

func (installer *Installer) kubeadmJoin() []cl.CommandAndParser {
	commandLib := ubuntu.Ubuntu2004CommandLib{}

	return []cl.CommandAndParser{
		commandLib.KubeadmJoin(netip.MustParseAddrPort("89.208.220.55:6443"), "7pnvri.q81vrsmbcblvqa3e", "sha256:3d31d94f1905bce8867a23026007308c9ba334d41eab740533796d3b423c901a"),
	}
}

func (installer *Installer) parseKubeadmInit(output []byte, extraData interface{}) error {
	// TODO: parse output
	fmt.Println(string(output))
	return nil
}

func (installer *Installer) InstallK8S(conn client_conn.ClientConn) error {
	kubeadmInstallCommands := installer.installKubeadm()
	kubeadmInstallCommands = append(kubeadmInstallCommands, installer.kubeadmJoin()...)
	for _, command := range kubeadmInstallCommands {
		exec, err := conn.Exec(string(command.Command))
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
