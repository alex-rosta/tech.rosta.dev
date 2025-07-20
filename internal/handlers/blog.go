package handlers

import (
	"html/template"
	"net/http"

	"goblog/internal/models"
	"goblog/pkg/markdown"

	"github.com/go-chi/chi/v5"
)

type PostData struct {
	Title string
	Post  *models.Post
	HTML  template.HTML
}

type TagData struct {
	Title string
	Tag   string
	Posts []*models.Post
}

type SearchData struct {
	Title   string
	Query   string
	Results []*models.Post
}

func (h *Handlers) handlePost(w http.ResponseWriter, r *http.Request) {
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

	// Convert markdown to HTML
	html := markdown.ToHTML(post.Content)

	data := PostData{
		Title: post.Title,
		Post:  post,
		HTML:  template.HTML(html),
	}

	h.renderTemplate(w, "post.html", data)
}

func (h *Handlers) handleTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	if tag == "" {
		http.NotFound(w, r)
		return
	}

	posts, err := h.storage.GetPostsByTag(tag)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	data := TagData{
		Title: "Posts tagged: " + tag,
		Tag:   tag,
		Posts: posts,
	}

	h.renderTemplate(w, "tag.html", data)
}

func (h *Handlers) handleAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.storage.ListPosts()
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	data := SearchData{
		Title:   "All Posts",
		Query:   "",
		Results: posts,
	}

	h.renderTemplate(w, "search.html", data)
}

func (h *Handlers) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var results []*models.Post
	var err error

	if query != "" {
		results, err = h.storage.SearchPosts(query)
		if err != nil {
			http.Error(w, "Error searching posts", http.StatusInternalServerError)
			return
		}
	}

	data := SearchData{
		Title:   "Search Results",
		Query:   query,
		Results: results,
	}

	h.renderTemplate(w, "search.html", data)
}
