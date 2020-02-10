package main

import (
	"log"
	"net/http"
	"time"

	"github.com/clintjedwards/go/config"
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
	router := mux.NewRouter()

	router.Handle("/create", handlers.MethodHandler{
		"POST": http.HandlerFunc(listLinks),
	})

	// router.Handle("/links", handlers.MethodHandler{
	// 	"GET": http.HandlerFunc(listLinks),
	// })

	// router.Handle("/links", handlers.MethodHandler{
	// 	"GET": http.HandlerFunc(listLinks),
	// })

	server := http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", server.Handler))
}
