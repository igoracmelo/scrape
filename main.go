package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	flag.Parse()

	assert(flag.NArg() >= 1, "first argument must be url")

	u, err := url.Parse(flag.Arg(0))
	assert(err == nil, fmt.Errorf("invalid url: %w", err))

	req, err := http.NewRequest("GET", flag.Arg(0), nil)
	assert(err == nil, err)
	req.Header.Set("Host", u.Host)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "curl/8.4.0")

	resp, err := http.DefaultClient.Do(req)
	assert(err == nil, err)
	defer resp.Body.Close()

	if flag.NArg() == 1 {
		_, err = io.Copy(os.Stdout, resp.Body)
		fmt.Println()
		assert(err == nil, err)
		return
	}

	switch flag.Arg(1) {
	case "first":
		assert(flag.NArg() == 3, "expected query after 'first'")
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		assert(err == nil, err)

		s := doc.Find(flag.Arg(2)).First()
		html, err := s.Html()
		assert(err == nil, err)
		fmt.Println(html)

	case "each":
		assert(flag.NArg() == 3, "expected query after 'each'")

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		assert(err == nil, err)

		doc.Find(flag.Arg(2)).Each(func(i int, s *goquery.Selection) {
			html, err := s.Html()
			assert(err == nil, err)
			fmt.Println(html)
		})

	default:
		panic("unkown option " + flag.Arg(1))
	}
}

func assert(ok bool, msg any) {
	if !ok {
		panic(msg)
	}
}
