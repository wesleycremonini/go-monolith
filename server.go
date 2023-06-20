package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/flow"
	"golang.org/x/exp/slog"
)

//go:embed public/*
var assetsFS embed.FS

// MIDDELWARES ##########

func (app *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.Log.Warn("Recovered from panic: ", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *App) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// ROUTES ###########

func (app *App) routes() *flow.Mux {
	mux := flow.New()

	mux.Use(app.recoverPanic)
	mux.Use(app.secureHeaders)
	mux.Use(app.SessionManager.LoadAndSave)

	mux.Handle("/public/...", http.FileServer(http.FS(assetsFS)), "GET")

	// GET	/user/signup	userSignup	Display a HTML form for signing up a new user
	// POST	/user/signup	userSignupPost	Create a new user
	// GET	/user/login	userLogin	Display a HTML form for logging in a user
	// POST	/user/login	userLoginPost	Authenticate and login the user

	mux.HandleFunc("/status", app.status, "GET")
	mux.HandleFunc("/", app.home, "GET")
	mux.HandleFunc("/", app.newItem, "POST")

	return mux
}

func (app *App) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for signing up a new user...")
}

func (app *App) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *App) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *App) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *App) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}

func (app *App) status(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{
		"status": "ok",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	items := map[string][]Item{
		"Items": {
			{X: "this is X", Y: "this is Y"},
			{X: "this is X2", Y: "this is Y2"},
			{X: "this is X3", Y: "this is Y3"},
		},
	}

	tmpl := template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html"))

	err := tmpl.ExecuteTemplate(w, "base", items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) newItem(w http.ResponseWriter, r *http.Request) {
	X := r.PostFormValue("X")
	Y := r.PostFormValue("Y")
	tmpl := template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html"))
	tmpl.ExecuteTemplate(w, "items-element", Item{X: X, Y: Y})
}

// SERVER #########

func (app *App) serve(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.Log.Warn("shutting down server", slog.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.Log.Info("server started")
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.Log.Warn("stopped server")
	return nil
}
