package dto

import (
	"go-template/internal/user/domain"
)

type PicResponse struct {
	UserID     string `json:"user_id"`
	FileName   string `json:"file_name"`
	URL        string `json:"url"`
	UploadDate string `json:"upload_date"`
}

func NewPicResponse(user *domain.User, profilePic *domain.ProfilePic, bucketName string) *PicResponse {
	return &PicResponse{
		UserID:     user.ID,
		FileName:   profilePic.Filename,
		URL:        bucketName + "/" + user.ID + "/" + profilePic.Filename,
		UploadDate: profilePic.UploadedAt.Format("2006-01-02"),
	}
}
