package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"go.etcd.io/bbolt"
)

type ksServer struct {
	port    uint16
	library string
	mux     *http.ServeMux
	fsmux   http.Handler
	db      *bbolt.DB
	server  *http.Server
}

func (k *ksServer) refreshHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ACTION: Refresh Catalog")
	k.genDB()
	log.Println("OK")
	io.WriteString(w, "OK")
}

func (k *ksServer) syncHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ACTION: Retrieve Catalog")
	http.ServeFile(w, r, "books.db")
	log.Println("OK")
}

func (k *ksServer) serve() {
	log.Printf("Serving on http://localhost:%d",k.port)
	http.ListenAndServe(":"+strconv.FormatUint(uint64(k.port), 10), k.mux)
}

func newServer(port uint16, library string, db *bbolt.DB) *ksServer {
	serv := &ksServer{
		mux:    &http.ServeMux{},
		fsmux:  http.FileServer(http.Dir(filepath.Join(".", library))),
		server: &http.Server{},
		port:   port,
		library: library,
		db:     db,
	}
	serv.mux.HandleFunc("/refresh", serv.refreshHandler)
	serv.mux.HandleFunc("/sync", serv.syncHandler)
	serv.mux.Handle("/"+library+"/", http.StripPrefix("/"+library, serv.fsmux))
	return serv
}

func (k *ksServer) genDB() error {
	files, err := os.ReadDir(filepath.Join(".", k.library))
	if err != nil {
		return err
	}
	tx, err := k.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Commit()
	b, err := tx.CreateBucketIfNotExists([]byte("books"))
	if err != nil {
		return err
	}

	for _, file := range files {
		// TODO: make file extensions configurable
		if !file.IsDir() {
			log.Printf("Indexing item %s",file.Name())
			err := b.Put([]byte(file.Name()), []byte("book"))
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}
