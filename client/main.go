package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
	"gopkg.in/ini.v1"
)

var (
	URL    string
	libdir string
)

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Insufficient args")
	}
	cfg, err := ini.Load(filepath.Join(".", "kindlesync.ini"))
	check(err)
	URL = "http://" + cfg.Section("").Key("address").String() + ":" + cfg.Section("").Key("port").String()
	libdir = cfg.Section("").Key("library").String() //Always the absolute path (/mnt/us/path/to/dir)
}

func main() {
	switch os.Args[1] {
	case "refresh": // force remote to refresh
		r, e := http.Get(URL + "/refresh")
		check(e)
		s, e := io.ReadAll(r.Body)
		check(e)
		log.Println("Refresh Catalog: " + string(s))
	case "sync-retrieve": // download database and download all books to library
		err := downloadFile(filepath.Join(".", "remote.db"), URL+"/sync")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Retrieve Catalog: OK")
		log.Println("Populating Library")

		db,err:=bbolt.Open("remote.db",0600,nil)
		if err != nil {
			log.Fatal(err)
		}
		tx,_:=db.Begin(false)
		tx.Bucket([]byte("books")).ForEach(func(k, v []byte) error {
			err:=downloadFile(filepath.Join(libdir,string(k)),fmt.Sprintf("%s/library/%s",URL,k))
			return err
		})
		tx.Rollback()
		db.Close()
	case "sync": // plain sync - only retrieve the database
		err := downloadFile(filepath.Join(".", "remote.db"), URL+"/sync")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Retrieve Catalog: OK")
	default:
		log.Println("Invalid command!")
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func downloadFile(filepath string, url string) (err error) {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
