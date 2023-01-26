package wrappers

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/openrelayxyz/plugeth-utils/core"
)

type Node struct {
	n *node.Node
}

func NewNode(n *node.Node) *Node {
	return &Node{n}
}

func (n *Node) Server() core.Server {
	return n.n.Server()
}

func (n *Node) DataDir() string {
	return n.n.DataDir()
}
func (n *Node) InstanceDir() string {
	return n.n.InstanceDir()
}
func (n *Node) IPCEndpoint() string {
	return n.n.IPCEndpoint()
}
func (n *Node) HTTPEndpoint() string {
	return n.n.HTTPEndpoint()
}
func (n *Node) WSEndpoint() string {
	return n.n.WSEndpoint()
}
func (n *Node) ResolvePath(x string) string {
	return n.n.ResolvePath(x)
}

func (n *Node) Attach() (core.Client, error) {
	return n.n.Attach()
}
func (n *Node) Close() (error) {
	return n.n.Close()
}
