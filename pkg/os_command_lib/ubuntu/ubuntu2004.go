package ubuntu

import (
	"fmt"
	"net/netip"

	cl "github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib"
)

type Ubuntu2004CommandLib struct{}

// Common commands for control-plane and workers

func (u *Ubuntu2004CommandLib) SudoUpdate() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo apt update",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) SudoFullUpgrade() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo apt -y full-upgrade",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddCRIORepos() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "echo \"deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_20.04/ /\"|sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list\necho \"deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable:/cri-o:/1.26/xUbuntu_20.04/ /\"|sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable:cri-o:1.26.list",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) ImportGPGKey() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "curl -L https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable:cri-o:1.26/xUbuntu_20.04/Release.key | sudo apt-key add -; curl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_20.04/Release.key | sudo apt-key add -",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) InstallCRIO() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo apt install -y cri-o cri-o-runc",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) StartCRIO() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo systemctl enable crio.service\nsudo systemctl start crio.service",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) DisableSWAP() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo swapoff -a",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) InstallUtils() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo apt-get install -y apt-transport-https ca-certificates curl",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) DownloadGoogleCloudSigningKey() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo mkdir -p /etc/apt/keyrings\nsudo touch /etc/apt/keyrings/kubernetes-archive-keyring.gpg\necho \"xsBNBGKItdQBCADWmKTNZEYWgXy73FvKFY5fRro4tGNa4Be4TZW3wZpct9Cj8EjykU7S9EPoJ3EdKpxFltHRu7QbDi6LWSNA4XxwnudQrYGxnxx6Ru1KBHFxHhLfWsvFcGMwit/znpxtIt9UzqCm2YTEW5NUnzQ4rXYqVQK2FLG4weYJ5bKwkY+ZsnRJpzxdHGJ0pBiqwkMT8bfQdJymUBown+SeuQ2HEqfjVMsIRe0dweD2PHWeWo9fTXsz1Q5abiGckyOVyoN9//DgSvLUocUcZsrWvYPaN+o8lXTO3GYFGNVsx069rxarkeCjOpiQOWrQmywXISQudcusSgmmgfsRZYW7FDBy5MQrABEBAAHNUVJhcHR1cmUgQXV0b21hdGljIFNpZ25pbmcgS2V5IChjbG91ZC1yYXB0dXJlLXNpZ25pbmcta2V5LTIwMjItMDMtMDctMDhfMDFfMDEucHViKcLAYgQTAQgAFgUCYoi11AkQtT3IDRPt7wUCGwMCGQEAAMGoCAB8QBNIIN3Q2D3aahrfkb6axd55zOwR0tnriuJRoPHoNuorOpCv9aWMMvQACNWkxsvJxEF8OUbzhSYjAR534RDigjTetjK2i2wKLz/kJjZbuF4ZXMynCm40eVm1XZqU63U9XR2RxmXppyNpMqQO9LrzGEnNJuh23icaZY6no12axymxcle/+SCmda8oDAfa0iyA2iyg/eU05buZv54MC6RB13QtS+8vOrKDGr7RYp/VYvQzYWm+ck6DvlaVX6VB51BkLl23SQknyZIJBVPm8ttU65EyrrgG1jLLHFXDUqJ/RpNKq+PCzWiyt4uy3AfXK89RczLu3uxiD0CQI0T31u/IzsBNBGKItdQBCADIMMJdRcg0Phv7+CrZz3xRE8Fbz8AN+YCLigQeH0B9lijxkjAFr+thB0IrOu7ruwNY+mvdP6dAewUur+pJaIjEe+4s8JBEFb4BxJfBBPuEbGSxbi4OPEJuwT53TMJMEs7+gIxCCmwioTggTBp6JzDsT/cdBeyWCusCQwDWpqoYCoUWJLrUQ6dOlI7s6p+iIUNIamtyBCwb4izs27HdEpX8gvO9rEdtcb7399HyO3oD4gHgcuFiuZTpvWHdn9WYwPGM6npJNG7crtLnctTR0cP9KutSPNzpySeAniHx8L9ebdD9tNPCWC+OtOcGRrcBeEznkYh1C4kzdP1ORm5upnknABEBAAHCwF8EGAEIABMFAmKItdQJELU9yA0T7e8FAhsMAABJmAgAhRPk/dFj71bU/UTXrkEkZZzE9JzUgan/ttyRrV6QbFZABByf4pYjBj+yLKw3280//JWurKox2uzEq1hdXPedRHICRuh1Fjd00otaQ+wGF3kY74zlWivB6Wp6tnL9STQ1oVYBUv7HhSHoJ5shELyedxxHxurUgFAD+pbFXIiK8cnAHfXTJMcrmPpC+YWEC/DeqIyEcNPkzRhtRSuERXcq1n+KJvMUAKMD/tezwvujzBaaSWapmdnGmtRjjL7IxUeGamVWOwLQbUr+34MwzdeJdcL8fav5LA8Uk0ulyeXdwiAK8FKQsixI+xZvz7HUs8ln4pZwGw/TpvO9cMkHogtgzQ==\" | base64 -d | sudo tee -a /etc/apt/keyrings/kubernetes-archive-keyring.gpg",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddK8SRepo() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "echo \"deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main\" | sudo tee /etc/apt/sources.list.d/kubernetes.list",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) InstallKubeadm() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo apt-get update\nsudo apt-get install -y kubelet kubeadm kubectl\nsudo apt-mark hold kubelet kubeadm kubectl",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) SetModprobe() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo modprobe br_netfilter",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) SudoSu() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo su",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) SetIpForward() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "echo '1' | sudo tee -a /proc/sys/net/ipv4/ip_forward",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) Exit() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "exit",
		Parser:    nil,
		Condition: cl.Required,
	}
}

// Control-plane

func (u *Ubuntu2004CommandLib) InitKubeadm(parser cl.Parser) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo kubeadm init --pod-network-cidr=10.244.0.0/16",
		Parser:    parser,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) UntaintControlPlane() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "kubectl taint nodes --all node-role.kubernetes.io/control-plane-",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddKubeConfig() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "rm -r  $HOME/.kube \n mkdir -p $HOME/.kube\nsudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config\nsudo chown $(id -u):$(id -g) $HOME/.kube/config",
		Parser:    nil,
		Condition: cl.Anyway,
	}
}

func (u *Ubuntu2004CommandLib) AddFlannel() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) InstallHelm() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo snap install helm --classic",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddBitnamiRepo() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "helm repo add bitnami https://charts.bitnami.com/bitnami",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) InstallPrometheus() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "helm install prometheus bitnami/kube-prometheus",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) CreateFolderForPV() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "sudo mkdir -p /devkube/postgres",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddStorageClass() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   "kubectl apply -f - <<EOF kind: StorageClass\napiVersion: storage.k8s.io/v1\nmetadata:\n  name: local-storage\n  annotations:\n    storageclass.kubernetes.io/is-default-class: \"true\"\nprovisioner: kubernetes.io/no-provisioner\nvolumeBindingMode: Immediate\nEOF",
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddMetallbConf(ip string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("kubectl apply -f - <<EOF ---\napiVersion: metallb.io/v1beta1\nkind: IPAddressPool\nmetadata:\n  name: default\n  namespace: default\nspec:\n  addresses:\n  - %s/32\n  autoAssign: true\n---\napiVersion: metallb.io/v1beta1\nkind: L2Advertisement\nmetadata:\n  name: default\n  namespace: default\nspec:\n  ipAddressPools:\n  - default\nEOF", ip)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddPostgresPV(hostname string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("kubectl apply -f - <<EOF apiVersion: v1\nkind: PersistentVolume\nmetadata:\n  name: pv-7\n  labels:\n    type: local\nspec:\n  capacity:\n    storage: 4Gi\n  volumeMode: Filesystem\n  accessModes:\n  - ReadWriteOnce\n  persistentVolumeReclaimPolicy: Retain\n  storageClassName: local-storage\n  local:\n    path: /devkube/postgresql\n  nodeAffinity:\n    required:\n      nodeSelectorTerms:\n      - matchExpressions:\n        - key: kubernetes.io/hostname\n          operator: In\n          values:\n          - %s\nEOF", hostname)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddGrafanaPV(hostname string) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("kubectl apply -f - <<EOF apiVersion: v1\nkind: PersistentVolume\nmetadata:\n  name: pv-grafana\n  labels:\n    type: local\nspec:\n  capacity:\n    storage: 10Gi\n  volumeMode: Filesystem\n  accessModes:\n  - ReadWriteOnce\n  persistentVolumeReclaimPolicy: Retain\n  storageClassName: local-storage\n  local:\n    path: /devkube/postgresql\n  nodeAffinity:\n    required:\n      nodeSelectorTerms:\n      - matchExpressions:\n        - key: kubernetes.io/hostname\n          operator: In\n          values:\n          - %s\nEOF", hostname)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

func (u *Ubuntu2004CommandLib) AddGrafanaIngress(userNum int) cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command(fmt.Sprintf("kubectl apply -f - <<EOF apiVersion: networking.k8s.io/v1\nkind: Ingress\nmetadata:\n  name: grafana-ingress\nspec:\n  ingressClassName: nginx\n  rules:\n  - host: \"grafana.user%d.huginn.pro\"\n    http:\n      paths:\n      - path: /\n        pathType: Prefix\n        backend:\n          service:\n            name: grafana\n            port:\n              number: 3000\nEOF", userNum)),
		Parser:    nil,
		Condition: cl.Required,
	}
}

// Join cluster

func (u *Ubuntu2004CommandLib) KubeadmJoin(ip netip.AddrPort, token, tokenHash string) cl.CommandAndParser {

	cp := cl.CommandAndParser{
		Command:   "sudo kubeadm join",
		Parser:    nil,
		Condition: 0,
	}
	cp = cp.WithArgs(ip.String(), "--token", token, "--discovery-token-ca-cert-hash", tokenHash)
	return cp
}

// Reset Kubeadm

func (u *Ubuntu2004CommandLib) KubeadmReset() cl.CommandAndParser {
	cp := cl.CommandAndParser{
		Command:   "sudo kubeadm reset -f",
		Parser:    nil,
		Condition: cl.Anyway,
	}
	return cp
}

func (u *Ubuntu2004CommandLib) StopKubelet() cl.CommandAndParser {
	cp := cl.CommandAndParser{
		Command:   "sudo systemctl stop kubelet",
		Parser:    nil,
		Condition: cl.Anyway,
	}
	return cp
}

func (u *Ubuntu2004CommandLib) StopCRIO() cl.CommandAndParser {
	cp := cl.CommandAndParser{
		Command:   "sudo systemctl stop crio.service",
		Parser:    nil,
		Condition: cl.Anyway,
	}
	return cp
}

func (u *Ubuntu2004CommandLib) LinkDownCNI0() cl.CommandAndParser {
	cp := cl.CommandAndParser{
		Command:   "sudo ip link set cni0 down",
		Parser:    nil,
		Condition: cl.Anyway,
	}
	return cp
}

func (u *Ubuntu2004CommandLib) IpconfigCNI0Down() cl.CommandAndParser {
	cp := cl.CommandAndParser{
		Command:   "sudo ifconfig cni0 down",
		Parser:    nil,
		Condition: cl.Anyway,
	}
	return cp
}

func (u *Ubuntu2004CommandLib) IpconfigFlannelDown() cl.CommandAndParser {
	cp := cl.CommandAndParser{
		Command:   "sudo ifconfig flannel.1 down",
		Parser:    nil,
		Condition: cl.Anyway,
	}
	return cp
}

func (u *Ubuntu2004CommandLib) BrctlDelbr() cl.CommandAndParser {
	cp := cl.CommandAndParser{
		Command:   "sudo brctl delbr cni0",
		Parser:    nil,
		Condition: cl.Anyway,
	}
	return cp
}

func (u *Ubuntu2004CommandLib) CatAdminConfFile() cl.CommandAndParser {
	return cl.CommandAndParser{
		Command:   cl.Command("sudo cat /etc/kubernetes/admin.conf"),
		Parser:    nil,
		Condition: cl.Required,
	}
}
