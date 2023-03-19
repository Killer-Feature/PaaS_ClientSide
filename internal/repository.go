package internal

import (
	"context"
	"net/netip"
)

type Repository interface {
	GetNodes(ctx context.Context) ([]FullNode, error)
	AddNode(ctx context.Context, node FullNode) (int, error)
	RemoveNode(ctx context.Context, id int) error
	IsNodeExists(ctx context.Context, ip netip.AddrPort) (bool, error)
}

type FullNode struct {
	ID       int
	Name     string
	IP       netip.AddrPort
	Login    string
	Password string
}
