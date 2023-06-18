package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"
)

type App struct {
	Logger         any
	DB             *sql.DB
	SessionManager any
}

func main() {
	addr := flag.String("addr", ":3000", "HTTP network address")
	dbDsn := flag.String("db_dsn", "db/db.db", "DB DSN")

	db, err := connectDB(*dbDsn)
	if err != nil {
		log.Fatal("cant start server", err)
	}
	defer db.Close()

	app := &App{DB: db}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("cant start server", err)
	}
}
