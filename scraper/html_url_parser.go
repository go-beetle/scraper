package scraper

import (
	"log"
	"strings"

	"golang.org/x/net/html"
)

func GetHref(h string, u *URL) []*URL {
	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		log.Fatal(err)
	}

	urls := make([]*URL, 0)
	extractorHelper(doc, u, &urls)
	return urls
}

func extractorHelper(n *html.Node, u *URL, urls *[]*URL) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				urlRef := ParseReferencedUrl(u, a.Val)
				if urlRef != nil && urlRef.Hostname == u.Hostname {
					*urls = append(*urls, urlRef)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractorHelper(c, u, urls)
	}
}
