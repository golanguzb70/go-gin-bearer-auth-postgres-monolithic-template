package v1

const (
	// 200 ... 2xx
	Success = "success"

	// 400
	BadRequest    = "bad_request"
	WrongPassword = "wrong_password"

	// 401
	UnAuthorized        = "unauthorized"
	AccessTokenExpired  = "access_token_expired"
	RefreshTokenExpired = "refresh_token_expired"
	HeaderRequired      = "header_required" // this can be custom authorization headers like session, access key...

	// 403
	PermissionDenied = "permission_denied"

	// 404
	NotFound = "not_found"

	// 413
	SizeExceeded = "size_exceeded"

	// 500
	InternalServerError = "internal_server_error"
)