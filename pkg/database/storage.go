package storage

import (
	"context"
	"github.com/box1bs/TelegraphicVault/pkg/model"

	"github.com/google/uuid"
)

type Storage interface {
	UserStorage
	tagStorage
	noteStorage
	bookmarkStorage
}

type JWTUserStorage interface {
	FindByID(uuid.UUID) (*model.User, error)
	FindByUsername(string) (*model.User, error)
	LastLoginUpdate(*model.User) error
}

type UserStorage interface {
	SaveUser(*model.User) error
	FindByID(uuid.UUID) (*model.User, error)
	FindByUsername(string) (*model.User, error)
	LastLoginUpdate(*model.User) error
}

type bookmarkStorage interface {
    CreateBookmark(context.Context, model.Bookmark) error
	SearchBookmark(context.Context, uuid.UUID, string) ([]model.Bookmark, error)
    UpdateBookmark(context.Context, uuid.UUID, string, string, string, []string) (*model.Bookmark, error)
    DeleteBookmark(context.Context, uuid.UUID, string) error
    ListBookmarks(context.Context, BookmarkFilter) ([]*model.Bookmark, error)
}

type noteStorage interface {
    CreateNote(context.Context, model.Note) error
    GetNote(context.Context, uuid.UUID, string) (*model.Note, error)
    UpdateNote(context.Context, uuid.UUID, string, string, string, []string) (*model.Note, error)
    DeleteNote(context.Context, uuid.UUID, string) error
    ListNotes(context.Context, NoteFilter) ([]*model.Note, error)
}

type tagStorage interface {
    createTag(context.Context, string) (*model.Tag, error)
	AddTagToNote(context.Context, *model.Note, []string) error
	AddTagToBookmark(context.Context, *model.Bookmark, []string) error
	FindByTag(context.Context, string) ([]*model.Note, []*model.Bookmark, error)
	GetPopularTags(context.Context, TagFilter, int) ([]*model.Tag, error)
}