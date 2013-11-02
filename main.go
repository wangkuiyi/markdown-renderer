package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/knieriem/markdown"
	"log"
	"net/http"
	"regexp"
)

var addr = flag.String("addr", ":8002", "address to listen on")
var data = flag.String("data", "http://localhost:8003/", "the data server")
var css  = flag.String("css", "/markdown.css", "path to CSS")

var validPath = regexp.MustCompile("^/([_a-zA-Z0-9]+)\\.md$")

func renderMarkdownHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.Error(w, r.URL.Path, http.StatusNotAcceptable)
		return
	}
	md := *data + m[1] + ".md"

	resp, err := http.Get(md)
	if err != nil {
		http.Error(w, err.Error() + md, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	bo := bufio.NewWriter(w)
	fmt.Fprintf(bo, "<link href=\"%s\" rel=\"stylesheet\"> </link>\n", *css)
	markdown.NewParser(nil).Markdown(resp.Body, markdown.ToHTML(bo))
	bo.Flush()
}

func main() {
	flag.Parse()
	http.HandleFunc("/", renderMarkdownHandler)
	e := http.ListenAndServe(*addr, nil)
	if e != nil {
		log.Fatal("ListenAndServe: ", e)
	}
}

