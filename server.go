package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"urlShortener/database"
	"urlShortener/handler"
)

func main() {

	db := database.MustConnect("db.db")
	router := chi.NewRouter()

	router.Post("/create", handler.PostShorthand(db))
	router.Get("/*", handler.GetURL(db))

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {

		text := "Method not allowed:\n" +
			"use POST /create \"url\":\"google.com\" to create a shortened url " +
			"and GET /{shorthand} to retrieve a shortened url target"

		http.Error(w, text, http.StatusNotFound)
	})

	log.Fatal(http.ListenAndServe(":4000", router))
}
