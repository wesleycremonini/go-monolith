package main

import (
	"flag"
	"log"
	"time"
	"wesleycremonini/go-monolith-template/internal/database"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

type App struct {
	Logger         any
	DB             *database.DB
	SessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":80", "HTTP network address")
	dbDsn := flag.String("db-dsn", "", "DB DSN")
	redisHost := flag.String("redis-host", "", "Redis host")
	redisPass := flag.String("redis-pass", "", "Redis password")

	db, err := database.Connect(*dbDsn)
	if err != nil {
		log.Fatal("cant connect to db: ", err)
	}
	defer db.Close()

	db.Config().MaxConnIdleTime = 5 * time.Minute
	db.Config().MaxConnLifetime = 2 * time.Hour
	db.Config().MaxConns = 25
	db.Config().MinConns = 5
	db.Config().HealthCheckPeriod = 1 * time.Minute

	redis := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", *redisHost, redis.DialPassword(*redisPass))
		},
	}

	sesM := scs.New()
	sesM.Store = redisstore.New(redis)
	sesM.Lifetime = 24 * time.Hour

	app := &App{DB: db, SessionManager: sesM}

	err = app.serve(*addr)
	if err != nil {
		log.Fatal("cant start server: ", err)
	}
}
