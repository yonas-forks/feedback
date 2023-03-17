package htmlgen

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func PrettyPrint() {
	input := "<html><head><title class=\"dfwfwefe\">Example Page</title></head><body><h1>Hello, world!</h1><p>This is an example page.</p></body></html>"
	doc, err := html.ParseFragment(strings.NewReader(input), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing HTML: %v\n", err)
		os.Exit(1)
	}
	for _, n := range doc {
		indent(os.Stdout, n, 0)
	}
}

func indent(w io.Writer, n *html.Node, depth int) {
	prefix := strings.Repeat(" ", depth*4)
	switch n.Type {
	case html.ElementNode:
		fmt.Fprintf(w, "%s<%s", prefix, n.Data)
		for _, a := range n.Attr {
			fmt.Fprintf(w, " %s=\"%s\"", a.Key, a.Val)
		}
		if n.FirstChild == nil {
			fmt.Fprint(w, "/>\n")
		} else {
			fmt.Fprint(w, ">\n")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				indent(w, c, depth+1)
			}
			fmt.Fprintf(w, "%s</%s>\n", prefix, n.Data)
		}
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Fprintf(w, "%s%s\n", prefix, text)
		}
	}
}
