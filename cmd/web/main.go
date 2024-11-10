package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type Application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", "localhost:4000", "HTTP Network Address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &Application{
		logger: logger,
	}
	logger.Info("Starting server on", "addr", *addr)
	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}
