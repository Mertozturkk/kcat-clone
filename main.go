package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
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

	scanNodeMappings(content, `$root`)

	colorizeKeys(content, `$root`)

	colorizeComments(content)

	// encode yaml
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	enc.Encode(content)

	fmt.Println(render(buf))
}

var keyPathMap = make(map[*yaml.Node]string)
var kvMpa = make(map[*yaml.Node]*yaml.Node)
var vkMap = make(map[*yaml.Node]*yaml.Node)

func scanNodeMappings(content *yaml.Node, s string) {
	if content.Kind == yaml.MappingNode {
		for i, child := range content.Content {
			if i%2 == 0 {
				keyPathMap[child] = fmt.Sprintf("%s.%s", s, child.Value)
				kvMpa[child] = content.Content[i+1]
				vkMap[content.Content[i+1]] = child
			}
		}
	}
	for _, child := range content.Content {
		key := s
		keyNode, ok := vkMap[child]
		if ok {
			key = keyPathMap[keyNode]

		}
		scanNodeMappings(child, key)
	}
}

func markComments(in string) string {
	re := regexp.MustCompile(`(?m)^#(.*)`)
	return re.ReplaceAllString(in, "#COMMENT_$1")
}

func colorizeComments(node *yaml.Node) {
	for _, child := range node.Content {
		child.HeadComment = markComments(child.HeadComment)
		child.LineComment = markComments(child.LineComment)
		child.FootComment = markComments(child.FootComment)

		colorizeComments(child)
	}
}
func render(buf bytes.Buffer) string {
	s := buf.String()
	// render keys
	s = regexp.MustCompile(`(?m)(KEY_BLUE_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgBlue, color.Bold).Sprint(`$2$3`))
	s = regexp.MustCompile(`(?m)(KEY_YELLOW_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgYellow, color.Bold).Sprint(`$2$3`))
	s = regexp.MustCompile(`(?m)(KEY_RED_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgRed, color.Bold).Sprint(`$2$3`))
	s = regexp.MustCompile(`(?m)(KEY_GRAY_)([^:]+)`).
		ReplaceAllString(s, color.New(color.FgHiBlack, color.Bold).Sprint(`$2$3`))

	// render comments
	s = regexp.MustCompile(`(?m)#COMMENT_(.*)`).
		ReplaceAllString(s, color.New(color.FgHiBlack).Sprint(`#$1`))

	// render comments
	s = regexp.MustCompile(`(?m)(#)(COMMENT_)(.*$)`).
		ReplaceAllString(s, color.New(color.FgHiBlack).Sprint(`#$1`))

	return s
}
func colorizeKeys(node *yaml.Node, path string) {
	var prevKey string
	for i, child := range node.Content {
		if node.Kind == yaml.SequenceNode && child.Kind == yaml.ScalarNode {
			continue
		}
		if i%2 == 0 && child.Value != "" {
			keyPath := path + "." + child.Value
			prevKey = child.Value
			child.Value = "KEY_" + colorForKey(keyPath) + "_" + child.Value
		}
		subPath := path
		if node.Kind == yaml.MappingNode {
			subPath = path + "." + prevKey
		}
		colorizeKeys(child, subPath)
	}

}

func colorForKey(path string) string {
	redSuffixes := []string{"$root.apiVersion",
		"$root.kind",
		".spec",
		"$root.metadata.name",
		".containers.name",
		".containers.image"}
	for _, f := range redSuffixes {
		if strings.HasSuffix(path, f) {
			return "RED"
		}

	}
	if strings.HasPrefix(path, "$root.metadata") {
		return "YELLOW"
	}

	if strings.HasPrefix(path, "$root.spec") {
		return "BLUE"
	}

	if strings.HasPrefix(path, "$root.status") {
		return "GRAY"
	}

	return "GRAY"
}
