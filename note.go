package main

type Note struct {
	ID uint64 `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}