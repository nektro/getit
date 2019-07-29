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
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
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
	eraseLine         = "\x1b[2K"
	moveToStartOfLine = "\x1b[0G"
	moveUp            = "\x1b[A"
)

func download(path string) {
	u, _ := url.Parse(path)
	n, _ := url.PathUnescape(u.Path)
	p, _ := filepath.Abs(strings.ReplaceAll(u.Hostname(), ":", "+") + n)
	fmt.Println(getTime(), p)

	res, _ := http.Get(path)
	defer res.Body.Close()

	if util.DoesFileExist(p) {
		info, err := os.Stat(p)
		if err == nil {
			if info.Size() == res.ContentLength {
				return
			}
		}
	}

	d := filepath.Dir(p)
	os.MkdirAll(d, os.ModePerm)
	f, _ := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	defer f.Close()

	l := res.ContentLength

	q := mpb.New(
		mpb.WithWidth(80),
		mpb.WithRefreshRate(150*time.Millisecond),
	)
	bar := q.AddBar(l, mpb.BarStyle("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% 6.1f / % 6.1f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_MMSS, float64(l)/2048),
			decor.Name(" ] "),
			decor.AverageSpeed(decor.UnitKiB, "% .2f"),
		),
	)

	reader := bar.ProxyReader(res.Body)
	io.Copy(f, reader)
	q.Wait()

	fmt.Print(moveUp)
	fmt.Print(eraseLine)
}

func getTime() string {
	return "[" + time.Now().String()[11:19] + "]"
}
