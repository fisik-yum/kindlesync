package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	case "refresh":
		r, e := http.Get(URL + "/refresh")
		check(e)
		s, e := io.ReadAll(r.Body)
		check(e)
		log.Println("Refresh: " + string(s))
	case "sync":
		r, e := http.Get(URL + "/books")
		check(e)
		defer r.Body.Close()
		scanner := bufio.NewScanner(r.Body)
		for scanner.Scan() {
			s := scanner.Text()
			_, e := os.Stat(filepath.Join(libdir, s))
			if e != nil {
				downloadFile(filepath.Join(libdir, s), URL+"/library/"+s)
				log.Printf("Download completed: %s\n", s)
			}
		}
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
