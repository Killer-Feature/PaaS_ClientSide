package service

import (
	"context"

	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
	"go.uber.org/zap"

	"github.com/Killer-Feature/PaaS_ClientSide/internal"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/executor"

	_ "github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
)

type Service struct {
	r  internal.Repository
	l  *zap.Logger
	tm *taskmanager.Manager
}

func NewService(r internal.Repository, l *zap.Logger) internal.Usecase {
	return &Service{
		r: r,
		l: l,
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

func (s *Service) AddNodeToCurrentCluster(ctx context.Context, id int) {
	// TODO: Add task

	// task, err := s.tm.AddTask()
	//if err != nil {
	//	return
	//}

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
