package compile

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Here ends the possession of the light

const (
	gOn      = "g-on:"
	gOnAlt   = "@"
	gBind    = "g-bind:"
	gBindAlt = ":"
	gDirs    = "g-"
	gFor     = "g-for"
)

// little alias
var trim = strings.TrimSpace

func Parse(reader io.Reader) ([]*html.Node, error) {
	return html.ParseFragment(reader, &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
}

func ParseTemplate(dir string, nodes []*html.Node) (*html.Node, map[string]string, error) {
	var tags = make(map[string]string)
	var template *html.Node
	for _, node := range nodes {
		if node.Type != html.ElementNode {
			continue
		}

		srcAttr := getXAttr(node.Attr, "src")
		if len(srcAttr) != 0 {
			fileBody, err := ioutil.ReadFile(dir + "/" + srcAttr)
			if err != nil {
				return nil, nil, err
			}

			if node.Data == "template" {
				template, err = html.Parse(bytes.NewBuffer(fileBody))
				if err != nil {
					return nil, nil, err
				}
			} else {
				tags[node.Data] = string(fileBody)
			}

			continue
		}

		if node.Data == "template" {
			template = node
		}

		tags[node.Data] = node.FirstChild.Data
	}

	if template == nil || tags["script"] == "" {
		return nil, nil, fmt.Errorf("invalid file syntax: %v", nodes)
	}

	return template, tags, nil
}

func TextChild(node *html.Node) string {
	return node.FirstChild.Data
}

func FirstChild(node *html.Node) *html.Node {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			return child
		}
	}
	return nil
}

func parseNode(t string) (*html.Node, error) {
	nodes, err := parseNodes(strings.NewReader(t))
	if err != nil {
		return nil, err
	}

	return nodes[0], nil
}

func parseNodes(reader io.Reader) ([]*html.Node, error) {
	return html.ParseFragment(reader, &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
}

func isExternal(tag string) bool {
	haveToUp := strings.Index(tag, "-") != -1
	haveLowDash := strings.Index(tag, "_") != -1
	haveBrackets := (strings.Index(tag, "(") != -1) || (strings.Index(tag, ")") != -1)
	haveDots := strings.Index(tag, ".") != -1
	haveBlocks := (strings.Index(tag, "[") != -1) || (strings.Index(tag, "]") != -1)
	return haveToUp || haveLowDash || haveBrackets || haveDots || haveBlocks
}

func stringNotEmpty(a string) bool {
	cleared := trim(a)
	return len(cleared) != 0
}

func removeInArr(a []byte, i int) string {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index
	a[len(a)-1] = 0      // Erase last element (write zero value)
	a = a[:len(a)-1]     // Truncate slice
	return string(a)
}
