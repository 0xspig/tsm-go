package main

import (
	"html/template"
	"net/http"
	"path/filepath"
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
	json, err := app.garden.ExportJSONData()
	if err != nil {
		app.logger.Error("json export error", err)
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

func (app *Application) resolveID(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("resolving ID")
	w.Header().Add("Server", "McServer")
	request_ID := filepath.Base(r.RequestURI)
	if app.garden.ContainsID(request_ID) {

		ts, err := template.ParseFiles("./ui/index.html")
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		//onload_string := "onload=\"targetNode('" + request_ID + "\"')"
		err = ts.Execute(w, map[string]string{"OnLoad": request_ID})
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else {
		app.home(w, r)
		return
	}
}
