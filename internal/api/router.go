package api

import (
	"fmt"
	"log"
	"net/http"
)

type Router struct {
	lookupHandler   LookupHandler
	lookupHandlerV2 LookupHandlerV2
}

func NewRouter(
	lookupHandler LookupHandler,
	lookupHandlerV2 LookupHandlerV2,
) *Router {
	return &Router{
		lookupHandler:   lookupHandler,
		lookupHandlerV2: lookupHandlerV2,
	}
}

func (r *Router) StartAPIServer() {
	mux := http.NewServeMux()
	handler := CorsMiddleware(mux)

	mux.Handle("/lookup", &r.lookupHandler)
	mux.Handle("/v2/lookup", &r.lookupHandlerV2)

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
