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
// @Summary		Check User status
// @Tags        User Authorzation
// @Description	Here user status is checked. If user is exists in database it should be logged in else registered
// @Accept      json
// @Produce		json
// @Param       email       path     string true "email"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.StandardResponse
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
	if h.HandleDatabaseLevelWithMessage(ctx, err, "UserCheck: Postgres().CheckIfExists()") {
		return
	}

	if exists.Exists {
		h.HandleResponse(ctx, nil, http.StatusOK, Success, "", &models.UserCheckRes{
			Status: "login",
		})
		return
	}

	otp := &models.Otp{
		Email: emailP,
		Code:  etc.GenerateCode(6),
	}
	// save to redis
	otpByte, err := json.Marshal(otp)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserCheck: json.Marshal(otp)", nil) {
		return
	}

	err = h.redis.SetWithTTL(emailP, string(otpByte), int(h.cfg.OtpTimeout))
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserCheck: redis.SetWithTTL()", nil) {
		return
	}

	// send otp email
	err = email.SendEmail([]string{emailP}, "GolangUzb70\n", h.cfg, "./api/helper/email/emailotp.html", otp)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserCheck: email.SendEmail()", nil) {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", &models.UserCheckRes{
		Status: "register",
	})
}

// @Router		/user/otp [GET]
// @Summary		Check Otp
// @Tags        User Authorzation
// @Description	Here otp can be checked if true.
// @Accept      json
// @Produce		json
// @Param       email       query     string true "email"
// @Param       otp       query     string true "otp"
// @Success		200 	{object}  models.OtpCheckResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) OtpCheck(ctx *gin.Context) {
	var (
		body   models.Otp
		emailP = ctx.Query("email")
		otp    = ctx.Query("otp")
	)

	otpAny, err := h.redis.Get(emailP)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "OtpCheck: json.redis.Get()", nil) {
		return
	}

	if otpAny == "" {
		if h.HandleResponse(ctx, fmt.Errorf(NotFound), http.StatusBadRequest, NotFound, "otp expired", nil) {
			return
		}
	}

	err = json.Unmarshal([]byte(cast.ToString(otpAny)), &body)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "OtpCheck: json.Unmarshal()", nil) {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", struct {
		IsRight bool `json:"is_right"`
	}{IsRight: body.Code == otp})
}

// @Router		/user [POST]
// @Summary		Register user
// @Tags        User Authorzation
// @Description	Here user can be registered.
// @Accept      json
// @Produce		json
// @Param       post   body       models.UserRegisterReq true "post info"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) UserRegister(ctx *gin.Context) {
	var (
		res = &models.UserResponse{}
	)
	body := &models.UserRegisterReq{}
	otpBody := models.Otp{}

	err := ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	otpAny, err := h.redis.Get(body.Email)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserRegister: redis.Get()", nil) {
		return
	}

	if cast.ToString(otpAny) == "" {
		h.HandleResponse(ctx, fmt.Errorf(BadRequest), http.StatusBadRequest, BadRequest, "otp expired", nil)
		return
	}

	err = json.Unmarshal([]byte(cast.ToString(otpAny)), &otpBody)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserRegister: json.Unmarshal()", nil) {
		return
	}

	if otpBody.Code != body.Otp {
		h.HandleResponse(ctx, fmt.Errorf(BadRequest), http.StatusBadRequest, BadRequest, "otp incorrect", nil)
		return
	}

	req := &models.UserCreateReq{}
	err = StructToStruct(body, &req)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserRegister: StructToStruct(body, &req)", nil) {
		return
	}
	req.Id = uuid.New().String()

	req.Password, err = etc.HashPassword(req.Password)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserRegister: etc.HashPassword(req.Password)", nil) {
		return
	}

	// Create access and refresh tokens JWT
	h.jwthandler = token.JWTHandler{
		Sub:       req.Id,
		Role:      "user",
		SigninKey: h.cfg.SignInKey,
		Aud:       []string{"template-front"},
		Log:       h.log,
		Timout:    h.cfg.AccessTokenTimout,
	}
	access, refresh, err := h.jwthandler.GenerateAuthJWT()
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserRegister: jwthandler.GenerateAuthJWT()", nil) {
		return
	}

	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	req.RefreshToken = refresh
	res, err = h.storage.Postgres().UserCreate(ctxWithCancel, req)
	if h.HandleDatabaseLevelWithMessage(ctx, err, "UserRegister: h.storage.Postgres().UserCreate()") {
		return
	}

	res.AccessToken = access
	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router 			/user/login	[POST]
// @Summary 		User Login
// @Description  	Through this api user is logged in
// @Tags       		User Authorzation
// @Accept 			json
// @Produce 		json
// @Param 	user 	body 	 	models.UserLoginRequest true "User Login"
// @Success 200 	{object} 	models.UserApiResponse
// @Failure default {object}  	models.StandardResponse
func (h *handlerV1) LoginUser(ctx *gin.Context) {
	var (
		body models.UserLoginRequest
		req  models.UserGetReq
	)
	err := ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
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
	if h.HandleDatabaseLevelWithMessage(ctx, err, "LoginUser:h.storage.Postgres().UserGet()") {
		return
	}

	if !etc.CheckPasswordHash(body.Password, res.Password) {
		h.HandleResponse(ctx, fmt.Errorf(BadRequest), http.StatusBadRequest, BadRequest, "incorrect password", nil)
		return
	}

	h.jwthandler = token.JWTHandler{
		Sub:       res.Id,
		Role:      "user",
		SigninKey: h.cfg.SignInKey,
		Aud:       []string{"template-front"},
		Log:       h.log,
		Timout:    h.cfg.AccessTokenTimout,
	}
	access, _, err := h.jwthandler.GenerateAuthJWT()
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "LoginUser: jwthandler.GenerateAuthJWT()", nil) {
		return
	}

	res.AccessToken = access
	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router 			/user/forgot-password/{user_name_or_email}	[GET]
// @Summary 		User forgot password
// @Description  	Through this api user forgot  password can be enabled.
// @Tags            User Authorzation
// @Accept 			json
// @Produce 		json
// @Param       	user_name_or_email       path     string true "user_name_or_email"
// @Success 200 	{object} 	models.StandardResponse
// @Failure default {object}  	models.StandardResponse
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
	if h.HandleDatabaseLevelWithMessage(ctx, err, "UserForgotPassword: storage.Postgres().UserGet()") {
		return
	}

	otp := &models.Otp{
		Email: res.Email,
		Code:  etc.GenerateCode(6),
	}

	// save to redis
	otpByte, err := json.Marshal(otp)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserForgotPassword: json.Marshal(otp)", nil) {
		return
	}

	err = h.redis.SetWithTTL(res.Email, string(otpByte), int(h.cfg.OtpTimeout))
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserForgotPassword: redis.SetWithTTL()", nil) {
		return
	}

	// send otp email
	err = email.SendEmail([]string{res.Email}, "GolangUzb70\n", h.cfg, "./api/helper/email/forgotpassword.html", otp)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserForgotPassword: email.SendEmail()", nil) {
		return
	}

	h.HandleResponse(ctx, err, http.StatusOK, Success, "We have sent otp to your email  address.", nil)
}

// @Router 	        /user/forgot-password/verify [POST]
// @Summary 		User forgot password
// @Description  	Through this api user forgot  password can be enabled.
// @Tags 			User Authorzation
// @Accept 			json
// @Produce 		json
// @Param 	user 	body 	 	models.UserForgotPasswordVerifyReq true "User Login"
// @Success 200 	{object} 	models.StandardResponse
// @Failure default {object}  	models.StandardResponse
func (h *handlerV1) UserForgotPasswordVerify(ctx *gin.Context) {
	var (
		body    models.UserForgotPasswordVerifyReq
		req     models.UserGetReq
		otpBody models.Otp
	)

	err := ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
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
	if h.HandleDatabaseLevelWithMessage(ctx, err, "UserForgotPasswordVerify: h.storage.Postgres().UserGet()") {
		return
	}

	otpAny, err := h.redis.Get(res.Email)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserForgotPasswordVerify: redis.Get(res.Email)", nil) {
		return
	}

	if cast.ToString(otpAny) == "" {
		h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "otp expired", nil)
		return
	}

	err = json.Unmarshal([]byte(cast.ToString(otpAny)), &otpBody)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserForgotPasswordVerify: json.Unmarshal()", nil) {
		return
	}

	if otpBody.Code != body.Otp {
		h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "otp incorrect", nil)
		return
	}

	res.Password, err = etc.HashPassword(body.NewPassword)
	if h.HandleResponse(ctx, err, http.StatusInternalServerError, InternalServerError, "UserForgotPasswordVerify: etc.HashPassword(body.NewPassword)", nil) {
		return
	}

	err = h.storage.Postgres().UpdateSingleField(ctxWithCancel, &models.UpdateSingleFieldReq{
		Id:       res.Id,
		Table:    "users",
		Column:   "hashed_password",
		NewValue: res.Password,
	})
	if h.HandleDatabaseLevelWithMessage(ctx, err, "UserForgotPasswordVerify:h.storage.Postgres().UpdateSingleField()") {
		return
	}

	h.HandleResponse(ctx, err, http.StatusOK, Success, "Password successfully updated", nil)
}

// @Router		/user/profile [GET]
// @Summary		Get user by key
// @Tags        User
// @Description	Here user profile info can be got by id.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) UserGet(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	res, err := h.storage.Postgres().UserGet(context.Background(), &models.UserGetReq{
		Id: claim.Sub,
	})
	if h.HandleDatabaseLevelWithMessage(ctx, err, "UserGet:h.storage.Postgres().UserGet()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/user [PUT]
// @Summary		Update user
// @Tags        User
// @Description	Here user can be updated.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Param       post   body       models.UserApiUpdateReq true "post info"
// @Success		200 	{object}  models.UserApiResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) UserUpdate(ctx *gin.Context) {
	var (
		body models.UserApiUpdateReq
	)

	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	err = ctx.ShouldBindJSON(&body)
	if h.HandleResponse(ctx, err, http.StatusBadRequest, BadRequest, "invalid body", nil) {
		return
	}

	res, err := h.storage.Postgres().UserUpdate(context.Background(), &models.UserUpdateReq{
		Id:       claim.Sub,
		UserName: body.UserName,
	})
	if h.HandleDatabaseLevelWithMessage(ctx, err, "h.storage.Postgres().UserUpdate()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "", res)
}

// @Router		/user [DELETE]
// @Summary		Delete user
// @Tags        User
// @Description	Here user can be deleted, user_id is taken from token.
// @Security    BearerAuth
// @Accept      json
// @Produce		json
// @Success		200 	{object}  models.StandardResponse
// @Failure     default {object}  models.StandardResponse
func (h *handlerV1) UserDelete(ctx *gin.Context) {
	claim, err := GetClaims(*h, ctx)
	if h.HandleResponse(ctx, err, http.StatusUnauthorized, UnAuthorized, "invalid authorization", nil) {
		return
	}

	err = h.storage.Postgres().UserDelete(context.Background(), &models.UserDeleteReq{Id: claim.Sub})
	if h.HandleDatabaseLevelWithMessage(ctx, err, "UserDelete: h.storage.Postgres().UserDelete()") {
		return
	}

	h.HandleResponse(ctx, nil, http.StatusOK, Success, "Successfully deleted", nil)
}