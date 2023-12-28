package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/class/add", addClass)
	mux.HandleFunc("/class/view", viewClass)
	log.Print("starting server on :5000")

	err := http.ListenAndServe(":5000", mux)
	log.Fatal(err)
}
