package api

import (
	"fmt"
	"log"
	"net/http"
)

type Router struct {
	lookupHandler LookupHandler
}

func NewRouter(lookupHandler LookupHandler) *Router {
	return &Router{
		lookupHandler: lookupHandler,
	}
}

func (r *Router) StartAPIServer() {
	mux := http.NewServeMux()
	handler := CorsMiddleware(mux)

	mux.Handle("/lookup", &r.lookupHandler)

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
