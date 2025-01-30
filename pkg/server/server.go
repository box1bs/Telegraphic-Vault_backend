package server

import (
	"log"
	"sync"
	"time"

	"github.com/box1bs/TelegraphicVault/pkg/auth"
	"github.com/box1bs/TelegraphicVault/pkg/config"
	"github.com/box1bs/TelegraphicVault/pkg/database"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

type server struct {
	store 		storage.Storage
	auth 		*auth.AuthService
	keyStore 	*sync.Map
	mu 			*sync.Mutex
}

func NewServer(store storage.Storage, conf *config.AuthConfig) *server {
	return &server{
		store: store,
		auth: auth.NewAuthService(conf, store),
		keyStore: &sync.Map{},
		mu: new(sync.Mutex),
	}
}

func (s *server) Run() error {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{
		"GET",
		"POST", 
		"PUT", 
		"DELETE", 
		"OPTIONS",
		},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
            return origin == "http://localhost:5173"
        },
    }))
	r.Use(func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Next()
    })

	rateLimiter := NewRateLimiter()
	r.Use(func(c *gin.Context) {
        ip := c.ClientIP()
        
        if !rateLimiter.Allow(ip) {
            c.Header("Retry-After", "3600")
            c.JSON(429, gin.H{
                "error": "Too many requests. Please try again in 1 hour.",
                "retry_after": "3600",
            })
            c.Abort()
            return
        }
        
        c.Next()
    })

	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		
		c.Next()
		
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			log.Printf("WARNING: %s %s %d %v\n", c.Request.Method, path, statusCode, latency)
		}
	})

	s.registerRoutes(r)
	return r.Run()
}

func (s *server) registerRoutes(r *gin.Engine) {
	r.GET("/auth", s.keyHandler) // encryption key for encrypt password
	r.POST("/auth", s.registerHandler) // for registration
	r.POST("/auth/login", s.loginHandler) // for login
	r.POST("/auth/refresh", s.refreshHandler)
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