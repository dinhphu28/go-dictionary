package api

import (
	"fmt"
	"log"
	"net/http"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
)

type Router struct {
	dictionaries []database.Dictionary
	globalConfig config.GlobalConfig
}

func NewRouter(
	dictionaries []database.Dictionary,
	globalConfig config.GlobalConfig,
) *Router {
	return &Router{
		dictionaries: dictionaries,
		globalConfig: globalConfig,
	}
}

func (r *Router) StartAPIServer() {
	mux := http.NewServeMux()
	handler := CorsMiddleware(mux)

	mux.Handle("/lookup",
		LookupHandler(r.dictionaries, r.globalConfig),
	)

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
