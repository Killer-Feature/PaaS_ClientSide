package internal

import (
	"context"
	"errors"
	"github.com/Killer-Feature/PaaS_ClientSide/internal/models"
	"net/netip"
)

type Usecase interface {
	ExecCommand(command string) ([]byte, error)
	GetClusterNodes(ctx context.Context) ([]Node, error)
	AddNode(ctx context.Context, node FullNode) (int, error)
	RemoveNode(ctx context.Context, id int) error
	AddNodeToCurrentCluster(ctx context.Context, id int) (int, error)
	AddResource(ctx context.Context, rType ResourceType, name string) error
	RemoveResource(ctx context.Context, rType ResourceType, name string) error
	GetAdminConfig(ctx context.Context, clusterId int) (*models.AdminConfig, error)
}

var (
	ErrNodeExists = errors.New("node with current ip exists")
)

type Node struct {
	ID   int            `json:"id"`
	IP   netip.AddrPort `json:"ip"`
	Name string         `json:"name"`
}

type ResourceType int

const (
	Undefined ResourceType = iota
	Postgres
	Redis
)
