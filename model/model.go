package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAlreadyExists = errors.New("record already exists")
)

type Bookmark struct {
	ID         	uuid.UUID	`json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	URL 	   	string		`json:"url" gorm:"not null"`
	Title      	string		`json:"title"`
	Description string		`json:"description"`
	Tags       	tags		`json:"tags" gorm:"many2many:bookmark_tags;"`
	UserID     	uuid.UUID	`json:"user_id" gorm:"not null"`
	CreatedAt  	time.Time	`json:"created_at" gorm:"autoCreateTime"`
}

func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	var count int64
	tx.Model(&Bookmark{}).Where("user_id = ? AND url = ?", b.UserID, b.URL).Count(&count)
	if count > 0 {
		return ErrAlreadyExists
	}
	return nil
}


type Note struct {
	ID          uuid.UUID	`json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Title       string		`json:"title"`
	Content     string		`json:"content"`
	Tags        tags		`json:"tags" gorm:"many2many:note_tags;"`
	UserID      uuid.UUID	`json:"user_id"`
	CreatedAt   time.Time	`json:"created_at" gorm:"autoCreateTime"`
}

func (n *Note) BeforeCreate(tx *gorm.DB) error {
	var count int64
	tx.Model(&Note{}).Where("user_id = ? AND title = ?", n.UserID, n.Title).Count(&count)
	if count > 0 {
		return ErrAlreadyExists
	}
	return nil
}

type User struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
    Username     string    `json:"username" gorm:"unique;not null"`
    Password     string    `json:"-" gorm:"not null"`
    Role         string    `json:"role"`
    CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
    LastLoginAt  time.Time `json:"last_login_at"`
}

type Tag struct {
	ID    uuid.UUID	`json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Name  string	`json:"name"`
	Count int64		`json:"count" gorm:"default:0"`
}

type tags []Tag

func (t tags) Names() []string {
	names := make([]string, len(t))
	for i, tag := range t {
		names[i] = tag.Name
	}
	return names
}

func (t tags) IDs() []string {
	ids := make([]string, len(t))
	for i, tag := range t {
		ids[i] = tag.ID.String()
	}
	return ids
}