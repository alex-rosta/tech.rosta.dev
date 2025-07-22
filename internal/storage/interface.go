package storage

import "goblog/internal/models"

type Storage interface {
	ListPosts() ([]*models.Post, error)
	GetPost(slug string) (*models.Post, error)
	SearchPosts(query string) ([]*models.Post, error)
	GetPostsByTag(tag string) ([]*models.Post, error)
	GetRecentPosts(limit int) ([]*models.Post, error)
}
