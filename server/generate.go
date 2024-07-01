package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func genDB() {
	files, err := os.ReadDir(filepath.Join(".", "library"))
	check(err)
	fd, err := os.Create("books")
	check(err)
	defer fd.Close()
	for _, j := range files {
		if !j.IsDir() && (strings.HasSuffix(j.Name(), ".mobi") || strings.HasSuffix(j.Name(), ".azw3")) {
			fmt.Fprintln(fd, j.Name())
		}
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
