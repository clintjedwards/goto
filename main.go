package main

import (
	"log"
	"net/http"
	"time"

	"github.com/clintjedwards/goto/config"
	"github.com/clintjedwards/toolkit/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func init() {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	logger.InitGlobalLogger(config.LogLevel, config.Debug)
}

func main() {
	app := newApp()
	router := mux.NewRouter()

	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	router.Handle("/links", handlers.MethodHandler{
		"GET": http.HandlerFunc(app.listLinksHandler),
	})

	router.Handle("/links/{id}", handlers.MethodHandler{
		"GET":    http.HandlerFunc(app.getLinkHandler),
		"DELETE": http.HandlerFunc(app.deleteLinksHandler),
	})

	router.Handle("/create", handlers.MethodHandler{
		"POST": http.HandlerFunc(app.createLinkHandler),
	})

	router.Handle("/{id}", handlers.MethodHandler{
		"GET": http.HandlerFunc(app.followLinkHandler),
	})

	server := http.Server{
		Addr:         config.Host,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(http.ListenAndServe(config.Host, server.Handler))
}
