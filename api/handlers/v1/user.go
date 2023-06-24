package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/helper/email"
	token "github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/tokens"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/pkg/etc"
	"github.com/golanguzb70/validator"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// @Router		/user/check/{email} [GET]
// @Summary		Create user
// @Tags        User
// @Description	Here user can be created.
// @Accept      json
// @Produce		json
// @Param       email       path     string true "email"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserCheck(ctx *gin.Context) {
	var (
		emailP = ctx.Param("email")
	)

	ctxTimout, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
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

// @Router		/user/otp [GET]
// @Summary		Check Otp
// @Tags        User
// @Description	Here otp can be checked if true.
// @Accept      json
// @Produce		json
// @Param       email       query     string true "email"
// @Param       otp       query     string true "otp"
// @Success		200 	{object}  models.OtpCheckResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) OtpCheck(ctx *gin.Context) {
	var (
		body   models.Otp
		emailP = ctx.Query("email")
		otpP   = ctx.Query("otp")
	)

	otpAny, err := h.redis.Get(emailP)
	if HandleInternalWithMessage(ctx, h.log, err, "OtpCheck.h.redis.Get()") {
		return
	}

	if otpAny == "" {
		if HandleBadRequestErrWithMessage(ctx, h.log, fmt.Errorf("otp is not found or expired"), "OtpCheck.h.redis.Get() Empty") {
			return
		}
	}

	err = json.Unmarshal([]byte(cast.ToString(otpAny)), &body)
	if HandleInternalWithMessage(ctx, h.log, err, "OtpCheck.json.Unmarshal()") {
		return
	}

	ctx.JSON(http.StatusOK, models.OtpCheckResponse{
		ErrorCode: ErrorSuccessCode,
		Body: struct {
			IsRight bool "json:\"is_right\""
		}{IsRight: body.Code == otpP},
	})
}

// @Router		/user [POST]
// @Summary		Register user
// @Tags        User
// @Description	Here user can be registered.
// @Accept      json
// @Produce		json
// @Param       post   body       models.UserRegisterReq true "post info"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UserRegister(ctx *gin.Context) {
	var (
		res = &models.UserResponse{}
	)
	body := &models.UserRegisterReq{}
	otpBody := models.Otp{}

	err := ctx.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(ctx, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	otpAny, err := h.redis.Get(body.Email)
	if HandleInternalWithMessage(ctx, h.log, err, "OtpCheck.h.redis.Get()") {
		return
	}

	if cast.ToString(otpAny) == "" {
		ctx.JSON(http.StatusBadRequest, models.DefaultResponse{
			ErrorCode:    ErrorCodeOtpIncorrect,
			ErrorMessage: "Otp not found",
		})
		return
	}

	err = json.Unmarshal([]byte(cast.ToString(otpAny)), &otpBody)
	if HandleInternalWithMessage(ctx, h.log, err, "OtpCheck.json.Unmarshal()") {
		return
	}

	if otpBody.Code != body.Otp {
		ctx.JSON(http.StatusBadRequest, models.DefaultResponse{
			ErrorCode:    ErrorCodeOtpIncorrect,
			ErrorMessage: "Otp incorrect",
		})
		return
	}

	req := &models.UserCreateReq{}
	err = StructToStruct(body, &req)
	if HandleInternalWithMessage(ctx, h.log, err, "UserRegister.StructToStruct(body, &req)") {
		return
	}
	req.Id = uuid.New().String()

	req.Password, err = etc.HashPassword(req.Password)
	if HandleInternalWithMessage(ctx, h.log, err, "UserRegister.etc.HashPassword(req.Password)") {
		return
	}

	// Create access and refresh tokens JWT
	h.jwthandler = token.JWTHandler{
		Sub:       req.Id,
		Role:      "user",
		SigninKey: h.cfg.SignInKey,
		Aud:       []string{"template-front"},
		Log:       h.log,
	}
	access, refresh, err := h.jwthandler.GenerateAuthJWT()
	if HandleInternalWithMessage(ctx, h.log, err, "UserRegister.h.jwthandler.GenerateAuthJWT()") {
		return
	}

	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	req.RefreshToken = refresh
	res, err = h.storage.Postgres().UserCreate(ctxWithCancel, req)
	if HandleDatabaseLevelWithMessage(ctx, h.log, err, "UserRegister.h.storage.Postgres().UserCreate()") {
		return
	}
	res.AccessToken = access
	ctx.JSON(http.StatusOK, &models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router 			/user/login	[POST]
// @Summary 		User Login
// @Description  	Through this api user is logged in
// @Tags 			User
// @Accept 			json
// @Produce 		json
// @Param 	user 	body 	 	models.UserLoginRequest true "User Login"
// @Success 200 	{object} 	models.UserApiResponse
// @Failure default {object}  	models.DefaultResponse
func (h *handlerV1) LoginUser(ctx *gin.Context) {
	var (
		body models.UserLoginRequest
		req  models.UserGetReq
	)
	err := ctx.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(ctx, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	isEmail := validator.IsEmail(body.UserNameOrEmail)
	if isEmail {
		req.Email = body.UserNameOrEmail
	} else {
		req.UserName = body.UserNameOrEmail
	}
	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	res, err := h.storage.Postgres().UserGet(ctxWithCancel, &req)
	if HandleBadRequestErrWithMessage(ctx, h.log, err, "LoginUser:h.storage.Postgres().UserGet()") {
		return
	}

	if !etc.CheckPasswordHash(body.Password, res.Password) {
		ctx.JSON(http.StatusConflict, models.DefaultResponse{
			ErrorCode:    ErrorCodeWrongPassword,
			ErrorMessage: "Your password is incorrect",
		})
		return
	}

	h.jwthandler = token.JWTHandler{
		Sub:       req.Id,
		Role:      "user",
		SigninKey: h.cfg.SignInKey,
		Aud:       []string{"template-front"},
		Log:       h.log,
	}
	access, _, err := h.jwthandler.GenerateAuthJWT()
	if HandleInternalWithMessage(ctx, h.log, err, "UserRegister.h.jwthandler.GenerateAuthJWT()") {
		return
	}

	res.AccessToken = access
	res.Password = ""
	ctx.JSON(http.StatusOK, &models.UserApiResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "",
		Body:         res,
	})
}

// @Router 			/user/forgot-password/{user_name_or_email}	[GET]
// @Summary 		User forgot password
// @Description  	Through this api user forgot  password can be enabled.
// @Tags 			User
// @Accept 			json
// @Produce 		json
// @Param       	user_name_or_email       path     string true "user_name_or_email"
// @Success 200 	{object} 	models.DefaultResponse
// @Failure default {object}  	models.DefaultResponse
func (h *handlerV1) UserForgotPassword(ctx *gin.Context) {
	var (
		req models.UserGetReq
	)

	uOre := ctx.Param("user_name_or_email")

	isEmail := validator.IsEmail(uOre)
	if isEmail {
		req.Email = uOre
	} else {
		req.UserName = uOre
	}
	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	res, err := h.storage.Postgres().UserGet(ctxWithCancel, &req)
	if HandleBadRequestErrWithMessage(ctx, h.log, err, "LoginUser:h.storage.Postgres().UserGet()") {
		return
	}

	otp := &models.Otp{
		Email: res.Email,
		Code:  etc.GenerateCode(6),
	}

	// save to redis
	otpByte, err := json.Marshal(otp)
	if HandleInternalWithMessage(ctx, h.log, err, "UserCheck.json.Marshal()") {
		return
	}

	err = h.redis.SetWithTTL(res.Email, string(otpByte), int(h.cfg.OtpTimeout))
	if HandleInternalWithMessage(ctx, h.log, err, "UserCheck.h.redis.SetWithTTL()") {
		return
	}

	// send otp email
	err = email.SendEmail([]string{res.Email}, "GolangUzb70\n", h.cfg, "./api/helper/email/forgotpassword.html", otp)
	if HandleInternalWithMessage(ctx, h.log, err, "UserCheck.email.SendEmail()") {
		return
	}

	ctx.JSON(http.StatusOK, &models.DefaultResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "We have sent otp to your email  address.",
	})
}

// @Router 			/user/forgot-password/verify	[POST]
// @Summary 		User forgot password
// @Description  	Through this api user forgot  password can be enabled.
// @Tags 			User
// @Accept 			json
// @Produce 		json
// @Param 	user 	body 	 	models.UserForgotPasswordVerifyReq true "User Login"
// @Success 200 	{object} 	models.DefaultResponse
// @Failure default {object}  	models.DefaultResponse
func (h *handlerV1) UserForgotPasswordVerify(ctx *gin.Context) {
	var (
		body    models.UserForgotPasswordVerifyReq
		req     models.UserGetReq
		otpBody models.Otp
	)

	err := ctx.ShouldBindJSON(&body)
	if HandleBadRequestErrWithMessage(ctx, h.log, err, "c.ShouldBindJSON(&body)") {
		return
	}

	isEmail := validator.IsEmail(body.UserNameOrEmail)
	if isEmail {
		req.Email = body.UserNameOrEmail
	} else {
		req.UserName = body.UserNameOrEmail
	}
	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	res, err := h.storage.Postgres().UserGet(ctxWithCancel, &req)
	if HandleBadRequestErrWithMessage(ctx, h.log, err, "UserForgotPasswordVerify:h.storage.Postgres().UserGet()") {
		return
	}

	otpAny, err := h.redis.Get(res.Email)
	if HandleInternalWithMessage(ctx, h.log, err, "UserForgotPasswordVerify.h.redis.Get()") {
		return
	}

	if cast.ToString(otpAny) == "" {
		ctx.JSON(http.StatusBadRequest, models.DefaultResponse{
			ErrorCode:    ErrorCodeOtpIncorrect,
			ErrorMessage: "Otp not found",
		})
		return
	}

	err = json.Unmarshal([]byte(cast.ToString(otpAny)), &otpBody)
	if HandleInternalWithMessage(ctx, h.log, err, "UserForgotPasswordVerify.json.Unmarshal()") {
		return
	}

	if otpBody.Code != body.Otp {
		ctx.JSON(http.StatusBadRequest, models.DefaultResponse{
			ErrorCode:    ErrorCodeOtpIncorrect,
			ErrorMessage: "Otp incorrect",
		})
		return
	}
	res.Password, err = etc.HashPassword(body.NewPassword)
	if HandleInternalWithMessage(ctx, h.log, err, "UserForgotPasswordVerify:etc.HashPassword(res.Password)") {
		return
	}

	err = h.storage.Postgres().UpdateSingleField(ctxWithCancel, &models.UpdateSingleFieldReq{
		Id:       res.Id,
		Table:    "users",
		Column:   "hashed_password",
		NewValue: res.Password,
	})
	if HandleDatabaseLevelWithMessage(ctx, h.log, err, "UserForgotPasswordVerify:h.storage.Postgres().UpdateSingleField()") {
		return
	}

	ctx.JSON(http.StatusOK, &models.DefaultResponse{
		ErrorCode:    ErrorSuccessCode,
		ErrorMessage: "Password successfully updated",
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
	// id, err := strconv.Atoi(c.Param("id"))
	// if HandleBadRequestErrWithMessage(c, h.log, err, "strconv.Atoi()") {
	// 	return
	// }

	res, err := h.storage.Postgres().UserGet(context.Background(), &models.UserGetReq{
		// Id: id,
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
	// id, err := strconv.Atoi(c.Param("id"))
	// if HandleBadRequestErrWithMessage(c, h.log, err, "strconv.Atoi()") {
	// 	return
	// }

	// err = h.storage.Postgres().UserDelete(context.Background(), &models.UserDeleteReq{Id: id})
	// if HandleDatabaseLevelWithMessage(c, h.log, err, "h.storage.Postgres().UserDelete()") {
	// 	return
	// }

	// c.JSON(http.StatusOK, models.DefaultResponse{
	// 	ErrorCode:    ErrorSuccessCode,
	// 	ErrorMessage: "Successfully deleted",
	// })

}
