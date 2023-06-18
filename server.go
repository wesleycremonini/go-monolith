package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/alexedwards/flow"
)

//go:embed public/*
var assetsFS embed.FS

func (a *App) routes() *flow.Mux {
	mux := flow.New()
	
	mux.Use(a.recoverPanic)
	
	mux.Handle("/public/...", http.FileServer(http.FS(assetsFS)), "GET")
	
	mux.HandleFunc("/status", a.status, "GET")
	mux.HandleFunc("/", a.home, "GET")
	mux.HandleFunc("/", a.newItem, "POST")

	return mux
}

func (a *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovered from panic: ", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *App) status(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{
		"status": "ok",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
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

func (a *App) newItem(w http.ResponseWriter, r *http.Request) {
	X := r.PostFormValue("X")
	Y := r.PostFormValue("Y")
	tmpl := template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html"))
	tmpl.ExecuteTemplate(w, "items-element", Item{X: X, Y: Y})
}
