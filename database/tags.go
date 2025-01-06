package storage

import (
	"context"
	"errors"
	"something/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagFilter struct {
	UserID 	uuid.UUID `json:"user_id"`
}

func (p *Postgres) CreateTag(ctx context.Context, name string) (*model.Tag, error) {
	tag := &model.Tag{ Name: normalizeName(name) }
	if err := p.db.WithContext(ctx).FirstOrCreate(tag, model.Tag{Name: normalizeName(name)}).Error; err != nil {
		return nil, err
	}

	if err := p.db.WithContext(ctx).Model(tag).UpdateColumn("counter", gorm.Expr("counter + ?", 1)).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

func (p *Postgres) AddTagToNote(ctx context.Context, note *model.Note, tagNames []string) error {
	var tags []model.Tag
	for _, name := range tagNames {
		tag, err := p.CreateTag(ctx, name)
		if err != nil {
			return err
		}
		tags = append(tags, *tag)
	}

	return p.db.WithContext(ctx).Model(note).Association("Tags").Append(tags)
}

func (p *Postgres) AddTagToBookmark(ctx context.Context, bookmark *model.Bookmark, tagNames []string) error {
	var tags []model.Tag
	for _, name := range tagNames {
		tag, err := p.CreateTag(ctx, name)
		if err != nil {
			return err
		}
		tags = append(tags, *tag)
	}

	return p.db.WithContext(ctx).Model(bookmark).Association("Tags").Append(tags)
}

func (p *Postgres) FindByTag(ctx context.Context, tagName string) ([]*model.Note, []*model.Bookmark, error) {
	var notes []*model.Note
	var bookmarks []*model.Bookmark
	var tag model.Tag

	if err := p.db.WithContext(ctx).Where("name = ?", normalizeName(tagName)).First(&tag).Error; err != nil {
		return nil, nil, err
	}

	if err := p.db.Joins("JOIN note_tags ON notes.id = note_tags.note_id").
		Where("note_tags.tag_id = ?", tag.ID).
		Find(&notes).Error; err != nil {
		return nil, nil, err
	}

	if err := p.db.Joins("JOIN bookmark_tags ON bookmarks.id = bookmark_tags.bookmark_id").
		Where("bookmark_tags.tag_id = ?", tag.ID).
		Find(&bookmarks).Error; err != nil {
		return nil, nil, err
	}

	return notes, bookmarks, nil
}

func (p *Postgres) SerchByTags(ctx context.Context, tagNames []string) ([]*model.Note, []*model.Bookmark, error) {
	var notes []*model.Note
	var bookmarks []*model.Bookmark
	var tags []model.Tag

	for _, name := range tagNames {
		tag := model.Tag{Name: normalizeName(name)}
		tags = append(tags, tag)
	}

	if err := p.db.WithContext(ctx).Where("name IN ?", tags).Find(&tags).Error; err != nil { // overwrites tags
		return nil, nil, err
	}

	if len(tags) == 0 {
		return nil, nil, nil
	}

	if err := p.db.Joins("JOIN note_tags ON notes.id = note_tags.note_id").
		Where("note_tags.tag_id IN ?", tags).Find(&notes).Error; err != nil {
		return nil, nil, err
	}

	if err := p.db.Joins("JOIN bookmark_tags ON bookmarks.id = bookmark_tags.bookmark_id").
		Where("bookmark_tags.tag_id IN ?", tags).Find(&bookmarks).Error; err != nil {
		return nil, nil, err
	}

	return notes, bookmarks, nil
}

func (p* Postgres) GetPopularTags(ctx context.Context, filter TagFilter, limit int) ([]*model.Tag, error) {
	var tags []*model.Tag
	query := p.db.WithContext(ctx)

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	} else {
		return nil, errors.New("user_id is required")
	}

	err := query.Order("counter DESC").Limit(limit).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(strings.Replace(name, " ", "-", -1)))
}