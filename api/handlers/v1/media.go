package v1

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
	"github.com/google/uuid"
)

// Upload photo
// @Summary 		Upload media
// @Description 	Through this api frontent can upload photo and get the link to the media.
// @Tags 			Media
// @Security        BearerAuth
// @Accept 			multipart/form-data
// @Produce         json
// @Param         	file                  formData file true "File"
// @Success         200					  {object} 	models.MediaResponse
// @Failure         default               {object}  models.StandardResponse
// @Router          /media/photo [post]
func (h *handlerV1) UploadMedia(ctx *gin.Context) {
	file := &models.File{}
	err := ctx.ShouldBind(&file)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	fileSize := file.File.Size
	if fileSize > (int64(h.cfg.MaxImageSize) << 20) {
		h.HandleResponse(ctx, fmt.Errorf(SizeExceeded), http.StatusRequestEntityTooLarge, SizeExceeded, fmt.Sprintf("Image size should be less than %d mb", h.cfg.MaxImageSize), nil)
		return
	}

	ext := filepath.Ext(file.File.Filename)
	if ext != ".jpg" && ext != ".png" {
		h.HandleResponse(ctx, fmt.Errorf(SizeExceeded), http.StatusBadRequest, BadRequest, "Only .jpg and .png images are allowed", nil)
		return
	}

	file.File.Filename = uuid.New().String() + filepath.Ext(file.File.Filename)

	err = ctx.SaveUploadedFile(file.File, "./media/"+file.File.Filename)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UploadMedia:SaveUploadedFile", nil) {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", models.UploadPhotoRes{
		URL: h.cfg.BaseUrl + "media/" + file.File.Filename,
	})
}
