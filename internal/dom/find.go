package dom

import (
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func FindElementsByClasses(n *html.Node, classNames []string) []*html.Node {
	var matchingNodes []*html.Node

	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				for c := range strings.FieldsSeq(attr.Val) {
					if slices.Contains(classNames, c) {
						matchingNodes = append(matchingNodes, n)
						break
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		matchingNodes = append(matchingNodes, FindElementsByClasses(c, classNames)...)
	}

	return matchingNodes
}

func FindElementsByTagName(n *html.Node, tagName string) []*html.Node {
	var elements []*html.Node
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tagName {
			elements = append(elements, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	traverse(n)
	return elements
}

func FindLabelByClass(n *html.Node, labels map[string][]string) (string, bool) {
	if n == nil {
		return "", false
	}

	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				for label, values := range labels {
					for _, v := range values {
						if strings.Contains(attr.Val, v) {
							return label, true
						}
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		label, found := FindLabelByClass(c, labels)
		if found {
			return label, true
		}
	}

	return "", false
}
