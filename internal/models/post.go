package models

import (
	"time"
)

type Post struct {
	Slug      string            `json:"slug"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Tags      []string          `json:"tags"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Metadata  map[string]string `json:"metadata"`
}
