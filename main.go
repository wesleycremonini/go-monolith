package main

import (
	"flag"
	"html/template"
	"os"
	"time"
	"wesleycremonini/go-monolith-template/internal/database"
	"wesleycremonini/go-monolith-template/internal/models"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/exp/slog"
)

type App struct {
	log            *slog.Logger
	db             *database.DB
	sessionManager *scs.SessionManager
	redis          *redis.Pool
	templateCache  map[string]*template.Template

	users *models.Users
}

func main() {
	addr := flag.String("addr", os.Getenv("ADDR"), "HTTP network address")
	dbDsn := flag.String("db-dsn", os.Getenv("DB_DSN"), "DB DSN")
	redisHost := flag.String("redis-host", os.Getenv("REDIS_HOST"), "Redis host")
	redisPass := flag.String("redis-pass", os.Getenv("REDIS_PASS"), "Redis password")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, err := database.Connect(*dbDsn)
	if err != nil {
		logger.Error("cant connect to db: " + err.Error())
		return
	}
	defer db.Close()

	db.Config().MaxConnIdleTime = 5 * time.Minute
	db.Config().MaxConnLifetime = 2 * time.Hour
	db.Config().MaxConns = 25
	db.Config().MinConns = 5
	db.Config().HealthCheckPeriod = 1 * time.Minute

	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", *redisHost, redis.DialPassword(*redisPass))
		},
	}
	defer redisPool.Close()

	sesM := scs.New()
	sesM.Store = redisstore.New(redisPool)
	sesM.Lifetime = 24 * time.Hour

	tmplCache, err := newTemplateCache()
	if err != nil {
		logger.Error("cant cache templates: " + err.Error())
		return
	}

	app := &App{
		db:             db,
		redis:          redisPool,
		sessionManager: sesM,
		log:            logger,
		templateCache:  tmplCache,

		users: &models.Users{Pool: db.Pool},
	}

	err = app.serve(*addr)
	if err != nil {
		logger.Error("cant start server: " + err.Error())
	}
}