package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(
		app.recoverPanic,
		app.logRequest,
		app.sessionMiddleware,
		secureHeaders)
	// Create a new middleware chain containing the middleware specific to
	// out dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.

	mux := chi.NewRouter()
	mux.Get("/", app.home)
	mux.Get("/snippet", app.showSnippet)
	mux.With(app.requireAuthenticatedUser).Get("/snippet/create", app.createSnippetForm)
	mux.With(app.requireAuthenticatedUser).Post("/snippet/create", app.createSnippet)
	mux.Get("/snippet/{id}", app.showSnippet)

	// User
	mux.Get("/user/signup", app.signupUserForm)
	mux.Post("/user/signup", app.signupUser)
	mux.Get("/user/login", app.loginUserForm)
	mux.Post("/user/login", app.loginUser)
	mux.Post("/user/logout", app.logoutUser)

	FileServer(mux, "/static", http.Dir("./ui/static/"))
	return standardMiddleware.Then(mux)
}


// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}