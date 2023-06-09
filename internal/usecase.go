package internal

import (
	"context"
	"errors"
	"github.com/Killer-Feature/PaaS_ClientSide/pkg/socketmanager"
	"github.com/gorilla/websocket"
	"net/netip"

	"github.com/Killer-Feature/PaaS_ClientSide/internal/models"
)

// Usecase interface implements functions for Huginn service
// Need to be implemented in main func
type Usecase interface {
	ExecCommand(command string) ([]byte, error)
	GetClusterNodes(ctx context.Context) ([]Node, error)
	AddNode(ctx context.Context, node FullNode) (int, error)
	RemoveNode(ctx context.Context, id int) error
	AddNodeToCurrentCluster(ctx context.Context, id int) (int, error)
	AddResource(ctx context.Context, rType ResourceType, name string) error
	RemoveResource(ctx context.Context, rType ResourceType, name string) error
	GetAdminConfig(ctx context.Context, clusterId int) (*models.AdminConfig, error)
	GetResources(ctx context.Context) ([]Resource, error)
	GetServices(ctx context.Context) ([]Service, error)
	RemoveNodeFromCurrentCluster(ctx context.Context, id int) (int, error)
	GetProgress(ctx context.Context, socket *websocket.Conn) error
	IsAdmin(ctx context.Context, session string) (bool, error)
	Login(ctx context.Context, data LoginData) (string, error)
	Logout(ctx context.Context, session string) error
}

var (
	ErrNodeExists = errors.New("node with current ip exists")
)

type Node struct {
	ID        int            `json:"id"`
	IP        netip.AddrPort `json:"ip"`
	GrafanaIP string         `json:"grafana_ip"`
	Name      string         `json:"name"`
	ClusterID int            `json:"clusterID"`
	IsMaster  bool           `json:"isMaster"`
}

type ResourceType int

type Resource struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	FirstDeployed string `json:"firstDeployed"`
	LastDeployed  string `json:"lastDeployed"`
	AppVersion    string `json:"appVersion"`
	ApiVersion    string `json:"apiVersion"`
	Description   string `json:"description"`
	ChartVersion  string `json:"chartVersion"`
	Type          string `json:"type"`
	ChartURL      string `json:"chartURL"`
}

type Service struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Age       string `json:"age"`
	Created   string `json:"created"`
}

const (
	Undefined ResourceType = iota
	Postgres
	Redis
	Prometheus
	Grafana
	NginxIngressController
	MetalLB
)

type TaskStatus string

const (
	STATUS_IN_QUEUE   TaskStatus = "in queue"
	STATUS_START      TaskStatus = "started"
	STATUS_IN_PROCESS TaskStatus = "in process"
	STATUS_ERROR      TaskStatus = "error"
	STATUS_SUCCESS    TaskStatus = "success"
)

type AddNodeToClusterProgressMsg struct {
	Log     string     `json:"log"`
	Percent int        `json:"percent"`
	Error   string     `json:"error"`
	Status  TaskStatus `json:"status"`
	NodeID  int        `json:"nodeID"`
}

type RemoveNodeFromClusterMsg struct {
	Log     string     `json:"log"`
	Percent int        `json:"percent"`
	Error   string     `json:"error"`
	Status  TaskStatus `json:"status"`
	NodeID  int        `json:"nodeID"`
}

const (
	AddNodeToClusterT      socketmanager.MessageType = "addNodeToCluster"
	RemoveNodeFromClusterT socketmanager.MessageType = "removeNodeFromCluster"
	MetricsT               socketmanager.MessageType = "Metrics"
)

type LoginData struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
