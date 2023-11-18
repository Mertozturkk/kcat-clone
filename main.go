package main

import (
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

	// encode yaml
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	if err := enc.Encode(content); err != nil {
		panic(err)
	}
}
