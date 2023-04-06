package service

import (
	"context"
	"errors"
	"github.com/Killer-Feature/PaaS_ClientSide/internal/models"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/os_command_lib/ubuntu"
	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"strconv"

	"github.com/Killer-Feature/PaaS_ClientSide/pkg/helm"
	k8s_installer "github.com/Killer-Feature/PaaS_ClientSide/pkg/k8s-installer"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
	"go.uber.org/zap"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/executor"

	"github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"

	_ "github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
)

type Service struct {
	r  internal.Repository
	l  *zap.Logger
	tm *taskmanager.Manager
	hi *helm.HelmInstaller

	k8sInstaller *k8s_installer.Installer
}

func NewService(r internal.Repository, l *zap.Logger, tm *taskmanager.Manager, k8sInstaller *k8s_installer.Installer, hi *helm.HelmInstaller) internal.Usecase {
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
			ID:   node.ID,
			IP:   node.IP,
			Name: node.Name,
		}
	}

	return respNodes, nil
}

func (s *Service) AddNodeToCurrentCluster(ctx context.Context, id int) (int, error) {
	// TODO: Add task

	node, err := s.r.GetFullNode(ctx, id)
	if err != nil {
		return 0, err
	}

	taskID, err := s.tm.AddTask(s.k8sInstaller.InstallK8S, node.IP, taskmanager.AuthData{
		Login:    node.Login,
		Password: node.Password,
	})
	return int(taskID), err
}

func (s *Service) AddNode(ctx context.Context, node internal.FullNode) (int, error) {
	exists, err := s.r.IsNodeExists(ctx, node.IP)
	if err != nil {
		return 0, err
	}
	if exists == false {
		return s.r.AddNode(ctx, node)
	}
	return 0, internal.ErrNodeExists
}

func (s *Service) RemoveNode(ctx context.Context, id int) error {
	return s.r.RemoveNode(ctx, id)
}

func (s *Service) AddResource(ctx context.Context, rType internal.ResourceType, name string) error {
	switch rType {
	case internal.Postgres:
		return s.hi.Install(name, "postgresql")
	default:
		return errors.New("resource not implemented")
	}
}

func (s *Service) RemoveResource(ctx context.Context, rType internal.ResourceType, name string) error {
	switch rType {
	case internal.Postgres:
		return s.hi.UninstallChart(name)
	default:
		return errors.New("resource not implemented")
	}
}

func (s *Service) GetAdminConfig(ctx context.Context, clusterId int) (*models.AdminConfig, error) {
	_, masterIdStr, _, err := s.r.GetClusterTokenIPAndHash(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	masterId, err := strconv.Atoi(masterIdStr)

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

	cl := ubuntu.Ubuntu2004CommandLib{}

	getAdminConfCommand := cl.CatAdminConfFile()

	output, err := cc.Exec(string(getAdminConfCommand.Command))

	if err == nil {
		adminConf := string(output)
		_ = s.r.UpdateAdminConf(ctx, clusterId, adminConf)
		return &models.AdminConfig{Config: adminConf}, nil
	}

	adminConf, err := s.r.GetAdminConf(ctx, clusterId)
	return &models.AdminConfig{Config: adminConf}, err
}
