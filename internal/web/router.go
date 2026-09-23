package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	templates *template.Template
}

func NewRouter() (http.Handler, error) {
	templates, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/home.html",
	)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		templates: templates,
	}

	router := chi.NewRouter()

	/* using closure
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "base", nil)

		if err != nil {
			http.Error(w, "Unable to render page", http.StatusInternalServerError)
			return
		}
	}) */
	router.Get("/", handler.home)

	// route group
	router.Route("/pages", func(r chi.Router) {
		r.Get("/{slug}", page)
	})

	return router, nil
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "unable to render page", http.StatusInternalServerError)
	}
}

/* func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/hom.html",
	)
	if err != nil {
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
} */

func page(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	fmt.Fprintf(w, "KDocs page: %s", slug)
}
