package goose

import (
	"log"
	"time"
)

func (engine *Engine) Logger() HandlerFunc {
	return func(ctx *Context) {
		// show param
		if engine.logRequestBody {
			log.Printf("Params: %v", ctx.QueryAll())
			if ctx.Method == "POST" || ctx.Method == "PATCH" || ctx.Method == "PUT" {
				log.Printf("Post Body: %v", ctx.PostFormAll())
			}
			log.Println("------")
		}

		// Start timer
		t := time.Now()

		// Process request
		ctx.Next()

		// response performance
		if engine.logPerformance {
			// Calculate resolution time
			log.Printf("Response In: %v", time.Since(t))
		}
	}
}
