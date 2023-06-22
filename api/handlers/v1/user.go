package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/helper/email"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/pkg/etc"
	"github.com/google/uuid"
)

// @Router		/user/check/{email} [GET]
// @Summary		Create user
// @Tags        User
// @Description	Here user can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       email       path     string true "email"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserCheck(ctx *gin.Context) {
	var (
		emailP = ctx.Param("email")
	)

	ctxTimout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContectTimeout))
	defer cancel()

	exists, err := h.storage.Postgres().CheckIfExists(ctxTimout, &models.CheckIfExistsReq{
		Table:  "users",
		Column: "email",
		Value:  emailP,
	})
	if HandleDatabaseLevelWithMessage(ctx, h.log, err, "UserCheck.Postgres().CheckIfExists()") {
		return
	}

	if exists.Exists {
		ctx.JSON(http.StatusOK, models.UserCheckResponse{
			ErrorCode:    ErrorSuccessCode,
			ErrorMessage: "",
			Body: &models.UserCheckRes{
				Status: "login",
			},
		})
		return
	}

	otp := &models.Otp{
		Email: emailP,
		Code:  etc.GenerateCode(6),
	}
	// save to redis
	otpByte, err := json.Marshal(otp)
	if HandleInternalWithMessage(ctx, h.log, err, "UserCheck.json.Marshal()") {
		return
	}
	err = h.redis.SetWithTTL(emailP, string(otpByte), int(h.cfg.OtpTimeout))
	if HandleInternalWithMessage(ctx, h.log, err, "UserCheck.h.redis.SetWithTTL()") {
		return
	}

	// send otp email
	err = email.SendEmail([]string{emailP}, "GolangUzb70\n", h.cfg, "./api/helper/email/emailotp.html", otp)
	if HandleInternalWithMessage(ctx, h.log, err, "UserCheck.email.SendEmail()") {
		return
	}

	ctx.JSON(http.StatusOK, models.UserCheckResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body: &models.UserCheckRes{
			Status: "register",
		},
	})
}

// @Router		/user [POST]
// @Summary		Create user
// @Tags        User
// @Description	Here user can be created.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.UserRegisterReq true "post info"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserRegister(c *gin.Context) {
	body := &models.UserRegisterReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	req := &models.UserCreateReq{}
	err = StructToStruct(body, &req)
	if HandleInternalWithMessage(c, h.log, err, "UserRegister.StructToStruct(body, &req)") {
		return
	}
	req.Id = uuid.New().String()

	req.Password, err = etc.HashPassword(req.Password)
	if HandleInternalWithMessage(c, h.log, err, "UserRegister.etc.HashPassword(req.Password)") {
		return
	}

	// save to redis
	userbyte, err := json.Marshal(models.Otp{
		Email: req.Email,
		Code:  etc.GenerateCode(6),
	})
	err = h.redis.SetWithTTL(req.Email, string(userbyte), int(h.cfg.OtpTimeout))

	res, err := h.storage.Postgres().UserCreate(context.Background(), req)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "h.storage.Postgres().UserCreate()") {
		return
	}

	c.JSON(http.StatusOK, &models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/user/{id} [GET]
// @Summary		Get user by key
// @Tags        User
// @Description	Here user can be got.
// @Accept      json
// @Produce		json
// @Param       id       path     int true "id"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if HandleBadRequestErrWithMessage(c, h.log, err, "strconv.Atoi()") {
		return
	}

	res, err := h.storage.Postgres().UserGet(context.Background(), &models.UserGetReq{
		Id: id,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "h.storage.Postgres().UserGet()") {
		return
	}

	c.JSON(http.StatusOK, models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/user/list [GET]
// @Summary		Get users list
// @Tags        User
// @Description	Here all users can be got.
// @Accept      json
// @Produce		json
// @Param       filters query models.UserFindReq true "filters"
// @Success		200 	{object}  models.UserApiFindResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserFind(c *gin.Context) {
	page, err := ParsePageQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "helper.ParsePageQueryParam(c)") {
		return
	}
	limit, err := ParseLimitQueryParam(c)
	if HandleBadRequestErrWithMessage(c, h.log, err, "helper.ParseLimitQueryParam(c)") {
		return
	}

	res, err := h.storage.Postgres().UserFind(context.Background(), &models.UserFindReq{
		Page:  page,
		Limit: limit,
	})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "h.storage.Postgres().UserFind()") {
		return
	}

	c.JSON(http.StatusOK, &models.UserApiFindResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Summary		Update user
// @Tags        User
// @Description	Here user can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.UserUpdateReq true "post info"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
// @Router		/user [PUT]
func (h *handlerV1) UserUpdate(c *gin.Context) {
	body := &models.UserUpdateReq{}
	err := c.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(c, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	res, err := h.storage.Postgres().UserUpdate(context.Background(), body)
	if HandleDatabaseLevelWithMessage(c, h.log, err, "h.storage.Postgres().UserUpdate()") {
		return
	}

	c.JSON(http.StatusOK, &models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router		/user/{id} [DELETE]
// @Summary		Delete user
// @Tags        User
// @Description	Here user can be deleted.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       id       path     int true "id"
// @Success		200 	{object}  models.DefaultResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if HandleBadRequestErrWithMessage(c, h.log, err, "strconv.Atoi()") {
		return
	}

	err = h.storage.Postgres().UserDelete(context.Background(), &models.UserDeleteReq{Id: id})
	if HandleDatabaseLevelWithMessage(c, h.log, err, "h.storage.Postgres().UserDelete()") {
		return
	}

	c.JSON(http.StatusOK, models.DefaultResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "Successfully deleted",
	})
}
