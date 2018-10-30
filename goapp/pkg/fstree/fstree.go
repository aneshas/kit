package fstree

import (
	"encoding/json"

	"github.com/xlab/treeprint"
)

// New constructs new FSTree
func New(ptree *AppTree) (*FSTree, error) {
	ntree, err := ptree.newNodeTree()
	if err != nil {
		return nil, err
	}

	return &FSTree{
		tree: ntree,
	}, nil
}

// NewFromJSON constructs new FSTree from json bytes
func NewFromJSON(bytes []byte) (*FSTree, error) {
	var pt AppTree

	err := json.Unmarshal(bytes, &pt)
	if err != nil {
		return nil, err
	}

	return New(&pt)
}

// AppTree represents a tree of app nodes (project, service, ...)
type AppTree struct {
	App AppNode `json:"app"`
}

// AppNode represents App node
// It is used to create and update file system trees eg. for App,
// service, or anything else that has a specific file system tree
type AppNode struct {
	Name  string            `json:"name"`
	Meta  map[string]string `json:"meta"`
	T     NodeType          `json:"t"`
	Nodes []AppNode         `json:"nodes"`
}

func (at *AppTree) newNodeTree() (*NodeTree, error) {
	tree := NodeTree{
		root: Node{
			name: at.App.Name,
			meta: at.App.Meta,
			t:    at.App.T,
			path: at.App.Name,
		},
	}

	var err error

	for _, an := range at.App.Nodes {
		err = at.walk(&an, &tree.root)
		if err != nil {
			return nil, err
		}
	}

	return &tree, nil
}

func (at *AppTree) walk(appNode *AppNode, node *Node) error {
	nextRoot, err := node.Add(appNode.Name, appNode.T, appNode.Meta)
	if err != nil {
		return err
	}

	for _, an := range appNode.Nodes {
		at.walk(&an, nextRoot)
	}

	return nil
}

// FSTree represents file system tree abstraction
type FSTree struct {
	tree *NodeTree
}

// Walk walks through the nodes in the tree depth first
func (ft *FSTree) Walk(f func(*Node) error) error {
	var err error

	err = f(&ft.tree.root)
	if err != nil {
		return err
	}

	for _, n := range ft.tree.root.nodes {
		err = ft.walk(n, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ft *FSTree) walk(root *Node, f func(*Node) error) error {
	err := f(root)
	if err != nil {
		return err
	}

	for _, n := range root.nodes {
		err = ft.walk(n, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// String prints file tree
func (ft *FSTree) String() string {
	tree := treeprint.New()

	app := tree.AddBranch(ft.tree.root.Name())

	for _, n := range ft.tree.root.nodes {
		ft.stringWalk(n, app)
	}

	return tree.String()
}

func (ft *FSTree) stringWalk(root *Node, b treeprint.Tree) {
	switch root.t {
	case Dir:
		b = b.AddBranch(root.Name())
	default:
		b.AddNode(root.Name())
	}

	for _, n := range root.nodes {
		ft.stringWalk(n, b)
	}
}
