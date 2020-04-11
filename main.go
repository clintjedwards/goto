package main

import (
	"net/http"
	"time"

	"github.com/clintjedwards/goto/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {

	config, err := config.FromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load env config")
	}

	setupLogging(config.LogLevel, config.Debug)

	app := newApp()
	router := mux.NewRouter()

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

	log.Fatal().Err(http.ListenAndServe(config.Host, server.Handler))
}
