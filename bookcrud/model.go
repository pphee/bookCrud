package model

import "time"

type Book struct {
	Book      string `json:"text" binding:"required"`
	Author    string `json:"author" binding:"required"`
	Title     string `json:"title" binding:"required"`
	ID        int    `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
