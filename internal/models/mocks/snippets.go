package mocks

import (
	"asniki/snippetbox/internal/models"
	"time"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

// SnippetModel mocks models.SnippetModel
type SnippetModel struct{}

// Insert mocks models.SnippetModel.Insert
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

// Get mocks models.SnippetModel.Get
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Latest mocks models.SnippetModel.Latest
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
