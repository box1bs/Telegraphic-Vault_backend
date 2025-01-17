package storage

import (
	"context"
	"something/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookmarkFilter struct {
	UserID uuid.UUID 	`json:"user_id"`
	Tag    string		`json:"tag"`
}

func (p *Postgres) CreateBookmark(ctx context.Context, bookmark model.Bookmark) error {
	return p.db.WithContext(ctx).Create(&bookmark).Error
}

func (p *Postgres) GetBookmark(ctx context.Context, userID uuid.UUID, uri string) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	err := p.db.WithContext(ctx).Where("user_id = ? AND url = ?", userID, uri).First(&bookmark).Error
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

func (p *Postgres) UpdateBookmark(ctx context.Context, userID uuid.UUID, uri, title, description string, newTagNames []string) (*model.Bookmark, error) {
	bookmark, err := p.GetBookmark(ctx, userID, uri)
	if err != nil {
		return nil, err
	}

	bookmark.Title = title
	bookmark.Description = description
	
	if newTagNames != nil {
		p.AddTagToBookmark(ctx, bookmark, newTagNames)
	}

	err = p.db.WithContext(ctx).Save(bookmark).Error
	return bookmark, err
}

func (p *Postgres) DeleteBookmark(ctx context.Context, userID uuid.UUID, uri string) error {
	result := p.db.WithContext(ctx).Where("user_id = ? AND url = ?", userID, uri).Delete(&model.Bookmark{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *Postgres) ListBookmarks(ctx context.Context, filter BookmarkFilter) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark
	query := p.db.WithContext(ctx).Model(&model.Bookmark{})

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Tag != "" {
	query = query.Joins("JOIN bookmark_tags ON bookmarks.id = bookmark_tags.bookmark_id").
		Joins("JOIN tags ON tags.id = bookmark_tags.tag_id").
		Where("tags.name = ?", filter.Tag)
	}

	err := query.Preload("Tags").Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}