package main

import (
	"fmt"
	"github.com/bartfrenk/miniscule-go/miniscule"
	"gopkg.in/yaml.v3"
)

func walk(depth int, node *yaml.Node, do func(int, *yaml.Node)) {
	do(depth, node)
	for i := range node.Content {
		walk(depth+1, node.Content[i], do)

	}
}

var m = map[yaml.Kind]string{
	yaml.DocumentNode: "DocumentNode",
	yaml.SequenceNode: "SequenceNode",
	yaml.MappingNode:  "MappingNode",
	yaml.ScalarNode:   "ScalarNode",
	yaml.AliasNode:    "AliasNode",
}

func printNode(depth int, node *yaml.Node) {
	for i := 0; i < depth; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("Kind %s, Value %s\n", m[node.Kind], node.Value)
}

func main() {
	var node yaml.Node
	data := []byte("foo: !or []")
	err := yaml.Unmarshal(data, &node)
	n, err := miniscule.Resolve(&node)
	if err != nil {
		panic(err)
	}
	walk(0, n, printNode)
}
