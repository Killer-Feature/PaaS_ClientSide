package service

import (
	"context"
	"net/netip"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"github.com/Killer-Feature/PaaS_ClientSide/internal/models"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/executor"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/helm"
	k8s_installer "github.com/Killer-Feature/PaaS_ClientSide/pkg/k8s-installer"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib/ubuntu"
	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
	"go.uber.org/zap"
)

type Service struct {
	r  internal.Repository
	l  *zap.Logger
	tm *taskmanager.Manager[netip.AddrPort]
	hi *helm.HelmInstaller

	k8sInstaller *k8s_installer.Installer
}

func NewService(r internal.Repository, l *zap.Logger, tm *taskmanager.Manager[netip.AddrPort], k8sInstaller *k8s_installer.Installer, hi *helm.HelmInstaller) internal.Usecase {
	return &Service{
		r:            r,
		l:            l,
		tm:           tm,
		k8sInstaller: k8sInstaller,
		hi:           hi,
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

	for i, node := range nodes {
		respNodes[i] = internal.Node{
			ID:        node.ID,
			IP:        node.IP,
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
	return int(taskID), err
}

func (s *Service) addNodeToCurrentClusterProgressTask(ctx context.Context, node internal.FullNode) func(taskId taskmanager.ID) error {
	return func(taskID taskmanager.ID) error {
		sshBuilder := ssh.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(node.IP, node.Login, node.Password)
		if err != nil {
			return err
		}
		defer func(cc cconn.ClientConn) {
			_ = cc.Close()
		}(cc)
		err = s.k8sInstaller.InstallK8S(cc, node.ID)
		if err != nil {
			return err
		}

		return nil
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
		return nil, err
	}

	defer func(cc cconn.ClientConn) {
		_ = cc.Close()
	}(cc)

	output, err := s.getAdminConf(ctx, cc)

	if err == nil {
		adminConf := string(output)
		_ = s.r.UpdateAdminConf(ctx, clusterId, adminConf)
		return &models.AdminConfig{Config: adminConf}, nil
	}

	adminConf, err := s.r.GetAdminConf(ctx, clusterId)
	return &models.AdminConfig{Config: adminConf}, err
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

func (s *Service) GetResources(ctx context.Context) ([]internal.Resourse, error) {
	resources, err := s.hi.GetResourcesList()

	if err != nil {
		return nil, err
	}

	resourceList := make([]internal.Resourse, 0, len(resources))

	for _, res := range resources {
		resourceList = append(resourceList, internal.Resourse{
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

func (s *Service) RemoveNodeFromCurrentCluster(ctx context.Context, id int) (int, error) {
	node, err := s.r.GetFullNode(ctx, id)
	if err != nil {
		return 0, err
	}

	taskID, err := s.tm.AddTask(s.removeNodeFromCurrentClusterProgressTask(context.Background(), node), node.IP)
	return int(taskID), err
}

func (s *Service) removeNodeFromCurrentClusterProgressTask(ctx context.Context, node internal.FullNode) func(taskId taskmanager.ID) error {
	return func(taskID taskmanager.ID) error {
		sshBuilder := ssh.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(node.IP, node.Login, node.Password)
		if err != nil {
			return err
		}
		defer func(cc cconn.ClientConn) {
			_ = cc.Close()
		}(cc)
		err = s.k8sInstaller.RemoveK8S(cc)
		if err != nil {
			return err
		}

		defer func(r internal.Repository, ctx context.Context, id int, clusterID int) {
			_ = r.SetNodeClusterID(ctx, id, clusterID)
		}(s.r, ctx, node.ID, 0)

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
