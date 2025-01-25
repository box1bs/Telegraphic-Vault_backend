package server

import (
	"github.com/box1bs/TelegraphicVault/pkg/auth"
	"github.com/box1bs/TelegraphicVault/pkg/config"
	"github.com/box1bs/TelegraphicVault/pkg/database"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type server struct {
	store 		storage.Storage
	auth 		*auth.AuthService
	keyStore 	*sync.Map
	rateLimiter map[string]*rate.Limiter
	mu 			*sync.Mutex
}

func NewServer(store storage.Storage, conf *config.AuthConfig) *server {
	return &server{
		store: store,
		auth: auth.NewAuthService(conf, store),
		keyStore: &sync.Map{},
		rateLimiter: make(map[string]*rate.Limiter),
		mu: new(sync.Mutex),
	}
}

func (s *server) Run() error {
	r := gin.Default()
	s.registerRoutes(r)
	return r.Run()
}

func (s *server) registerRoutes(r *gin.Engine) {
	r.GET("/auth", s.keyHandler) // encryption key for encrypt password
	r.POST("/auth", s.registerHandler) // for registration
	r.POST("/auth/login", s.loginHandler) // for login
	app := r.Group("/app", s.auth.AuthMiddleware())
	{
		bookmarks := app.Group("/bookmarks")
		{
			bookmarks.GET("", s.getAllBookmarkHandler)
			bookmarks.POST("", s.postBookmarkHandler)
			bookmarks.PUT("", s.putBookmarkHandler)
			bookmarks.DELETE("", s.deleteBookmarkHandler)
			bookmarks.GET("/search", s.searchBookmarkHandler)
		}
		
		notes := app.Group("/notes")
		{
			notes.GET("", s.getAllNoteHandler)
			notes.POST("", s.postNoteHandler)
			notes.PUT("", s.putNoteHandler)
			notes.DELETE("", s.deleteNoteHandler)
			notes.GET("/search", s.searchNoteHandler)
		}
	}
}