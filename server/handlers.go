package server

import (
	"context"
	"encoding/json"
	"errors"
	"something/database"
	"something/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *server) getAllBookmarkHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	bookmarks, err := s.store.ListBookmarks(context.Background(), storage.BookmarkFilter{UserID: id})
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, bookmarks)
}

func (s *server) postBookmarkHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	var bookmark model.Bookmark
	if err := json.NewDecoder(c.Request.Body).Decode(&bookmark); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	bookmark.UserID = id
	bookmark.ID = uuid.New()
	tags := make([]string, len(bookmark.Tags))
	copy(tags, bookmark.Tags.Names())
	bookmark.Tags = []model.Tag{}
	s.store.CreateBookmark(context.Background(), bookmark)
	s.store.AddTagToBookmark(context.Background(), &bookmark, tags)

	s.index.IndexBookmark(bookmark)

	c.JSON(201, bookmark)
}

func (s *server) putBookmarkHandler(c *gin.Context) {
}

func (s *server) deleteBookmarkHandler(c *gin.Context) {
}


func (s *server) getAllNoteHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	notes, err := s.store.ListNotes(context.Background(), storage.NoteFilter{UserID: id})
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, notes)
}

func (s *server) postNoteHandler(c *gin.Context) {
}

func (s *server) putNoteHandler(c *gin.Context) {
}

func (s *server) deleteNoteHandler(c *gin.Context) {
}


func (s *server) searchHandler(c *gin.Context) {
}

func (s *server) authHandler(c *gin.Context) {
	var LoginData loginData
	if err := json.NewDecoder(c.Request.Body).Decode(&LoginData); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
	}

	tokenPair, err := s.auth.Login(LoginData.Username, LoginData.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
	}

	c.JSON(200, tokenPair)
}

func extractUserId(c *gin.Context) (uuid.UUID, error) {
	uid, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, errors.New("internal error")
	}

	id, err := uuid.Parse(uid.(string))
	if err != nil {
		return uuid.Nil, errors.New("internal error")
	}

	return id, nil
}