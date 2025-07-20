package storage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"goblog/internal/models"
)

type FilesystemStorage struct {
	postsDir string
}

func NewFilesystemStorage(postsDir string) *FilesystemStorage {
	// Ensure posts directory exists
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create posts directory: %v", err))
	}

	return &FilesystemStorage{
		postsDir: postsDir,
	}
}

func (fs *FilesystemStorage) ListPosts() ([]*models.Post, error) {
	var posts []*models.Post

	err := filepath.WalkDir(fs.postsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			post, err := fs.parseMarkdownFile(path)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}
			posts = append(posts, post)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by creation time, newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	return posts, nil
}

func (fs *FilesystemStorage) GetPost(slug string) (*models.Post, error) {
	filename := slug + ".md"
	path := filepath.Join(fs.postsDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("post not found: %s", slug)
	}

	return fs.parseMarkdownFile(path)
}

func (fs *FilesystemStorage) SavePost(post *models.Post) error {
	filename := post.Slug + ".md"
	path := filepath.Join(fs.postsDir, filename)

	content := fs.formatMarkdownContent(post)

	return os.WriteFile(path, []byte(content), 0644)
}

func (fs *FilesystemStorage) DeletePost(slug string) error {
	filename := slug + ".md"
	path := filepath.Join(fs.postsDir, filename)

	return os.Remove(path)
}

func (fs *FilesystemStorage) SearchPosts(query string) ([]*models.Post, error) {
	posts, err := fs.ListPosts()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var matches []*models.Post

	for _, post := range posts {
		if fs.matchesQuery(post, query) {
			matches = append(matches, post)
		}
	}

	return matches, nil
}

func (fs *FilesystemStorage) GetPostsByTag(tag string) ([]*models.Post, error) {
	posts, err := fs.ListPosts()
	if err != nil {
		return nil, err
	}

	var matches []*models.Post
	for _, post := range posts {
		for _, postTag := range post.Tags {
			if strings.EqualFold(postTag, tag) {
				matches = append(matches, post)
				break
			}
		}
	}

	return matches, nil
}

func (fs *FilesystemStorage) GetRecentPosts(limit int) ([]*models.Post, error) {
	posts, err := fs.ListPosts()
	if err != nil {
		return nil, err
	}

	if len(posts) <= limit {
		return posts, nil
	}

	return posts[:limit], nil
}

func (fs *FilesystemStorage) parseMarkdownFile(path string) (*models.Post, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	slug := strings.TrimSuffix(filepath.Base(path), ".md")
	post := &models.Post{
		Slug:      slug,
		CreatedAt: stat.ModTime(),
		UpdatedAt: stat.ModTime(),
		Metadata:  make(map[string]string),
		Tags:      []string{},
	}

	scanner := bufio.NewScanner(file)
	var contentLines []string
	inFrontMatter := false
	frontMatterEnded := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" && !frontMatterEnded {
			if !inFrontMatter {
				inFrontMatter = true
			} else {
				frontMatterEnded = true
				inFrontMatter = false
			}
			continue
		}

		if inFrontMatter {
			fs.parseFrontMatterLine(post, line)
		} else if frontMatterEnded || !strings.HasPrefix(line, "---") {
			contentLines = append(contentLines, line)
		}
	}

	post.Content = strings.Join(contentLines, "\n")

	// Extract title from content if not in front matter
	if post.Title == "" && len(contentLines) > 0 {
		for _, line := range contentLines {
			if strings.HasPrefix(line, "# ") {
				post.Title = strings.TrimSpace(strings.TrimPrefix(line, "#"))
				break
			}
		}
	}

	return post, scanner.Err()
}

func (fs *FilesystemStorage) parseFrontMatterLine(post *models.Post, line string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	value = strings.Trim(value, `"'`)

	switch key {
	case "title":
		post.Title = value
	case "tags":
		if value != "" {
			tags := strings.Split(value, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
			post.Tags = tags
		}
	case "created":
		if t, err := time.Parse("2006-01-02", value); err == nil {
			post.CreatedAt = t
		}
	case "updated":
		if t, err := time.Parse("2006-01-02", value); err == nil {
			post.UpdatedAt = t
		}
	default:
		post.Metadata[key] = value
	}
}

func (fs *FilesystemStorage) formatMarkdownContent(post *models.Post) string {
	var builder strings.Builder

	// Front matter
	builder.WriteString("---\n")
	builder.WriteString(fmt.Sprintf("title: \"%s\"\n", post.Title))

	if len(post.Tags) > 0 {
		builder.WriteString(fmt.Sprintf("tags: %s\n", strings.Join(post.Tags, ", ")))
	}

	builder.WriteString(fmt.Sprintf("created: %s\n", post.CreatedAt.Format("2006-01-02")))
	builder.WriteString(fmt.Sprintf("updated: %s\n", post.UpdatedAt.Format("2006-01-02")))

	for key, value := range post.Metadata {
		builder.WriteString(fmt.Sprintf("%s: \"%s\"\n", key, value))
	}

	builder.WriteString("---\n\n")

	// Content
	builder.WriteString(post.Content)

	return builder.String()
}

func (fs *FilesystemStorage) matchesQuery(post *models.Post, query string) bool {
	// Search in title
	if strings.Contains(strings.ToLower(post.Title), query) {
		return true
	}

	// Search in tags
	for _, tag := range post.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	// Search in content (headers)
	lines := strings.Split(post.Content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") && strings.Contains(strings.ToLower(line), query) {
			return true
		}
	}

	return false
}
