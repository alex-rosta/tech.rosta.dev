package handlers

import (
	"net/http"

	"goblog/internal/models"
)

type HomeData struct {
	Title       string
	RecentPosts []*models.Post
	AllTags     []string
	Breadcrumb  string
}

func (h *Handlers) handleHome(w http.ResponseWriter, r *http.Request) {
	// Get recent posts
	recentPosts, err := h.storage.GetRecentPosts(6)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	// Get all unique tags
	allPosts, err := h.storage.ListPosts()
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	tagSet := make(map[string]bool)
	for _, post := range allPosts {
		for _, tag := range post.Tags {
			tagSet[tag] = true
		}
	}

	var allTags []string
	for tag := range tagSet {
		allTags = append(allTags, tag)
	}

	data := HomeData{
		Title:       "Dashboard",
		RecentPosts: recentPosts,
		AllTags:     allTags,
	}

	h.renderTemplate(w, "home.html", data)
}
