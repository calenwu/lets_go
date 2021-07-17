package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := chi.NewRouter()
	mux.Get("/", app.home)
	mux.Get("/snippet", app.showSnippet)
	mux.Get("/snippet/create", app.createSnippet)
	mux.Post("/snippet/create", app.createSnippet)
	mux.Get("/snippet/{snippetId}", app.showSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}