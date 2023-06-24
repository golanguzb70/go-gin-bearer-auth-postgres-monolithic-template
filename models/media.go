package models

import "mime/multipart"

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadPhotoRes struct {
	URL string `json:"photo_url"`
}

type MediaResponse struct {
	ErrorCode    int            `json:"error_code"`
	ErrorMessage string         `json:"error_message"`
	Body         UploadPhotoRes `json:"body"`
}
