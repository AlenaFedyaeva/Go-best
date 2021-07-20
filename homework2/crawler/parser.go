package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html"
)
func errcheck(err error) {
    if err != nil {
        log.Println(err)//log.Fatal(err)
    }
}
func dclose(c io.Closer) {
    if err := c.Close(); err != nil {
        log.Println(err)//log.Fatal(err)
    }
}
// парсим страницу
func parse(url string) (*html.Node, error) {
	// что здесь должно быть вместо http.Get? :)
	cli := http.Client{
		Timeout: 20 * time.Second,
	}
	cli.Get(url)

	resp, err := cli.Get(url)//http.Get(url)
	errcheck(err)
	defer dclose(resp.Body)
	

	b, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't parse page")
	}
	
	return b, err
}

// ищем заголовок на странице
func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

// ищем все ссылки на страницы. Используем мапку чтобы избежать дубликатов
func pageLinks(links map[string]struct{}, n *html.Node) map[string]struct{} {
	if links == nil {
		links = make(map[string]struct{})
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			//fmt.Println("a.val",a.Val)
			// костылик для простоты
			if _, ok := links[a.Val]; !ok && len(a.Val) > 2 && a.Val[:2] == "//" {
				
				links["http://"+a.Val[2:]] = struct{}{}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {

		links = pageLinks(links, c)
	}
	return links
}
