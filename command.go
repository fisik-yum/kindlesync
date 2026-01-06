package main

import (
	"kindlesync/client"
	"kindlesync/server"
	"log"
	"os"

	"go.etcd.io/bbolt"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments!")
	}
	cmd := os.Args[1]
	opt := os.Args[2]
	switch cmd {
	case "client":
		client.Execute(opt)
	case "server":
		db, err := bbolt.Open("books.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		serv := server.NewServer(8080, opt, db)
		serv.Serve()
	default:
		log.Fatal("not a valid command")
	}

}
