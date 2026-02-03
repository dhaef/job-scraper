package dom

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func GetTextContent(n *html.Node) string {
	var text strings.Builder

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return strings.TrimSpace(text.String())
}

func GetHref(n *html.Node) (string, bool) {
	if n == nil {
		return "", false
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return attr.Val, true
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if val, ok := GetHref(c); ok {
			return val, true
		}
	}

	return "", false
}

func RenderNode(n *html.Node) (string, error) {
	var buf bytes.Buffer

	if err := html.Render(&buf, n); err != nil {
		return "", err
	}

	return buf.String(), nil
}
