package internal

import (
	"context"
	"net/netip"

	"github.com/Killer-Feature/PaaS_ClientSide/internal/models"
)

type Repository interface {
	GetNodes(ctx context.Context) ([]FullNode, error)
	GetFullNode(ctx context.Context, id int) (FullNode, error)
	AddNode(ctx context.Context, node FullNode) (int, error)
	RemoveNode(ctx context.Context, id int) error
	IsNodeExists(ctx context.Context, ip netip.AddrPort) (int, error)
	SetNodeClusterID(ctx context.Context, id int, clusterID int) error

	AddResource(ctx context.Context, rType, name string) error
	GetResources(ctx context.Context) ([]models.ResourceData, error)

	AddCluster(ctx context.Context, clusterName string) (int, error)
	GetClusterID(ctx context.Context, clusterName string) (int, error)
	GetClusterName(ctx context.Context, id int) (string, error)
	AddClusterTokenIPAndHash(ctx context.Context, clusterID int, token, masterIP, hash string) error
	CheckClusterTokenIPAndHash(ctx context.Context, clusterID int) (bool, error)
	GetClusterTokenIPAndHash(ctx context.Context, clusterID int) (token, masterIP, hash string, err error)
	DeleteClusterTokenIPAndHash(ctx context.Context, clusterID int) (err error)
	UpdateAdminConf(ctx context.Context, clusterID int, adminConf string) error
	GetAdminConf(ctx context.Context, clusterID int) (conf string, err error)
}

type FullNode struct {
	ID        int
	Name      string
	IP        netip.AddrPort
	Login     string
	Password  string
	ClusterID int
	IsMaster  bool
}
