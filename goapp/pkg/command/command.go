package command

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tonto/kit/goapp/pkg/fstree"
)

// Command interface
type Command interface {
	Execute() error
	Rollback() error
}

func newCreateTree(data []byte, f func(*fstree.Node) ([]byte, bool)) (*createTree, error) {
	tree, err := fstree.NewFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &createTree{
		tree: tree,
		f:    f,
	}, nil
}

type createTree struct {
	tree *fstree.FSTree
	f    func(*fstree.Node) ([]byte, bool)
}

// Execute command
func (ct *createTree) Execute() error {
	err := ct.tree.Walk(
		func(n *fstree.Node) error {
			data, ok := ct.f(n)
			if !ok {
				// TODO - mark as node not printable
				return nil
			}

			fmt.Printf("creating %s", n.Path())

			err := ct.createNode(n, data)
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		fmt.Println("")
		return err
	}

	fmt.Println(ct.tree)
	fmt.Printf("Created 0 files 5 dirs.\n")

	return nil
}

func (ct *createTree) createNode(n *fstree.Node, data []byte) error {
	var err error
	switch n.T() {
	case fstree.Dir:
		err = ct.createDir(n.Path())
	case fstree.File:
		err = ct.createFile(n.Path(), data)
	default:
		return fmt.Errorf("command: invalid node type")
	}
	return err
}

func (ct *createTree) createDir(dir string) error {
	err := os.Mkdir(dir, 0766)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf(" [Node Exists]\n")
			return nil
		}
		return err
	}

	fmt.Printf(" [Done]\n")

	return nil
}

func (ct *createTree) createFile(file string, data []byte) error {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		err = ioutil.WriteFile(file, data, 0766)
		if err != nil {
			return err
		}
		fmt.Printf(" [Done]\n")
		return nil
	}

	if err == nil {
		fmt.Printf(" [Node Exists]\n")
	}

	return err
}

// Rollback cleans up anything created by execute
func (ct *createTree) Rollback() error {
	return nil
}
