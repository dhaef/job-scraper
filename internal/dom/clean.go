package dom

import "golang.org/x/net/html"

func CleanNodes(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "noscript" || n.Data == "style" || n.Data == "iframe" || n.Data == "form") {
		n.Parent.RemoveChild(n)
		return nil
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		CleanNodes(c)
		c = next
	}
	return n
}
