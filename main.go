package main

import (
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {

	file := os.Args[1]

	data, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	// decode yaml
	var v yaml.Node
	if err := yaml.Unmarshal(data, &v); err != nil {
		panic(err)
	}

	if len(v.Content) == 0 {
		panic("No yaml docs found")
	}

	content := v.Content[0]

	colorizeKeys(content)

	// encode yaml
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	if err := enc.Encode(content); err != nil {
		panic(err)
	}

}
func colorizeKeys(node *yaml.Node) {

	for i, child := range node.Content {
		if node.Kind == yaml.SequenceNode && child.Kind == yaml.ScalarNode {
			continue
		}
		if i%2 == 0 && child.Value != "" {
			child.Value = color.WhiteString(child.Value)
		}
		colorizeKeys(child)
	}

}
