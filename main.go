package main

import (
	"fmt"
	"log"
	"net/http"
	// own define
	goose "./goose"
)

func main() {
	engine:=goose.New()

	engine.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is Home page!")
	})
	engine.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World page!")
	})

	log.Fatal(engine.BoostEngine(":9999"))
}