package fstree

import (
	"errors"
)

var (
	// ErrNodeExists is thrown when trying to add a subnode
	// if a subnode with the same name exists
	ErrNodeExists = errors.New("filetree: sub node with the same name already exists")
)

// NodeType type
type NodeType int

const (
	// Dir - directory node type
	Dir NodeType = iota

	// File - file node type
	File
)

// NodeTree represents a tree of nodes
type NodeTree struct {
	root Node
}

// Node represents project node
type Node struct {
	name  string
	path  string
	meta  map[string]string
	t     NodeType
	nodes []*Node
}

// Name returns node name
func (n *Node) Name() string {
	return n.name
}

// Path returns node path
func (n *Node) Path() string {
	return n.path
}

// ByPath fetches sub-node by its path
func (n *Node) ByPath(path string) (*Node, error) {
	panic("not impl")
}

// Remove removes sub-node by name
func (n *Node) Remove(name string) error {
	panic("not impl")
}

// T returns node type
func (n *Node) T() NodeType {
	return n.t
}

// Add adds new sub-node ensuring name is unique
func (n *Node) Add(name string, t NodeType, meta map[string]string) (*Node, error) {
	for _, n := range n.nodes {
		if n.name == name {
			return nil, ErrNodeExists
		}
	}

	node := Node{
		name: name,
		t:    t,
		meta: meta,
		path: n.path + "/" + name,
	}

	n.nodes = append(n.nodes, &node)

	return &node, nil
}
