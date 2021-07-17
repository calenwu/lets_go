package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"strconv"

	"calenwu.com/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", 404)
		return
	}
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
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
	err = ts.Execute(w, s)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "snippetId"))
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
	type templateData struct {
		Snippet *models.Snippet
	}
	files := []string {
		"./ui/html/show.page.gohtml",
		"./ui/html/base.layout.gohtml",
		"./ui/html/footer.partial.gohtml",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, &templateData{s})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
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
