package storage

import (
	"context"
	"something/model"

	"github.com/google/uuid"
)

type BookmarkFilter struct {
	UserID uuid.UUID 	`json:"user_id"`
	Tag    string		`json:"tag"`
}

func (p *Postgres) CreateBookmark(ctx context.Context, bookmark model.Bookmark) error {
	return p.db.WithContext(ctx).Create(&bookmark).Error
}

func (p *Postgres) GetBookmark(ctx context.Context, id uuid.UUID) (*model.Bookmark, error) {
	var bookmark *model.Bookmark
	err := p.db.WithContext(ctx).Where("id = ?", id).First(bookmark).Error
	if err != nil {
		return nil, err
	}
	return bookmark, nil
}

func (p *Postgres) UpdateBookmark(ctx context.Context, bookmark *model.Bookmark) error {
	return p.db.WithContext(ctx).Save(bookmark).Error
}

func (p *Postgres) DeleteBookmark(ctx context.Context, id uuid.UUID) error {
	return p.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Bookmark{}).Error
}

func (p *Postgres) ListBookmarks(ctx context.Context, filter BookmarkFilter) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark
	query := p.db.WithContext(ctx)

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.Tag != "" {
		query = query.Where("? = ANY(tags)", filter.Tag)
	}

	err := query.Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}