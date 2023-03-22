package internal

import (
	"context"
	"net/netip"
)

type Repository interface {
	GetNodes(ctx context.Context) ([]FullNode, error)
	GetFullNode(ctx context.Context, id int) (FullNode, error)
	AddNode(ctx context.Context, node FullNode) (int, error)
	RemoveNode(ctx context.Context, id int) error
	IsNodeExists(ctx context.Context, ip netip.AddrPort) (bool, error)

	AddCluster(ctx context.Context, clusterName string) (int, error)
	GetClusterID(ctx context.Context, clusterName string) (int, error)
	GetClusterName(ctx context.Context, clusterName string) (int, error)
	AddClusterTokenIPAndHash(ctx context.Context, clusterID int, token, masterIP, hash string) error
	CheckClusterTokenIPAndHash(ctx context.Context, clusterID int) (bool, error)
	GetClusterTokenIPAndHash(ctx context.Context, clusterID int) (token, masterIP, hash string, err error)
}

type FullNode struct {
	ID       int
	Name     string
	IP       netip.AddrPort
	Login    string
	Password string
}
