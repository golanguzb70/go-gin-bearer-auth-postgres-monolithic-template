package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
)

// @Router		/template [POST]
// @Summary		Create template
// @Tags        Template
// @Description	Here template can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.TemplateCreateReq true "post info"
// @Success		200 	{object}  models.TemplateResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) TemplateCreate(ctx *gin.Context) {
	body := &models.TemplateCreateReq{}
	err := ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	res, err := h.storage.Postgres().TemplateCreate(context.Background(), body)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "TemplateCreate: h.storage.Postgres().TemplateCreate()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/template/{id} [GET]
// @Summary		Get template by key
// @Tags        Template
// @Description	Here template can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     int true "id"
// @Success		200 	{object}  models.TemplateResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) TemplateGet(ctx *gin.Context) {
	res, err := h.storage.Postgres().TemplateGet(context.Background(), &models.TemplateGetReq{
		Id: ctx.Param("id"),
	})
	
	if h.HandleDatabaseLevelWithMessage(ctx, err, "TemplateGet: h.storage.Postgres().TemplateGet()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/template/list [GET]
// @Summary		Get templates list
// @Tags        Template
// @Description	Here all templates can be got.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       filters query models.TemplateFindReq true "filters"
// @Success		200 	{object}  models.TemplateFindResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) TemplateFind(ctx *gin.Context) {
	var (
		dbReq = &models.TemplateFindReq{}
		err   error
	)

	dbReq.Page, err = ParsePageQueryParam(ctx)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid page param", nil) {
		return
	}

	dbReq.Limit, err = ParseLimitQueryParam(ctx)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid limit param", nil) {
		return
	}

	dbReq.Search = ctx.Query("search")
	dbReq.OrderByCreatedAt, _ = strconv.ParseUint(ctx.Query("order_by_created_at"), 10, 8)

	res, err := h.storage.Postgres().TemplateFind(context.Background(), dbReq)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "TemplateFind: h.storage.Postgres().TemplateFind()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/template [PUT]
// @Summary		Update template
// @Tags        Template
// @Description	Here template can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.TemplateUpdateReq true "post info"
// @Success		200 	{object}  models.TemplateResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) TemplateUpdate(ctx *gin.Context) {
	body := &models.TemplateUpdateReq{}
	err := ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	res, err := h.storage.Postgres().TemplateUpdate(context.Background(), body)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "TemplateUpdate: h.storage.Postgres().TemplateUpdate()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/template/{id} [DELETE]
// @Summary		Delete template
// @Tags        Template
// @Description	Here template can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     int true "id"
// @Success		200 	{object}  models.StandardResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) TemplateDelete(ctx *gin.Context) {
	err := h.storage.Postgres().TemplateDelete(context.Background(), &models.TemplateDeleteReq{Id: ctx.Param("id")})
	if h.HandleDatabaseLevelWithMessage(ctx, err, "TemplateDelete: h.storage.Postgres().TemplateDelete()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "Successfully deleted", nil)
}
