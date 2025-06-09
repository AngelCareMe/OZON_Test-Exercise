package models

import "time"

type Comment struct {
	ID              int
	PostID          int
	ParentCommentID *int
	Text            string
	Author          string
	CreatedAt       time.Time
}
