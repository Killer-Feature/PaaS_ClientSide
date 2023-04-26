package service

import (
	"context"
	"fmt"
	"net/netip"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"github.com/Killer-Feature/PaaS_ClientSide/internal/models"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/executor"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/helm"
	k8s_installer "github.com/Killer-Feature/PaaS_ClientSide/pkg/k8s-installer"

	"github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib/ubuntu"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/socketmanager"
	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"
	"github.com/gorilla/websocket"

	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Service struct {
	r          internal.Repository
	l          *zap.Logger
	tm         *taskmanager.Manager[netip.AddrPort]
	sm         *socketmanager.SocketManager[internal.Message]
	hi         *helm.HelmInstaller
	progressCh chan internal.Message

	k8sInstaller *k8s_installer.Installer
}

// NewService returns instance of Huginn service
// Receives repository, logger and taskmanager structs as pointer
func NewService(r internal.Repository, l *zap.Logger, tm *taskmanager.Manager[netip.AddrPort], k8sInstaller *k8s_installer.Installer, hi *helm.HelmInstaller) internal.Usecase {
	progressCh := make(chan internal.Message)
	return &Service{
		r:            r,
		l:            l,
		tm:           tm,
		sm:           socketmanager.NewSocketManager[internal.Message](progressCh, nil, l),
		k8sInstaller: k8sInstaller,
		hi:           hi,
		progressCh:   progressCh,
	}
}

func (s *Service) ExecCommand(command string) ([]byte, error) {
	return executor.Exec(command)
}

func (s *Service) GetClusterNodes(ctx context.Context) ([]internal.Node, error) {
	nodes, err := s.r.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	respNodes := make([]internal.Node, len(nodes))
	var masterIP netip.AddrPort
	for _, node := range nodes {
		if node.IsMaster {
			masterIP = node.IP
		}
	}
	for i, node := range nodes {
		respNodes[i] = internal.Node{
			ID:        node.ID,
			IP:        node.IP,
			GrafanaIP: fmt.Sprintf("http://%s:3000/d/nMnqQpEVk/kubernetes-cluster-monitoring-via-prometheus?orgId=1&refresh=10s", masterIP.Addr().String()),
			Name:      node.Name,
			ClusterID: node.ClusterID,
			IsMaster:  node.IsMaster,
		}
	}

	return respNodes, nil
}

func (s *Service) AddNodeToCurrentCluster(ctx context.Context, id int) (int, error) {
	node, err := s.r.GetFullNode(ctx, id)
	if err != nil {
		return 0, err
	}

	taskID, err := s.tm.AddTask(s.addNodeToCurrentClusterProgressTask(context.Background(), node), node.IP)

	if err == nil {
		s.progressCh <- internal.Message{Type: internal.AddNodeToClusterT, Payload: internal.AddNodeToClusterProgressMsg{NodeID: node.ID, Status: internal.STATUS_IN_QUEUE, Percent: 0}}
	}

	return int(taskID), err
}

func (s *Service) addNodeToCurrentClusterProgressTask(ctx context.Context, node internal.FullNode) func(taskId taskmanager.ID) error {
	return func(taskID taskmanager.ID) error {
		sendProgress := func(percent int, status internal.TaskStatus, log string, err string) {
			s.progressCh <- internal.Message{Type: internal.AddNodeToClusterT, Payload: internal.AddNodeToClusterProgressMsg{NodeID: node.ID, Status: status, Percent: percent, Log: log, Error: err}}
		}
		sendProgress(1, internal.STATUS_START, "", "")
		sshBuilder := ssh.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(node.IP, node.Login, node.Password)
		if err != nil {
			sendProgress(1, internal.STATUS_ERROR, "", err.Error())
			return err
		}
		sendProgress(1, internal.STATUS_IN_PROCESS, "", "")
		defer func(cc cconn.ClientConn) {
			_ = cc.Close()
		}(cc)

		err = s.k8sInstaller.InstallK8S(cc, node.ID, node.IP.Addr().String(), sendProgress)

		return err
	}
}

func (s *Service) AddNode(ctx context.Context, node internal.FullNode) (int, error) {
	exists, err := s.r.IsNodeExists(ctx, node.IP.Addr())
	if err != nil {
		return 0, err
	}
	if exists == 0 {
		return s.r.AddNode(ctx, node)
	}
	return 0, internal.ErrNodeExists
}

func (s *Service) RemoveNode(ctx context.Context, id int) error {
	return s.r.RemoveNode(ctx, id)
}

func (s *Service) AddResource(ctx context.Context, rType internal.ResourceType, name string) error {
	err := s.hi.Install(name, rType)
	if err != nil {
		return err
	}
	return s.r.AddResource(ctx, convertResourceTypeToString(rType), name)
}

func convertResourceTypeToString(rtype internal.ResourceType) string {
	switch rtype {
	case internal.Postgres:
		return "postgres"
	case internal.Redis:
		return "redis"
	case internal.Prometheus:
		return "kube-prometheus"
	case internal.Grafana:
		return "grafana"
	default:
	}
	return "unknown"
}

func (s *Service) RemoveResource(ctx context.Context, rType internal.ResourceType, name string) error {
	return s.hi.UninstallChart(name)
}

func (s *Service) GetAdminConfig(ctx context.Context, clusterId int) (*models.AdminConfig, error) {
	_, masterIpStr, _, err := s.r.GetClusterTokenIPAndHash(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	ipport, err := netip.ParseAddrPort(masterIpStr)
	if err != nil {
		return nil, err
	}

	masterId, err := s.r.IsNodeExists(ctx, ipport.Addr())
	if err != nil {
		return nil, err
	}

	node, err := s.r.GetFullNode(ctx, masterId)
	if err != nil {
		return nil, err
	}

	sshBuilder := ssh.NewSSHBuilder()
	cc, err := sshBuilder.CreateCC(node.IP, node.Login, node.Password)

	if err != nil {
		configFile, err := os.ReadFile("./config")
		if err != nil {
			return nil, err
		}
		return &models.AdminConfig{Config: string(configFile)}, nil
	}

	defer func(cc cconn.ClientConn) {
		_ = cc.Close()
	}(cc)

	output, err := s.getAdminConf(ctx, cc)

	if err != nil {
		configFile, err := os.ReadFile("./config")
		if err != nil {
			return nil, err
		}
		return &models.AdminConfig{Config: string(configFile)}, nil
	}

	_ = os.WriteFile("./config", output, 666)
	return &models.AdminConfig{Config: string(output)}, nil
}

func (s *Service) getAdminConf(ctx context.Context, cc cconn.ClientConn) ([]byte, error) {
	cl := ubuntu.Ubuntu2004CommandLib{}
	getAdminConfCommand := cl.CatAdminConfFile()
	output, err := cc.Exec(string(getAdminConfCommand.Command))
	if err != nil {
		s.l.Error("error getting admin.conf", zap.String("error", err.Error()))
		return nil, err
	}
	return output, nil
}

func (s *Service) GetResources(ctx context.Context) ([]internal.Resource, error) {
	resources, err := s.hi.GetResourcesList()

	if err != nil {
		return nil, err
	}

	resourceList := make([]internal.Resource, 0, len(resources))

	for _, res := range resources {
		resourceList = append(resourceList, internal.Resource{
			Name:          res.Name,
			Status:        res.Status,
			FirstDeployed: res.FirstDeployed,
			LastDeployed:  res.LastDeployed,
			AppVersion:    res.AppVersion,
			Description:   res.Description,
			ChartVersion:  res.ChartVersion,
			ApiVersion:    res.ApiVersion,
			Type:          res.Type,
			ChartURL:      res.ChartURL,
		})
	}
	return resourceList, nil
}

func (s *Service) GetServices(ctx context.Context) ([]internal.Service, error) {
	kubeconfig := "./config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	services, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	serviceList := make([]internal.Service, 0, len(services.Items))

	for _, res := range services.Items {
		serviceList = append(serviceList, internal.Service{
			Name:      res.Name,
			Namespace: res.ObjectMeta.Namespace,
			Type:      string(res.Spec.Type),
			Created:   res.ObjectMeta.CreationTimestamp.String(),
			Age:       time.Now().Sub(res.ObjectMeta.CreationTimestamp.Time).Round(time.Second).String(),
		})
	}
	return serviceList, nil
}

func (s *Service) RemoveNodeFromCurrentCluster(ctx context.Context, id int) (int, error) {
	node, err := s.r.GetFullNode(ctx, id)
	if err != nil {
		return 0, err
	}

	taskID, err := s.tm.AddTask(s.removeNodeFromCurrentClusterProgressTask(context.Background(), node), node.IP)
	if err == nil {
		s.progressCh <- internal.Message{Type: internal.RemoveNodeFromClusterT, Payload: internal.AddNodeToClusterProgressMsg{NodeID: node.ID, Status: internal.STATUS_IN_QUEUE, Percent: 0}}
	}
	return int(taskID), err
}

func (s *Service) removeNodeFromCurrentClusterProgressTask(ctx context.Context, node internal.FullNode) func(taskId taskmanager.ID) error {
	return func(taskID taskmanager.ID) error {
		sendProgress := func(percent int, status internal.TaskStatus, log string, err string) {
			s.progressCh <- internal.Message{Type: internal.RemoveNodeFromClusterT, Payload: internal.AddNodeToClusterProgressMsg{NodeID: node.ID, Status: status, Percent: percent, Log: log, Error: err}}
		}

		sendProgress(1, internal.STATUS_START, "", "")
		sshBuilder := ssh.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(node.IP, node.Login, node.Password)
		if err != nil {
			return err
		}
		defer func(cc cconn.ClientConn) {
			_ = cc.Close()
		}(cc)
		err = s.k8sInstaller.RemoveK8S(cc, sendProgress)
		if err != nil {
			return err
		}

		defer func(r internal.Repository, ctx context.Context, id int) {
			_ = r.ResetNodeCluster(ctx, id)
		}(s.r, ctx, node.ID)

		_, masterIpStr, _, err := s.r.GetClusterTokenIPAndHash(ctx, 1)
		if err != nil {
			return err
		}

		ipport, err := netip.ParseAddrPort(masterIpStr)
		if err != nil {
			return err
		}

		masterId, err := s.r.IsNodeExists(ctx, ipport.Addr())
		if err != nil {
			return err
		}

		if masterId == node.ID {
			s.l.Info("Removing cluster token and hash from DB")
			err = s.r.DeleteClusterTokenIPAndHash(ctx, masterId)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (s *Service) GetProgress(ctx context.Context, socket *websocket.Conn) error {
	s.sm.SetSocket(socket)

	go s.sm.SendResultToSocketByTicker(time.Second, func() interface{} {
		return interface{}(struct{ Message string }{"Message"})
	})
	return nil
}
