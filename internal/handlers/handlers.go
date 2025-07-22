package handlers

import (
	"html/template"
	"log"
	"net/http"

	"goblog/internal/config"
	"goblog/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	storage storage.Storage
	config  *config.Config
	//templates *template.Template
}

func New(storage storage.Storage, config *config.Config) *Handlers {
	return &Handlers{
		storage: storage,
		config:  config,
	}
}

func (h *Handlers) SetupRoutes(r *chi.Mux) {
	// Public routes
	r.Get("/", h.handleHome)
	r.Get("/posts", h.handleAllPosts)
	r.Get("/post/{slug}", h.handlePost)
	r.Get("/tag/{tag}", h.handleTag)
	r.Get("/search", h.handleSearch)
}

func (h *Handlers) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	// Load templates dynamically for each request to avoid conflicts
	tmpl := template.New("base.html")

	// Parse base template
	tmpl, err := tmpl.ParseFiles("web/templates/layouts/base.html")
	if err != nil {
		log.Printf("Template parse error (base): %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the specific page template
	tmpl, err = tmpl.ParseFiles("web/templates/pages/" + name)
	if err != nil {
		log.Printf("Template parse error (%s): %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the base template
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Template execute error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
