package main

import (
	"calenwu.com/snippetbox/pkg/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type ho struct {
	name string
}

func (h *ho) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strings.Join([]string{"This", h.name}, " ")))
}

func (app *application) home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", 404)
		return
	}
	files := []string{
		"./ui/html/home.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Fprintf(w, "%v", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	title := "O snail"
	content := "O snail\\nClimb Moun Fuji,\nBut slowly, slowly!\n\n- Kobayashi xyz"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

