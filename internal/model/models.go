package model

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

var layoutModules = []string{"overview", "resources", "posts", "ai-lab"}

type Resource struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Type        string     `json:"type" db:"type"`
	Description string     `json:"description" db:"description"`
	URL         *string    `json:"url,omitempty" db:"url"`
	FilePath    *string    `json:"filePath,omitempty" db:"file_path"`
	Tags        []string   `json:"tags" db:"tags"`
	AuthorID    string     `json:"authorId" db:"author_id"`
	Status      string     `json:"status" db:"status"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" db:"published_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

type Post struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Summary     string     `json:"summary" db:"summary"`
	Body        string     `json:"body" db:"body"`
	CoverImage  *string    `json:"coverImage,omitempty" db:"cover_image"`
	Tags        []string   `json:"tags" db:"tags"`
	AuthorID    string     `json:"authorId" db:"author_id"`
	Status      string     `json:"status" db:"status"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" db:"published_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

type AIProduct struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Summary     string     `json:"summary" db:"summary"`
	URL         string     `json:"url" db:"url"`
	CoverImage  *string    `json:"coverImage,omitempty" db:"cover_image"`
	Tags        []string   `json:"tags" db:"tags"`
	Status      string     `json:"status" db:"status"`
	SortOrder   int        `json:"sortOrder" db:"sort_order"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" db:"published_at"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

type Page[T any] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type ListFilter struct {
	Page     int
	PageSize int
	Type     string
	Tag      string
}

func (f ListFilter) Normalized() ListFilter {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	return f
}

type Dashboard struct {
	Resources  []Resource  `json:"resources"`
	Posts      []Post      `json:"posts"`
	AIProducts []AIProduct `json:"aiProducts"`
}

type Layout struct {
	Modules []string `json:"modules"`
	Hidden  []string `json:"hidden"`
	Density string   `json:"density"`
	Theme   string   `json:"theme"`
}

func DefaultLayout() Layout {
	return Layout{
		Modules: slices.Clone(layoutModules),
		Hidden:  []string{},
		Density: "comfortable",
		Theme:   "light",
	}
}

func (l Layout) Validate() error {
	if len(l.Modules) != len(layoutModules) {
		return fmt.Errorf("modules must contain all supported modules")
	}
	seen := make(map[string]bool, len(l.Modules))
	for _, module := range l.Modules {
		if !slices.Contains(layoutModules, module) {
			return fmt.Errorf("unknown module %q", module)
		}
		if seen[module] {
			return fmt.Errorf("duplicate module %q", module)
		}
		seen[module] = true
	}
	for _, module := range layoutModules {
		if !seen[module] {
			return fmt.Errorf("missing module %q", module)
		}
	}
	hidden := make(map[string]bool, len(l.Hidden))
	for _, module := range l.Hidden {
		if !seen[module] {
			return fmt.Errorf("hidden module %q is not enabled", module)
		}
		if hidden[module] {
			return fmt.Errorf("duplicate hidden module %q", module)
		}
		hidden[module] = true
	}
	if l.Density != "compact" && l.Density != "comfortable" {
		return errors.New("density must be compact or comfortable")
	}
	if l.Theme != "light" && l.Theme != "dark" {
		return errors.New("theme must be light or dark")
	}
	return nil
}

func ValidStatus(status string) bool {
	return status == StatusDraft || status == StatusPublished || status == StatusArchived
}
