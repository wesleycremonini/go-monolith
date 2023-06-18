package main

import (
	"flag"
	"log"
)

type App struct {
	Logger         any
	DB             any
	SessionManager any
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")

	app := App{}

	err := app.server(*addr).ListenAndServe()
	if err != nil {
		log.Fatal("cant start server", err)
	}
}
