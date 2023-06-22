package models

type UserCreateReq struct {
	UserName string `json:"user_name"`
}

type UserUpdateReq struct {
	Id           int    `json:"id"`
	UserName string `json:"user_name"`
}

type UserGetReq struct {
	Id int `json:"id"`
}

type UserFindReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type UserDeleteReq struct {
	Id int `json:"id"`
}

type UserFindResponse struct {
	Users []*UserResponse `json:"users"`
	Count     int                 `json:"count"`
}

type UserResponse struct {
	Id           int    `json:"id"`
	UserName string `json:"user_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UserApiResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *UserResponse
}

type UserApiFindResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         *UserFindResponse
}
