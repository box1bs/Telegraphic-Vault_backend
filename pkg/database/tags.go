package storage

import (
	"context"
	"errors"
	"slices"
	"something/pkg/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagFilter struct {
	UserID 	uuid.UUID `json:"user_id"`
}

func (p *Postgres) createTag(ctx context.Context, name string) (*model.Tag, error) {
	tag := &model.Tag{ Name: normalizeName(name) }
	if err := p.db.WithContext(ctx).FirstOrCreate(tag, model.Tag{Name: normalizeName(name)}).Error; err != nil {
		return nil, err
	}

	if tag.Count > 0 {
		if err := p.db.WithContext(ctx).Model(tag).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			return nil, err
		}
	}

	return tag, nil
}

func (p *Postgres) AddTagToNote(ctx context.Context, note *model.Note, tagNames []string) error {
	var tags []model.Tag
	for _, name := range tagNames {
		tag, err := p.createTag(ctx, name)
		if err != nil {
			return err
		}
		tags = append(tags, *tag)
	}

	return p.db.WithContext(ctx).Model(note).Association("Tags").Append(tags)
}

func (p *Postgres) updateNoteTags(ctx context.Context, note *model.Note, newTagNames []string) error {
	newTagNames = removeDuplicates(newTagNames)

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var removingTags []model.Tag
		for _, tag := range note.Tags {
			if !slices.Contains(newTagNames, tag.Name) {
				removingTags = append(removingTags, tag)
			} else {
				newTagNames = remove(newTagNames, tag.Name)
			}
		}

		if len(removingTags) > 0 {
			if err := tx.Model(note).Association("Tags").Delete(removingTags); err != nil {
				return err
			}

			if err := p.db.WithContext(ctx).
				Model(&model.Tag{}).
				Where("id IN ?", ExtractTagIDs(removingTags)).
				UpdateColumn("count", gorm.Expr("GREATEST(count - ?, 0)", 1)).
				Error; err != nil {
				return err
			}
		}

		if len(newTagNames) > 0 {
			if err := p.AddTagToNote(ctx, note, newTagNames); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := p.db.WithContext(ctx).Preload("Tags").First(note, note.ID).Error; err != nil {
		return err
	}

	return err
}

func (p *Postgres) AddTagToBookmark(ctx context.Context, bookmark *model.Bookmark, tagNames []string) error {
	var tags []model.Tag
	for _, name := range tagNames {
		tag, err := p.createTag(ctx, name)
		if err != nil {
			return err
		}
		tags = append(tags, *tag)
	}

	return p.db.WithContext(ctx).Model(bookmark).Association("Tags").Append(tags)
}

func (p *Postgres) updateBookmarkTags(ctx context.Context, bookmark *model.Bookmark, newTagNames []string) error {
	newTagNames = removeDuplicates(newTagNames)

	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var removingTags []model.Tag
		for _, tag := range bookmark.Tags {
			if !slices.Contains(newTagNames, tag.Name) {
				removingTags = append(removingTags, tag)
			} else {
				newTagNames = remove(newTagNames, tag.Name)
			}
		}

		if len(removingTags) > 0 {
			if err := tx.Model(bookmark).Association("Tags").Delete(removingTags); err != nil {
				return err
			}

			if err := p.db.WithContext(ctx).
				Model(&model.Tag{}).
				Where("id IN ?", ExtractTagIDs(removingTags)).
				UpdateColumn("count", gorm.Expr("GREATEST(count - ?, 0)", 1)).
				Error; err != nil {
				return err
			}
		}

		if len(newTagNames) > 0 {
			if err := p.AddTagToBookmark(ctx, bookmark, newTagNames); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := p.db.WithContext(ctx).Preload("Tags").First(bookmark, bookmark.ID).Error; err != nil {
		return err
	}

	return err
}

func remove(list []string, item string) []string {
    for i, v := range list {
        if v == item {
            copy(list[i:], list[i+1:])
            list[len(list)-1] = ""
            list = list[:len(list)-1]
        }
    }
    return list
}

func removeDuplicates(names []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, name := range names {
        if _, exists := seen[name]; !exists {
            seen[name] = struct{}{}
            result = append(result, name)
        }
    }
    return result
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

func ExtractTagIDs(tags []model.Tag) []uuid.UUID {
    ids := make([]uuid.UUID, len(tags))
    for i, tag := range tags {
        ids[i] = tag.ID
    }
    return ids
}