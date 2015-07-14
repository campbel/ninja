package main

import (
	"fmt"
	"net/http"

	"github.com/campbel/ninja"
)

func main() {
	api := ninja.New()

	api.Route("/echo", ninja.Route{
		"GET":  get,
		"POST": post,
	}, always)

	http.ListenAndServe(":8080", api)
}

func always(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "You used: ")
	return nil
}

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET")
}

func post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "POST")
}
