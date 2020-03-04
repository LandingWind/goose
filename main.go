package main

import (
	// "fmt"
	"log"
	// "net/http"
	// own define
	goose "./goose"
)

func main() {
	engine:=goose.New()

	engine.GET("/", func(ctx *goose.Context) {
		ctx.Send("dddd hello main home", 200)
	})
	engine.GET("/hello", func(ctx *goose.Context) {
		ctx.Html("<h1>loooks bigger ,right?</h1>", 200)
	})

	log.Fatal(engine.BoostEngine(":9999"))
}