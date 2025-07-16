package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("POST /snippet/create", app.createSnippet)
	mux.HandleFunc("GET /snippet/view", app.viewSnippet)

	return mux
}
