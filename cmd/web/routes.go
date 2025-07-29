package main

import (
	"net/http"

	"github.com/heisenberg8055/gosts/ui"
)

func (app *application) routes() http.Handler {

	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	mux := http.NewServeMux()
	mux.Handle("GET /static/", fileServer)
	mux.Handle("GET /healthz", http.HandlerFunc(healthCheck))

	noAuthHandler := New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	authHandler := noAuthHandler.Append(app.requireAuthentication)
	mux.Handle("GET /", noAuthHandler.ThenFunc(app.home))
	mux.Handle("GET /snippet/create", authHandler.ThenFunc(app.createSnippet))
	mux.Handle("POST /snippet/create", authHandler.ThenFunc(app.createSnippetPost))
	mux.Handle("GET /snippet/view/{id}", noAuthHandler.ThenFunc(app.viewSnippet))
	mux.Handle("GET /user/signup", noAuthHandler.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", noAuthHandler.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", noAuthHandler.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", noAuthHandler.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", authHandler.ThenFunc(app.userLogoutPost))

	stack := New(app.recoverPanic, app.logRequest, secureHeaders)

	return stack.Then(mux)
}
