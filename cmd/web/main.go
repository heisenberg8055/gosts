package main

import (
	"log"
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func main() {

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /", home)
	mux.HandleFunc("POST /snippet/create", createSnippet)
	mux.HandleFunc("POST /snippet/view", viewSnippet)
	log.Println("Starting server on: 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
