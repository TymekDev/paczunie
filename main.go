package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"git.sr.ht/~tymek/rpi-paczunie/packages"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	dbName := flag.String("database", "packages.db", "path to SQLite3 database")
	initIfEmpty := flag.Bool("init", false, "initialize Packages table if it does not exist")
	flag.Parse()

	c, err := packages.NewClientWithSQLiteStorage(*dbName, *initIfEmpty)
	if err != nil {
		log.Fatalln("FATAL", err)
	}

	log.Println("INFO", "started listening on", *port)
	if err := http.ListenAndServe(fmt.Sprint(":", *port), c); err != nil {
		log.Fatalln("FATAL", err)
	}
}
