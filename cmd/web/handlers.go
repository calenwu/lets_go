package main

import (
	"calenwu.com/snippetbox/pkg/forms"
	"calenwu.com/snippetbox/pkg/models"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"net/url"
	"strconv"
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
	type templateData struct {
		Snippets []*models.Snippet
	}
	app.render(w, r, "home.page.gohtml", &templateData{s})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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
	app.render(w, r, "show.page.gohtml", &templateData{s})
}

type CreateSnippetTemplateData struct {
	FormErrors map[string]string
	FormData   url.Values
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// http.MaxBytesReader()
	// r.PostForm["title"] does the same thing
	form := forms.New(r.PostForm)
	form.Required([]string{"title", "content", "expires"})
	form.MaxLength([]string{"title"}, 100)
	form.PermittedValues([]string{"expires"}, "1", "7", "365")
	if !form.Valid() {
		app.render(w, r, "create.page.gohtml", form)
		return
	}
	expires := r.PostForm.Get("expires")
	expiresInt, _ := strconv.Atoi(expires)
	id, err := app.snippets.Insert(
		r.PostForm.Get("title"),
		r.PostForm.Get("content"),
		expiresInt)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w  http.ResponseWriter, r *http.Request) {
	//Form: forms.New(nil),
	app.render(w, r, "create.page.gohtml", forms.New(nil))
}