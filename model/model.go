package model

import (
	"time"

	"github.com/google/uuid"
)

type Bookmark struct {
	ID         	uuid.UUID	`json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	URL 	   	string		`json:"url" gorm:"not null"`
	Title      	string		`json:"title"`
	Description string		`json:"description"`
	Tags       	tags		`json:"tags" gorm:"many2many:bookmark_tags;"`
	UserID     	uuid.UUID	`json:"user_id" gorm:"not null"`
	CreatedAt  	time.Time	`json:"created_at" gorm:"autoCreateTime"`
}

type Note struct {
	ID          uuid.UUID	`json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	Title       string		`json:"title"`
	Content     string		`json:"content"`
	Tags        tags		`json:"tags" gorm:"many2many:note_tags;"`
	UserID      uuid.UUID	`json:"user_id"`
	CreatedAt   time.Time	`json:"created_at" gorm:"autoCreateTime"`
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

type Tag struct {
	ID    uuid.UUID	`json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
	Name  string	`json:"name"`
	Count int64		`json:"count" gorm:"default:0"`
}

type User struct {
    ID           uuid.UUID `json:"id" gorm:"primaryKey;default:uuid_generate_v4()"`
    Username     string    `json:"username"`
    Password     string    `json:"-" gorm:"not null"`
	RefreshToken string
    Role         string    `json:"role"`
    CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
    LastLoginAt  time.Time `json:"last_login_at"`
}