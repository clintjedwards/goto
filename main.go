package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func test(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("test"))
}

func main() {
	router := mux.NewRouter()

	router.Handle("/test", handlers.MethodHandler{
		"GET": http.HandlerFunc(test),
	})

	server := http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", server.Handler))
}
