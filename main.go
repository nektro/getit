package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	mapset "github.com/deckarep/golang-set"
	"github.com/nektro/go-util/util"
	"github.com/schollz/progressbar"
)

var (
	root  string
	links = mapset.NewSet()
)

func main() {
	val := flag.String("url", "", "The URL to search")
	flag.Parse()

	if len(*val) == 0 {
		fmt.Println(getTime(), "-url is required!")
		return
	}

	u, _ := url.Parse(*val)
	u.Path = filepath.Dir(u.Path)
	root = u.String()

	fmt.Println(getTime(), "Input:", *val)
	fmt.Println(getTime(), "Root: ", root)

	crawlPage(*val)
	fmt.Println()
}

func crawlPage(path string) {
	fmt.Println(getTime(), "[DIR]", path)
	//
	res, err := http.Get(path)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return
	}

	//
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	//
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		//
		v, o := s.Attr("href")
		if !o {
			return
		}
		if strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "http://") {
			return
		}
		if strings.HasPrefix(v, "#") {
			return
		}
		if strings.HasSuffix(v, "index.html") {
			v = v[0 : len(v)-10]
		}
		if strings.HasSuffix(v, "index.htm") {
			v = v[0 : len(v)-9]
		}
		n, _ := url.PathUnescape(v)
		u, _ := url.Parse(path)
		if strings.HasPrefix(n, "/") {
			u.Path = n
		} else {
			u.Path = filepath.Join(u.Path, n)
		}
		p := u.String()

		if !strings.HasPrefix(p, root) {
			return
		}
		if !links.Add(p) {
			return
		}

		if !hasFileExtension(p) {
			crawlPage(p)
			return
		}
		download(p)
	})
}

func hasFileExtension(page string) bool {
	a := strings.Split(page, "/")
	b := a[len(a)-1]
	return strings.Contains(b, ".")
}

const (
	moveToStartOfLine = "\x1b[0G"
)

func download(path string) {
	u, _ := url.Parse(path)
	n, _ := url.PathUnescape(u.Path)
	p, _ := filepath.Abs(u.Hostname() + n)

	if util.DoesFileExist(p) {
		return
	}
	fmt.Println(getTime(), p)

	res, _ := http.Get(path)
	defer res.Body.Close()

	d := filepath.Dir(p)
	os.MkdirAll(d, os.ModePerm)
	f, _ := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	defer f.Close()

	l := int(res.ContentLength)
	bar := progressbar.NewOptions(l, progressbar.OptionSetBytes(l))
	out := io.MultiWriter(f, bar)
	io.Copy(out, res.Body)

	fmt.Print(moveToStartOfLine)
}

func getTime() string {
	return "[" + time.Now().String()[11:19] + "]"
}
