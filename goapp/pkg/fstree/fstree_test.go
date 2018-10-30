package fstree_test

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/tonto/kit/goapp/pkg/fstree"
)

func TestTreeWalk(t *testing.T) {
	cases := []struct {
		name string
	}{
		{
			name: "service",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(filepath.Join("testdata", tc.name+".json"))
			checkFatal(t, err)

			tree, err := fstree.NewFromJSON(data)
			checkFatal(t, err)

			tree.Walk(func(n *fstree.Node) error {
				t.Logf("path: %s - name: %s\n", n.Path(), n.Name())
				return nil
			})

			log.Println(tree)
			t.Fail()
		})
	}
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
