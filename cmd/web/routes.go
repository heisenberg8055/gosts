package main

import "net/http"

func (app *application) routes() http.Handler {

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /snippet/create", app.createSnippet)
	mux.HandleFunc("POST /snippet/create", app.createSnippetPost)
	mux.HandleFunc("GET /snippet/view/{id}", app.viewSnippet)

	stack := CreateStack(app.recoverPanic, app.logRequest, secureHeaders)

	return stack(mux)
}
