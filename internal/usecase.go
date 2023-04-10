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
	GetResources(ctx context.Context) ([]Resourse, error)
}

var (
	ErrNodeExists = errors.New("node with current ip exists")
)

type Node struct {
	ID        int            `json:"id"`
	IP        netip.AddrPort `json:"ip"`
	Name      string         `json:"name"`
	ClusterID int            `json:"clusterID"`
	IsMaster  bool           `json:"isMaster"`
}

type ResourceType int

type Resourse struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	FirstDeployed string `json:"firstDeployed"`
	LastDeployed  string `json:"lastDeployed"`
	AppVersion    string `json:"appVersion"`
	ApiVersion    string `json:"apiVersion"`
	Description   string `json:"description"`
	ChartVersion  string `json:"chartVersion"`
}

const (
	Undefined ResourceType = iota
	Postgres
	Redis
	Prometheus
	Grafana
	NginxIngressController
)
