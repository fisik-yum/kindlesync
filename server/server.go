package server

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pion/mdns"
	"go.etcd.io/bbolt"
	"golang.org/x/net/ipv4"
	"net"
)

type ksServer struct {
	port       uint16
	library    string
	mux        *http.ServeMux
	fsmux      http.Handler
	db         *bbolt.DB
	httpserver *http.Server
	mdnsserver *mdns.Config // should be a better idea to store the configuration
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

func (k *ksServer) Serve() {
	log.Println("Indexing library")
	k.genDB()
	log.Println("Index completed")
	log.Printf("Serving on http://localhost:%d", k.port)
	go http.ListenAndServe(":"+strconv.FormatUint(uint64(k.port), 10), k.mux)

	log.Printf("Advertising mDNS Service")
	log.Printf(mdns.DefaultAddress)
	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddress)
	if err != nil {
		panic(err)
	}

	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		panic(err)
	}

	_, err = mdns.Server(ipv4.NewPacketConn(l4), k.mdnsserver)
	if err != nil {
		return
	}
	select {}
}

func NewServer(port uint16, library string, db *bbolt.DB) *ksServer {
	serv := &ksServer{
		mux:        &http.ServeMux{},
		fsmux:      http.FileServer(http.Dir(filepath.Join(".", library))),
		httpserver: &http.Server{},
		port:       port,
		library:    library,
		db:         db,
	}
	serv.mux.HandleFunc("/refresh", serv.refreshHandler)
	serv.mux.HandleFunc("/sync", serv.syncHandler)
	serv.mux.Handle("/"+library+"/", http.StripPrefix("/"+library, serv.fsmux))
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	serv.mdnsserver = &mdns.Config{LocalNames: []string{host + "_kindlesync_.local"}}
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
			log.Printf("Indexing item %s", file.Name())
			err := b.Put([]byte(file.Name()), []byte("book"))
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}
