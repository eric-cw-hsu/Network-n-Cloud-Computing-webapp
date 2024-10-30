package domain

import "time"

type User struct {
	ID        string
	UpdatedAt time.Time
}

func NewUser(ID string) *User {
	return &User{
		ID:        ID,
		UpdatedAt: time.Now(),
	}
}
