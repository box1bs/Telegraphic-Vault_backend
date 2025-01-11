package fts

import (
	"something/model"
	"time"
)

type IndexBookmark struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TagNames    []string  `json:"tag_names"`
	TagIds      []string  `json:"tag_ids"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type IndexNote struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TagNames  []string  `json:"tag_names"`
	TagIds    []string  `json:"tag_ids"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *SearchService) IndexBookmark(bookmark model.Bookmark) error {
	return s.index.Index(bookmark.ID.String(), IndexBookmark{
		ID:          bookmark.ID.String(),
		URL:         bookmark.URL,
		Title:       bookmark.Title,
		Description: bookmark.Description,
		TagNames:    bookmark.Tags.Names(),
		TagIds:      bookmark.Tags.IDs(),
		UserID:      bookmark.UserID.String(),
		CreatedAt:   bookmark.CreatedAt,
	})
}

func (s *SearchService) DeleteBookmark(id string) error {
	return s.index.Delete(id)
}

func (s *SearchService) IndexNote(note model.Note) error {
	return s.index.Index(note.ID.String(), IndexNote{
		ID:        note.ID.String(),
		Title:     note.Title,
		Content:   note.Content,
		TagNames:  note.Tags.Names(),
		TagIds:    note.Tags.IDs(),
		UserID:    note.UserID.String(),
		CreatedAt: note.CreatedAt,
	})
}

func (s *SearchService) DeleteNote(id string) error {
	return s.index.Delete(id)
}