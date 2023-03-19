package ubuntu

import cl "github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib"

type Ubuntu2004CommandLib struct{}

// Common commands for control-plane and workers

func (u *Ubuntu2004CommandLib) SudoUpdate() cl.Command {
	return cl.Command("sudo apt update")
}

func (u *Ubuntu2004CommandLib) SudoFullUpgrade() cl.Command {
	return cl.Command("sudo apt -y full-upgrade")
}

func (u *Ubuntu2004CommandLib) AddCRIORepos() cl.Command {
	return cl.Command("echo \"deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/$OS/ /\"|sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list\necho \"deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable:/cri-o:/$CRIO_VERSION/$OS/ /\"|sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable:cri-o:$CRIO_VERSION.list\n")
}

func (u *Ubuntu2004CommandLib) ImportGPGKey() cl.Command {
	return cl.Command("curl -L https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable:cri-o:$CRIO_VERSION/$OS/Release.key | sudo apt-key add -\ncurl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/$OS/Release.key | sudo apt-key add -\n")
}

func (u *Ubuntu2004CommandLib) InstallCRIO() cl.Command {
	return cl.Command("sudo apt update\nsudo apt install -y cri-o cri-o-runc")
}

func (u *Ubuntu2004CommandLib) StartCRIO() cl.Command {
	return cl.Command("sudo systemctl enable crio.service\nsudo systemctl start crio.service")
}

func (u *Ubuntu2004CommandLib) DisableSWAP() cl.Command {
	return cl.Command("sudo swapoff -a")
}

func (u *Ubuntu2004CommandLib) InstallUtils() cl.Command {
	return cl.Command("sudo apt-get install -y apt-transport-https ca-certificates curl")
}

func (u *Ubuntu2004CommandLib) DownloadGoogleCloudSigningKey() cl.Command {
	return cl.Command("sudo curl -fsSLo /etc/apt/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg")
}

func (u *Ubuntu2004CommandLib) AddK8SRepo() cl.Command {
	return cl.Command("echo \"deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main\" | sudo tee /etc/apt/sources.list.d/kubernetes.list\n")
}

func (u *Ubuntu2004CommandLib) InstallKubeadm() cl.Command {
	return cl.Command("sudo apt-get update\nsudo apt-get install -y kubelet kubeadm kubectl\nsudo apt-mark hold kubelet kubeadm kubectl")
}

// Control-plane

func (u *Ubuntu2004CommandLib) InitKubeadm() cl.Command {
	return cl.Command("kubeadm init --podCIDR=10.244.0.0/16")
}

func (u *Ubuntu2004CommandLib) AddKubeConfig() cl.Command {
	return cl.Command("mkdir -p $HOME/.kube\nsudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config\nsudo chown $(id -u):$(id -g) $HOME/.kube/config")
}

func (u *Ubuntu2004CommandLib) AddFlannel() cl.Command {
	return cl.Command("kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml")
}
