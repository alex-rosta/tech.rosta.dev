package handlers

import (
	"net/http"
	"strings"
	"time"

	"goblog/internal/models"

	"github.com/go-chi/chi/v5"
)

type AdminDashboardData struct {
	Title string
	Posts []*models.Post
}

type PostFormData struct {
	Title string
	Post  *models.Post
	Mode  string // "new" or "edit"
}

func (h *Handlers) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	posts, err := h.storage.ListPosts()
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	data := AdminDashboardData{
		Title: "Admin Dashboard",
		Posts: posts,
	}

	h.renderTemplate(w, "admin_dashboard.html", data)
}

func (h *Handlers) handleNewPost(w http.ResponseWriter, r *http.Request) {
	data := PostFormData{
		Title: "New Post",
		Post:  &models.Post{},
		Mode:  "new",
	}

	h.renderTemplate(w, "admin_post_form.html", data)
}

func (h *Handlers) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	tagsStr := r.FormValue("tags")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Generate slug from title
	slug := generateSlug(title)

	// Parse tags
	var tags []string
	if tagsStr != "" {
		tagList := strings.Split(tagsStr, ",")
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	post := models.NewPost(slug, title, content)
	post.Tags = tags

	if err := h.storage.SavePost(post); err != nil {
		http.Error(w, "Error saving post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handlers) handleEditPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	post, err := h.storage.GetPost(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := PostFormData{
		Title: "Edit Post: " + post.Title,
		Post:  post,
		Mode:  "edit",
	}

	h.renderTemplate(w, "admin_post_form.html", data)
}

func (h *Handlers) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	post, err := h.storage.GetPost(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	tagsStr := r.FormValue("tags")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Parse tags
	var tags []string
	if tagsStr != "" {
		tagList := strings.Split(tagsStr, ",")
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	post.Title = title
	post.Content = content
	post.Tags = tags
	post.UpdatedAt = time.Now()

	if err := h.storage.SavePost(post); err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handlers) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	if err := h.storage.DeletePost(slug); err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove special characters
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	return result.String()
}
