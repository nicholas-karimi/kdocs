package web

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicholas-karimi/kdocs/internal/domain"
)

type Handler struct {
	homeTemplate  *template.Template
	spaceTemplate *template.Template
	pageTemplate  *template.Template
}

// view/application data — data specifically needed to render the homepage.
type HomeData struct {
	Title       string
	Description string
	Spaces      []domain.Space
}

type SpaceData struct {
	Space domain.Space
	Pages []domain.Page
}

func NewRouter() (http.Handler, error) {
	homeTemplate, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/home.html",
	)
	if err != nil {
		return nil, err
	}
	spaceTemplate, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/space.html",
	)
	if err != nil {
		return nil, err
	}
	pageTemplate, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/page.html",
	)
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		homeTemplate:  homeTemplate,
		pageTemplate:  pageTemplate,
		spaceTemplate: spaceTemplate,
	}

	router := chi.NewRouter()

	// load static files
	fileServer := http.FileServer(http.Dir("./static"))

	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

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
		r.Get("/{slug}", handler.page)
	})

	//
	router.Route("/spaces", func(r chi.Router) {
		r.Get("/{slug}", handler.space)
	})

	return router, nil
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {

	data := HomeData{
		Title:       "KDocs",
		Description: "Internal Engineering Knowledge System",
		Spaces:      domain.SampleSpaces(),
	}
	err := h.homeTemplate.ExecuteTemplate(w, "base", data)
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

func (h *Handler) page(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, found := domain.FindPage(slug)

	if !found {
		http.NotFound(w, r)
		return
	}
	err := h.pageTemplate.ExecuteTemplate(w, "base", page)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}

func (h *Handler) space(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	space, found := domain.FindSpace(slug)

	if !found {
		http.NotFound(w, r)
		return
	}

	data := SpaceData{
		Space: space,
		Pages: domain.FindPagesBySpace(space.Slug),
	}
	/* fmt.Fprintln(w, "Space:", data.Space.Name)

	for _, page := range data.Pages {
		fmt.Fprintln(w, "Page:", page.Title)
	} */
	err := h.spaceTemplate.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}
