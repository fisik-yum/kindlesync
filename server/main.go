package main

import (
	"log"

	"go.etcd.io/bbolt"
)



func init() {
}

func main() {
	//serv := &http.Server{}
	db,err:=bbolt.Open("books.db",0600,nil)
	if err!=nil{
		log.Fatal(err)
	}
	serv := newServer(8080, "library",db)
	log.Println("Indexing library")
	serv.genDB()
	log.Println("Index completed")
	serv.serve()
}
