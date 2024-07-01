package main

import (
	//"errors"
	"io"
	"log"
	"net/http"
	"path/filepath"
	//"os"
)

func init() {
	log.Println("Indexing library")
	genDB()
	log.Println("Index completed")
}

func main() {
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ACTION: Refresh")
		genDB()
		log.Println("OK")
		io.WriteString(w, "OK")
	})

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ACTION: Catalog")
		http.ServeFile(w, r, "books")
		log.Println("OK")
	})

	lib := http.FileServer(http.Dir(filepath.Join(".", "library")))
	http.Handle("/library/", http.StripPrefix("/library", lib))

	http.ListenAndServe(":8080", nil)
}
