package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/box1bs/TelegraphicVault/pkg/auth"
	"github.com/box1bs/TelegraphicVault/pkg/database"
	"github.com/box1bs/TelegraphicVault/pkg/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *server) getAllBookmarkHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	bookmarks, err := s.store.ListBookmarks(context.Background(), storage.BookmarkFilter{UserID: id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "bookmark not found"})
			return
		}
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

	var payload struct {
		Url 		string 			`json:"url"`
		Title 		string 			`json:"title"`
		Description string 			`json:"description"`
		Tags 		[]string 		`json:"tags"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	bookmark := model.Bookmark{
		ID: uuid.New(),
		URL: payload.Url,
		Title: payload.Title,
		Description: payload.Description,
		UserID: id,
	}
	s.store.CreateBookmark(context.Background(), bookmark)
	s.store.AddTagToBookmark(context.Background(), &bookmark, payload.Tags)

	c.JSON(201, bookmark)
}

func (s *server) putBookmarkHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	var payload struct {
		URL 			string `json:"url"`
		Title 			string `json:"title"`
		Description 	string `json:"description"`
		Tags 			[]string `json:"tags"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	bookmark, err := s.store.UpdateBookmark(context.Background(), id, payload.URL, payload.Title, payload.Description, payload.Tags)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "bookmark not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, bookmark)
}

func (s *server) deleteBookmarkHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	uri := c.Query("uri")
	if uri == "" {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if err := s.store.DeleteBookmark(context.Background(), id, uri); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "bookmark not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(204, nil)
}

func (s *server) searchBookmarkHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	bookmarks, err := s.store.SearchBookmark(context.Background(), id, query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "note not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, bookmarks)
}

func (s *server) getAllNoteHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	notes, err := s.store.ListNotes(context.Background(), storage.NoteFilter{UserID: id})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "note not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, notes)
}

func (s *server) postNoteHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	var payload struct {
		Title 		string 			`json:"title"`
		Content 	string 			`json:"content"`
		Tags 		[]string 		`json:"tags"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	note := model.Note{
		ID: uuid.New(),
		Title: payload.Title,
		Content: payload.Content,
		UserID: id,
	}
	s.store.CreateNote(context.Background(), note)
	s.store.AddTagToNote(context.Background(), &note, payload.Tags)

	c.JSON(201, note)
}

func (s *server) putNoteHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	var payload struct {
		CurrentTitle 	string `json:"current_title"`
		NewTitle 		string `json:"new_title"`
		Content 		string `json:"content"`
		Tags 			[]string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	note, err := s.store.UpdateNote(context.Background(), id, payload.CurrentTitle, payload.NewTitle, payload.Content, payload.Tags)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "note not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, note)
}

func (s *server) deleteNoteHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	title := c.Query("title")
	if title == "" {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if err := s.store.DeleteNote(context.Background(), id, title); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "note not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(204, nil)
}

func (s *server) searchNoteHandler(c *gin.Context) {
	id, err := extractUserId(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	note, err := s.store.GetNote(context.Background(), id, query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "note not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, note)
}

func (s *server) refreshHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		if token[:7] == "Bearer " {
			tokenPair, err := s.auth.RefreshTokens(token[7:])
			if err != nil {
				c.JSON(401, gin.H{"error": "invalid token"})
				return
			}
			c.JSON(200, tokenPair)
			return
		}
		c.JSON(401, gin.H{"error": "invalid token format"})
		return
	}
	c.JSON(401, gin.H{"error": "empty token"})
}

func (s *server) registerHandler(c *gin.Context) {
	var payload struct {
		Username 			string `json:"username"`
		EncryptedPassword 	string `json:"password"`
		TempKey 			string `json:"key"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	existingUser, err := s.store.FindByUsername(payload.Username)
	if err == nil && existingUser != nil {
		c.JSON(409, gin.H{"error": "username already exists"})
		return
	}

	if _, exist := s.keyStore.Load(payload.TempKey); !exist {
		c.JSON(401, gin.H{"error": "invalid key"})
		return
	}

	password, err := auth.Decode(payload.EncryptedPassword, payload.TempKey)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid key"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	user := &model.User{
		ID: uuid.New(),
		Username: payload.Username,
		Password: string(hashedPassword),
		Role: "user",
	}

	if err := s.store.SaveUser(user); err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	tokenPair, err := s.auth.Register(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(201, tokenPair)
}

func (s *server) keyHandler(c *gin.Context) {
	key, err := auth.GenerateServerKey()
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	s.keyStore.Store(key, time.Now())
	go func(key string) {
		time.Sleep(1 * time.Minute)
		s.keyStore.Delete(key)
	}(key)

	c.JSON(200, gin.H{"key": key})
}

func (s *server) loginHandler(c *gin.Context) {
	var payload struct {
		Username 	string `json:"username"`
		Password 	string `json:"password"`
		TempKey 	string `json:"key"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if _, exist := s.keyStore.Load(payload.TempKey); !exist {
		c.JSON(401, gin.H{"error": "invalid key"})
		return
	}

	password, err := auth.Decode(payload.Password, payload.TempKey)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid key"})
		return
	}

	tokenPair, err := s.auth.Login(payload.Username, password)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
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