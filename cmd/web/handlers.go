package main

import (
	"calenwu.com/snippetbox/pkg/forms"
	"calenwu.com/snippetbox/pkg/models"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
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

	session, err := app.session.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.AddFlash("Your snippet has been created")
	err = session.Save(r, w)
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

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.gohtml", forms.New(nil))
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	session, _ := app.session.Get(r, "session-name")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required([]string{"name", "email", "password"})
	form.IsEmail([]string{"email"})
	form.MinLength([]string{"password"}, 8)

	if !form.Valid() {
		app.render(w, r, "signup.page.gohtml", &form)
		return
	}

	err = app.users.Insert(
		form.Get("name"),
		form.Get("email"),
		form.Get("password"),
	)

	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Email is already in use")
		app.render(w, r, "signup.page.gohtml", &form)
	} else if err != nil {
		app.serverError(w, err)
	}
	session.AddFlash("Your signup was successful, Please log in.")
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", forms.New(nil))
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	session, _ := app.session.Get(r, "session-name")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or password is incorrect")
		app.render(w, r, "login.page.gohtml", form)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	session.Values["userID"] = id
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
	}
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := app.session.Get(r, "session-name")
	delete(session.Values, "userID")
	session.AddFlash("You have been logged out successfully!")
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
