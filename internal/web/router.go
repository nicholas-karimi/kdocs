package web

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicholas-karimi/kdocs/internal/database"
	"github.com/nicholas-karimi/kdocs/internal/domain"
	"github.com/nicholas-karimi/kdocs/internal/markdown"
)

type Handler struct {
	db               *sql.DB
	homeTemplate     *template.Template
	spaceTemplate    *template.Template
	pageTemplate     *template.Template
	newPageTemplate  *template.Template
	editPageTemplate *template.Template

	markdownRenderer *markdown.Renderer
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

type PageData struct {
	Page    domain.Page
	Content template.HTML
	Title   string
}

func NewRouter(db *sql.DB) (http.Handler, error) {

	markdownRenderer := markdown.NewRenderer()

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

	// add new page

	newPageTemplate, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/new.html",
	)
	if err != nil {
		return nil, err
	}

	editPageTemplate, err := template.ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/edit.html",
	)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		db:               db,
		homeTemplate:     homeTemplate,
		pageTemplate:     pageTemplate,
		spaceTemplate:    spaceTemplate,
		newPageTemplate:  newPageTemplate,
		editPageTemplate: editPageTemplate,

		markdownRenderer: markdownRenderer,
	}

	router := chi.NewRouter()

	// load static files
	fileServer := http.FileServer(http.Dir("./static"))

	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	router.Get("/", h.home)

	// route group
	router.Route("/pages", func(r chi.Router) {
		r.Get("/{slug}", h.page)
	})

	// spaces
	router.Route("/spaces", func(r chi.Router) {
		r.Get("/{slug}", h.space)
	})

	// new page
	router.Get("/pages/new", h.newPage)

	//create page
	router.Post("/pages", h.createPage)

	router.Get("/pages/{slug}/edit", h.editPage)
	router.Post("/pages/{slug}/edit", h.updatePage)

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

	renderedContent, err := h.markdownRenderer.Render(page.Content)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		return
	}
	log.Printf("RENDERED HTML:\n%s", renderedContent)
	data := PageData{
		Page:    page,
		Content: template.HTML(renderedContent),
		Title:   page.Title,
	}
	err = h.pageTemplate.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	http.Error(w, "Unable to render page", http.StatusInternalServerError)
	// }
	if err != nil {
		log.Println("Failed to execute page template:", err)
		return
	}
}

// create new page
func (h *Handler) newPage(w http.ResponseWriter, r *http.Request) {
	spaces, err := database.FindSpaces(h.db)
	if err != nil {
		log.Println("Failed to find spaces:", err)
		http.Error(w, "Unable to load page form", http.StatusInternalServerError)
		return
	}

	// anonymous newpagestruct-since data structure only resides here
	data := struct {
		Title  string
		Spaces []domain.Space
	}{
		Title:  "Create Page",
		Spaces: spaces,
	}

	if err := h.newPageTemplate.ExecuteTemplate(w, "base", data); err != nil {
		log.Println("Failed to execute new page template:", err)
	}
}

func (h *Handler) createPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	slug := r.FormValue("slug")
	spaceSlug := r.FormValue("space_slug")
	content := r.FormValue("content")

	if title == "" || slug == "" || spaceSlug == "" || content == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	log.Printf("CONTENT RECEIVED:\n%s", content)
	err := database.CreatePage(
		h.db,
		title,
		slug,
		spaceSlug,
		content,
	)

	if err != nil {
		log.Println("Failed to create page:", err)
		http.Error(w, "Unable to create page", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/pages/"+slug,
		http.StatusSeeOther,
	)
}

func (h *Handler) editPage(w http.ResponseWriter, r *http.Request) {
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

	spaces, err := database.FindSpaces(h.db)
	if err != nil {
		log.Println("Failed to find spaces:", err)
		http.Error(w, "Unable to load spaces", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title  string
		Page   domain.Page
		Spaces []domain.Space
	}{
		Title:  "Edit " + page.Title,
		Page:   page,
		Spaces: spaces,
	}

	if err := h.editPageTemplate.ExecuteTemplate(w, "base", data); err != nil {
		log.Println("Failed to execute edit page template:", err)
	}
}

func (h *Handler) updatePage(w http.ResponseWriter, r *http.Request) {
	originalSlug := chi.URLParam(r, "slug")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	slug := r.FormValue("slug")
	spaceSlug := r.FormValue("space_slug")
	content := r.FormValue("content")

	if title == "" || slug == "" || spaceSlug == "" || content == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := database.UpdatePage(
		h.db,
		originalSlug,
		title,
		slug,
		spaceSlug,
		content,
	)

	if err != nil {
		log.Println("Failed to update page:", err)
		http.Error(w, "Unable to update page", http.StatusInternalServerError)
		return
	}

	http.Redirect(
		w,
		r,
		"/pages/"+slug,
		http.StatusSeeOther,
	)
}
