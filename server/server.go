package server

import (
	"log"
	"something/FTS"
	"something/auth"
	"something/config"
	"something/database"

	"github.com/gin-gonic/gin"
)

type server struct {
	store 	storage.Storage
	index 	*fts.SearchService
	auth 	*auth.AuthService
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewServer(store storage.Storage, conf *config.AuthConfig, indexPath string) *server {
	index, err := fts.NewSearchService(indexPath)
	if err != nil {
		log.Fatalf("Failed to create search service: %v", err)
	}

	return &server{
		store: store,
		index: index,
		auth: auth.NewAuthService(conf, store),
	}
}

func (s *server) Run() error {
	r := gin.Default()
	s.registerRoutes(r)
	return r.Run()
}

func (s *server) registerRoutes(r *gin.Engine) {
	r.POST("/auth", s.authHandler)
	app := r.Group("/app", s.auth.AuthMiddleware())
	{
		bookmarks := app.Group("/bookmarks")
		{
			bookmarks.GET("", s.getAllBookmarkHandler)
			bookmarks.POST("", s.postBookmarkHandler)
			bookmarks.PUT("", s.putBookmarkHandler)
			bookmarks.DELETE("", s.deleteBookmarkHandler)
		}
		
		notes := app.Group("/notes")
		{
			notes.GET("", s.getAllNoteHandler)
			notes.POST("", s.postNoteHandler)
			notes.PUT("", s.putNoteHandler)
			notes.DELETE("", s.deleteNoteHandler)
		}

		search := app.Group("/search")
		{
			search.GET("", s.searchHandler)
		}
	}
}