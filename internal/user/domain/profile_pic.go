package domain

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type ProfilePic struct {
	Filename      string
	UploadedAt    time.Time
	Url           string
	S3Key         string
	ETag          string
	Encryption    string
	EncryptionKey string
}

func NewProfilePic(
	filename, url, s3Key, eTag string,
	encryption types.ServerSideEncryption,
	encryptionKey string,
) *ProfilePic {
	return &ProfilePic{
		Filename:      filename,
		UploadedAt:    time.Now(),
		Url:           url,
		S3Key:         s3Key,
		ETag:          eTag,
		Encryption:    string(encryption),
		EncryptionKey: encryptionKey,
	}
}
