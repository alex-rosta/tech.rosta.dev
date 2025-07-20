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
	storage   storage.Storage
	config    *config.Config
	templates *template.Template
}

func New(storage storage.Storage, config *config.Config) *Handlers {
	// Load templates
	templates := template.New("")

	// Parse all template files
	template.Must(templates.ParseGlob("web/templates/layouts/*.html"))
	template.Must(templates.ParseGlob("web/templates/pages/*.html"))

	return &Handlers{
		storage:   storage,
		config:    config,
		templates: templates,
	}
}

func (h *Handlers) SetupRoutes(r *chi.Mux) {
	// Public routes
	r.Get("/", h.handleHome)
	r.Get("/posts", h.handleAllPosts)
	r.Get("/post/{slug}", h.handlePost)
	r.Get("/tag/{tag}", h.handleTag)
	r.Get("/search", h.handleSearch)

	// Auth routes
	r.Get("/login", h.handleLogin)
	r.Get("/auth/github", h.handleGitHubLogin)
	r.Get("/auth/callback", h.handleAuthCallback)
	r.Post("/logout", h.handleLogout)

	// Admin routes (protected)
	r.Route("/admin", func(r chi.Router) {
		r.Use(h.requireAuth)
		r.Get("/", h.handleAdminDashboard)
		r.Get("/new", h.handleNewPost)
		r.Post("/posts", h.handleCreatePost)
		r.Get("/posts/{slug}/edit", h.handleEditPost)
		r.Put("/posts/{slug}", h.handleUpdatePost)
		r.Delete("/posts/{slug}", h.handleDeletePost)
	})
}

func (h *Handlers) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handlers) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := h.getSessionUser(r); !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
