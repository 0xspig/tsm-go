package main

import (
	"html/template"
	"net/http"
	"os"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("home Function run")
	w.Header().Add("Server", "McServer")

	ts, err := template.ParseFiles("./ui/index.html")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *Application) getJSON(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("getJSON Function Run")
	w.Header().Add("Server", "McServer")
	json, err := os.ReadFile("ui/static/gen/graph-data.json")
	if err != nil {
		app.logger.Error("Error loading json", err)
		return
	}

	w.Write(json)
}

func (app *Application) getPostHTML(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("getPostHTML function run")
	w.Header().Add("Server", "McServer")
	html := app.garden.NodeToHTML(r.PathValue("id"))
	w.Write(html)
}
