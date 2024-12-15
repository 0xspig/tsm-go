package main

import (
	"html/template"
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("home Function run")
	w.Header().Add("Server", "McServer")

	ts, err := template.ParseFiles("./ui/vite/index.html")
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

func (app *Application) getNode(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("getNode Function Run")
	w.Header().Add("Server", "McServer")

}
