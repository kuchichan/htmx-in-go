package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, there!</h1>")
}

func main() {
	http.HandleFunc("/hello", hello)

	server := &http.Server{
		Addr: ":4000",
	}
	log.Fatal(server.ListenAndServe())
}
