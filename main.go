package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"flag"

	"strings"

	"github.com/withnic/gowiki/models/pages"
	"github.com/withnic/gowiki/models/templates"
)

var renderer *templates.Template

func init() {
	renderer = templates.Util.New()
}

var validPath = regexp.MustCompile(`^/(edit|save|view)/([a-zA-Z0-9]+)$`)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
	return
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		err := renderer.Render(w, "errorPage", "404 Page Not found.")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if status == http.StatusInternalServerError {
		err := renderer.Render(w, "errorPage", "Internal Server Error")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		fn(w, r, m[2])
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	if strings.HasPrefix(r.URL.Path, "/static/") && len(strings.TrimPrefix(r.URL.Path, "/static/")) > 0 {
		fileServer.ServeHTTP(w, r)
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pages.Util.Load(title)

	if err != nil {
		p = &pages.Page{
			Title: title,
		}
	}

	err = renderer.Render(w, "editPage", p)
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pages.Util.Load(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderer.Render(w, "viewPage", p)
}

func markdownHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	body := r.FormValue("body")
	output := pages.Util.Parse([]byte(body))

	if err := json.NewEncoder(w).Encode(string(output)); err != nil {
		log.Fatal(err)
	}

}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &pages.Page{
		Title: title,
		Body:  []byte(body),
	}
	err := p.Save()
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var (
	port = flag.Int("port", 8000, "port number")
)

func main() {
	flag.Parse()

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/api/mkd", markdownHandler)
	http.HandleFunc("/", handler)

	addr := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal(err)
	}
}
