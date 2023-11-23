package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
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
	colorizeComments(content)

	// encode yaml
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	enc.Encode(content)

	fmt.Println(render(buf))
}

func mark(in string) string {
	re := regexp.MustCompile(`(?m)^#(.*)`)
	return re.ReplaceAllString(in, "#COMMENT_$1")
}

func colorizeComments(node *yaml.Node) {
	for _, child := range node.Content {
		colorizeComments(child)
	}
}
func render(buf bytes.Buffer) string {
	s := buf.String()
	// render keys
	re := regexp.MustCompile(`(?m)(KEY_)([^:]+)`)
	s = re.ReplaceAllString(s, color.RedString("$2$3"))

	// render comments
	s = regexp.MustCompile(`(?m)(#)(COMMENT_)(.*$)`).ReplaceAllString(s, color.YellowString("$1$3"))

	return s
}
func colorizeKeys(node *yaml.Node) {

	for i, child := range node.Content {
		if node.Kind == yaml.SequenceNode && child.Kind == yaml.ScalarNode {
			continue
		}
		if i%2 == 0 && child.Value != "" {
			child.Value = "KEY_" + child.Value
		}
		colorizeKeys(child)
	}

}
