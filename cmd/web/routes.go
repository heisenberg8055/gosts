package main

import "net/http"

func (app *application) routes() http.Handler {

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.Handle("GET /", app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	mux.Handle("GET /snippet/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.createSnippet)))
	mux.Handle("POST /snippet/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.createSnippetPost)))
	mux.Handle("GET /snippet/view/{id}", app.sessionManager.LoadAndSave(http.HandlerFunc(app.viewSnippet)))
	mux.Handle("GET /user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignup)))
	mux.Handle("POST /user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignupPost)))
	mux.Handle("GET /user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogin)))
	mux.Handle("POST /user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLoginPost)))
	mux.Handle("POST /user/logout", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogoutPost)))

	stack := CreateStack(app.recoverPanic, app.logRequest, secureHeaders)

	return stack(mux)
}
