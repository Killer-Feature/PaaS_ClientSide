package internal

import (
	"context"
	"errors"
	"net/netip"
)

type Usecase interface {
	ExecCommand(command string) ([]byte, error)
	GetClusterNodes(ctx context.Context) ([]Node, error)
	AddNode(ctx context.Context, node FullNode) (int, error)
	RemoveNode(ctx context.Context, id int) error
	AddNodeToCurrentCluster(ctx context.Context, id int) (int, error)
}

var (
	ErrNodeExists = errors.New("node with current ip exists")
)

type Node struct {
	ID   int
	IP   netip.AddrPort
	Name string
}
