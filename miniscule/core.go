package miniscule

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// func resolve(node *yaml.Node, resolver func(*yaml.Node) error) (*yaml.Node, error) {
// 	fmt.Printf("TOP: Tag %s Value %s Kind %#v\n", node.Tag, node.Value, node.Kind)
// 	for _, n := range node.Content {
// 		fmt.Printf("Tag %s Value %s Kind %#v\n", n.Tag, n.Value, n.Kind)
// 	}
// }

type Error struct {
	Tag     string
	Message string
	Line    int
	Column  int
}

func (err Error) Error() string {
	return err.Message
}

// EnvResolver TO
func EnvResolver(cont cont, node *yaml.Node) (*yaml.Node, error) {
	if node.Tag == "!env" {
		value, ok := os.LookupEnv(node.Value)
		if ok {
			fmt.Println(value)
			node.Value = value
			return node, nil
		}

		return nil, nil
	}
	return node, nil
}

// Resolver TODO
type cont = func(*yaml.Node) (*yaml.Node, error)

// Resolver TODO
type Resolver = func(cont, *yaml.Node) (*yaml.Node, error)

// BaseResolver TODO
func BaseResolver(cont cont, node *yaml.Node) (*yaml.Node, error) {
	for i := range node.Content {
		n, err := cont(node.Content[i])
		if err != nil {
			return nil, err
		}
		node.Content[i] = n
	}
	return node, nil
}

// OrResolver TODO
func OrResolver(cont cont, node *yaml.Node) (*yaml.Node, error) {
	if node.Tag == "!or" && node.Kind == yaml.SequenceNode {
		for i := range node.Content {
			child, err := cont(node.Content[i])
			if err != nil {
				return nil, err
			}
			if child != nil {
				return child, nil
			}
		}
		return nil, nil
	}
	return node, nil
}

// Sequence TODO
func Sequence(resolvers ...Resolver) Resolver {
	return func(cont cont, node *yaml.Node) (*yaml.Node, error) {
		var n *yaml.Node = node
		var err error
		for i := range resolvers {
			n, err = resolvers[i](cont, n)
			if err != nil {
				return nil, err
			}
			if n == nil {
				return &yaml.Node{Kind: yaml.ScalarNode, Value: ""}, nil
			}
		}
		return n, nil
	}
}

// Resolve TODO
func Resolve(node *yaml.Node) (*yaml.Node, error) {
	seq := Sequence(EnvResolver, OrResolver, BaseResolver)
	return ResolveWith(seq, node)
}

// ResolveWith TODO
func ResolveWith(r Resolver, node *yaml.Node) (*yaml.Node, error) {
	var cont cont = func(n *yaml.Node) (*yaml.Node, error) {
		return ResolveWith(r, n)
	}
	return r(cont, node)
}

// func Unmarshal(in []byte, out interface{}) error {
// 	var original yaml.Node
// 	err := yaml.Unmarshal(in, &original)
// 	if err != nil {
// 		return err
// 	}
// 	resolved, err := resolve(&original)
// 	if err != nil {
// 		return err
// 	}
// 	err = node.Decode(out)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
