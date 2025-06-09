package models

import "time"

type Post struct {
	ID            int
	Title         string
	Text          string
	AllowComments bool
	Author        string
	CreatedAt     time.Time
}
