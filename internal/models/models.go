package models

import "time"

// User is the user model
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Room is the room model
type Room struct {
	ID        string
	UserID    int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
