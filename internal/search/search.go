package search

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/alex-held/gofile/internal/log"
)

const (
	perPage = 10

	spClass = "LegacySearchSnippet"
	hdClass = "LegacySearchSnippet-header"
)

func SearchPackages(pkgs ...string) (res []PackageResult) {
	count := len(pkgs)
	log.Log.Debugf("searching for %d packages\n", count)

	pkgC := make(chan *PackageResult, count)
	wg := new(sync.WaitGroup)

	for _, p := range pkgs {
		wg.Add(1)
		go search2(p, pkgC, wg)
	}

	go func() {
		wg.Wait()
		close(pkgC)
	}()

	for p := range pkgC {
		if p.Err != nil {
			log.Log.Error(p.Err)
			continue
		}

		log.Log.Debugf("found url '%s' for pkg '%s'\n", p.Repo, p.Package)
		res = append(res, *p)
	}

	return res
}

func search2(pkg string, c chan *PackageResult, wg *sync.WaitGroup) {
	defer wg.Done()

	baseURL := "https://pkg.go.dev/search"
	fullURL := fmt.Sprintf("%s?q=%s&page=1", baseURL, pkg)

	resp, err := http.Get(fullURL)
	if err != nil {
		c <- &PackageResult{
			Package: pkg,
			Err:     err,
		}
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		c <- &PackageResult{
			Package: pkg,
			Err:     err,
		}
	}

	resultElems := find(doc, condHasClass(spClass))
	if len(resultElems) != 0 {
		if r := resultElems[0]; r != nil {
			rawRepo := r.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.Data
			repo := strings.TrimSpace(rawRepo)

			if isValidRepo(pkg, repo) {
				c <- &PackageResult{
					Package: pkg,
					Repo:    repo,
				}
			} else {
				c <- &PackageResult{
					Package: pkg,
					Err:     NoResultsErr(fmt.Errorf("no results for pkg '%s'", pkg)),
				}
			}
		}
	}
}

type NoResultsErr error

func isValidRepo(pkg, url string) bool {
	regex := regexp.MustCompile("(:?.*\\.)*(:?.*\\/)*ginkgo")
	isMatch := regex.MatchString(url)
	return len(url) > 0 &&
		strings.Count(url, "/") >= 2 &&
		isMatch &&
		strings.HasSuffix(url, "/"+pkg)
}

type PackageResult struct {
	Package string
	Repo    string
	Err     error
}

func find(node *html.Node, by cond) []*html.Node {
	nodes := make([]*html.Node, 0)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if by(c) {
			nodes = append(nodes, c)
		}
		nodes = append(nodes, find(c, by)...)
	}
	return nodes
}

type cond func(*html.Node) bool

func condHasClass(class string) cond {
	return func(node *html.Node) bool {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == class {
				return true
			}
		}
		return false
	}
}

func condValidTxt() cond {
	return func(node *html.Node) bool {
		return node.Type == html.TextNode && strings.TrimSpace(node.Data) != "" && node.Data != "|"
	}
}
