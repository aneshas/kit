package command

import (
	"fmt"

	"github.com/tonto/kit/goapp/pkg/fstree"
	"github.com/tonto/kit/goapp/pkg/template"
)

type app struct {
	Name string
	// ...
}

// NewCreateApp creates new CreateApp command instance
// CreateApp can generate different project layouts based on the provided tplPath
// TODO - pass options
func NewCreateApp(name string, tpl string) (*CreateApp, error) {
	var cmd CreateApp

	data, err := template.Parse(
		"app/"+tpl,
		app{
			Name: name,
		},
	)
	if err != nil {
		return nil, err
	}

	ct, err := newCreateTree(data, cmd.filter)
	if err != nil {
		return nil, err
	}

	cmd.createTree = ct
	cmd.name = name

	return &cmd, nil
}

// CreateApp command
type CreateApp struct {
	name string
	*createTree
}

// Execute create app command
func (ca *CreateApp) Execute() error {
	fmt.Printf("Creating new goapp application under %s directory\n\n", ca.name)
	return ca.createTree.Execute()
}

// TODO - filter should also return raw file to be created after going through the template
func (cmd *CreateApp) filter(n *fstree.Node) ([]byte, bool) {
	return []byte("// Autogenerated by Goapp"), true
}
