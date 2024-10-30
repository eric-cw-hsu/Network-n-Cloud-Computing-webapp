package domain

import (
	"database/sql"
	"time"
)

type ProfilePic struct {
	Filename   string
	UploadedAt sql.NullTime
}

func NewProfilePic(filename string) *ProfilePic {
	return &ProfilePic{
		Filename:   filename,
		UploadedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
}
