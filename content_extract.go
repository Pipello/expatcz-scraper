package expatczscraper

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/html"
)


const BaseURL = "https://www.expats.cz"

func GetArticleContent(link string) ([]Section, error) {
	response, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}

	htmlBody, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	return ExtractArticleContentWithTitle(htmlBody), nil
}

func FindLinkWith(linkContains string, pageLink string) (string, error) {
	response, err := http.Get(pageLink)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}

	htmlBody, err := html.Parse(response.Body)
	if err != nil {
		return "", err
	}
	l := FindFirstLinkWithContent(htmlBody, linkContains)
	if l == "" {
		return "", errors.New("no link found")
	}
	return BaseURL+l, err
}

func FindFirstLinkWithContent(n *html.Node, content string) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && strings.Contains(a.Val, content) {
				return a.Val
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if l := FindFirstLinkWithContent(c, content); l != "" {
			return l
		}
	}
	return ""
}

type Section struct {
	Title   string
	Content string
}

func ExtractArticleContentWithTitle(n *html.Node) []Section {
	if n.Type == html.ElementNode && n.Data == "div" && slices.ContainsFunc(n.Attr, containContentClass) {
		ct := []Section{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "div" && slices.ContainsFunc(c.Attr, containTitleClass) {
				ct = append(ct, Section{
					Title:   extractText(c),
					Content: extractContent(c.NextSibling),
				})
			}
		}
		if len(ct) > 0 {
			return ct
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if ct := ExtractArticleContentWithTitle(c); len(ct) > 0 {
			return ct
		}
	}
	return nil
}

func containTextWrapperClass(a html.Attribute) bool {
	return a.Key == "class" && strings.Contains(a.Val, "widget text")
}

func containTitleClass(a html.Attribute) bool {
	return a.Key == "class" && strings.Contains(a.Val, "headinglevel2")
}

func containContentClass(a html.Attribute) bool {
	return a.Key == "class" && strings.Contains(a.Val, "content")
}

func extractContent(n *html.Node) string {
	r := ""
	for c := n; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "div" && slices.ContainsFunc(c.Attr, containTextWrapperClass) {
			r += extractText(c.FirstChild)
		}
		if c.Type == html.ElementNode && c.Data == "div" && slices.ContainsFunc(c.Attr, containTitleClass) {
			break
		}
	}
	return r
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	r := ""
	if n.Type == html.ElementNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			r += extractText(c)
		}
	}
	return r
}
