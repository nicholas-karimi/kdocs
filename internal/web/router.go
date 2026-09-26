package web

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicholas-karimi/kdocs/internal/database"
	"github.com/nicholas-karimi/kdocs/internal/domain"
)

type Handler struct {
	db            *sql.DB
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

func NewRouter(db *sql.DB) (http.Handler, error) {
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
		db:            db,
		homeTemplate:  homeTemplate,
		pageTemplate:  pageTemplate,
		spaceTemplate: spaceTemplate,
	}

	router := chi.NewRouter()

	// load static files
	fileServer := http.FileServer(http.Dir("./static"))

	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

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

	spaces, err := database.FindSpaces(h.db)
	if err != nil {
		log.Println("Failed to find spaces:", nil)
		http.Error(w, "Unable to load spaces", http.StatusInternalServerError)
		return
	}
	data := HomeData{
		Title:       "KDocs",
		Description: "Internal Engineering Knowledge System",
		Spaces:      spaces,
	}
	err = h.homeTemplate.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "unable to render page", http.StatusInternalServerError)
	}
}

// space
func (h *Handler) space(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	space, found, err := database.FindSpace(h.db, slug)

	if err != nil {
		log.Println("Failed to find space:", err)
		http.Error(w, "Unabe to load space", http.StatusInternalServerError)
	}

	if !found {
		http.NotFound(w, r)
		return
	}
	pages, err := database.FindPagesBySpace(h.db, space.Slug)
	if err != nil {
		log.Println("Failed to find pages:", err)
		http.Error(w, "Unable to load pages", http.StatusInternalServerError)
		return
	}
	data := SpaceData{
		Space: space,
		Pages: pages,
	}

	err = h.spaceTemplate.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}

// page route

func (h *Handler) page(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, found, err := database.FindPage(h.db, slug)

	if err != nil {
		log.Println("Failed to find page:", err)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}

	if !found {
		http.NotFound(w, r)
		return
	}
	err = h.pageTemplate.ExecuteTemplate(w, "base", page)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
	}
}
