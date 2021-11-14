package main

type Note struct {
	ID uint `gorm:"primaryKey"`
	Title string
	Content string
}