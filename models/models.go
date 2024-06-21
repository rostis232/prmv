package models

import (
	"time"
)

type Post struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title" validate:"required,min=3,max=100"`
	Content   string    `db:"content" json:"content" validate:"required,min=3"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
