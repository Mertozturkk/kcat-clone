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
	var v map[string]interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		panic(err)
	}

	// print type
	// fmt.Printf("%T\n", v) -> map[string]interface {}
}
