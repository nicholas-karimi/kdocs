package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/", home)

	// route group
	router.Route("/pages", func(r chi.Router) {
		r.Get("/{slug}", page)
	})

	return router
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to KDocs")
}

func page(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	fmt.Fprintf(w, "KDocs page: %s", slug)
}
