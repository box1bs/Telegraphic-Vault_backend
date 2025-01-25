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
	Tags       	[]Tag		`json:"tags" gorm:"many2many:bookmark_tags;constraint:OnDelete:CASCADE;"`
	UserID     	uuid.UUID	`json:"-" gorm:"not null"`
	CreatedAt  	time.Time	`json:"created_at" gorm:"autoCreateTime"`
}

func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	var exist Bookmark
	if err := tx.Model(&Bookmark{}).Where("user_id = ? AND url = ?", b.UserID, b.URL).First(&exist).Error; err == nil {
		return ErrAlreadyExists
	}
	return nil
}

type Note struct {
	ID          	uuid.UUID	`json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Title       	string		`json:"title"`
	Content     	string		`json:"content"`
	Tags        	[]Tag		`json:"tags" gorm:"many2many:note_tags;constraint:OnDelete:CASCADE;"`
	UserID      	uuid.UUID	`json:"-"`
	CreatedAt   	time.Time	`json:"created_at" gorm:"autoCreateTime"`
}

func (n *Note) BeforeCreate(tx *gorm.DB) error {
	var exist Note
	if err := tx.Model(&Note{}).Where("user_id = ? AND title = ?", n.UserID, n.Title).First(&exist).Error; err == nil {
		return ErrAlreadyExists
	}
	return nil
}

type User struct {
	ID           uuid.UUID `json:"-" gorm:"primaryKey;default:gen_random_uuid()"`
    Username     string    `json:"username" gorm:"unique;not null"`
    Password     string    `json:"-" gorm:"not null"`
    Role         string    `json:"role"`
    CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
    LastLoginAt  time.Time `json:"last_login_at"`
}

type Tag struct {
	ID    uuid.UUID	`json:"id,omitempty" gorm:"primaryKey;default:gen_random_uuid()"`
	Name  string	`json:"name"`
	Count int64		`json:"count,omitempty" gorm:"default:0"`
}