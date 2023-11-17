package model

import "time"

type Book struct {
	Book      string    `json:"book" faker:"sentence"`
	Author    string    `json:"author" faker:"name"`
	Title     string    `json:"title" faker:"sentence"`
	ID        int       `gorm:"primarykey"`
	CreatedAt time.Time `faker:"-"`
	UpdatedAt time.Time `faker:"-"`
}
