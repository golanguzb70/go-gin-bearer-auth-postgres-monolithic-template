package models

type UserCheckRes struct {
	Status string `json:"status"`
}

type UserCheckResponse struct {
	ErrorCode    int           `json:"error_code"`
	ErrorMessage string        `json:"error_message"`
	Body         *UserCheckRes `json:"body"`
}

type UserRegisterReq struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

type UserCreateReq struct {
	Id           string `json:"id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginRequest struct {
	UserNameOrEmail string `json:"user_name_or_email"`
	Password        string `json:"password"`
}

type UserUpdateReq struct {
	Id       string `json:"id"`
	UserName string `json:"user_name"`
}

type UserApiUpdateReq struct {
	UserName string `json:"user_name"`
}

type UserGetReq struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
}

type UserFindReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type UserDeleteReq struct {
	Id string `json:"id"`
}

type UserFindResponse struct {
	Users []*UserResponse `json:"users"`
	Count int             `json:"count"`
}

type UserResponse struct {
	Id           string `json:"id"`
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UserForgotPasswordVerifyReq struct {
	NewPassword     string `json:"new_password"`
	Otp             string `json:"otp"`
	UserNameOrEmail string `json:"user_name_or_email"`
}
