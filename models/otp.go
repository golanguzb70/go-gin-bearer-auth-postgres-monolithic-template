package models

type Otp struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type OtpCheckResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Body         struct {
		IsRight bool `json:"is_right"`
	} `json:"body"`
}
