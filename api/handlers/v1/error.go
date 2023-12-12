package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handlerV1) HandleDatabaseLevelWithMessage(c *gin.Context, err error, message string, args ...interface{}) bool {
	status_err, _ := status.FromError(err)
	if err != nil {
		errorCode := InternalServerError
		statuscode := http.StatusInternalServerError
		message := status_err.Message()
		switch status_err.Code() {
		case codes.NotFound:
			errorCode = NotFound
			statuscode = http.StatusNotFound
		case codes.Unknown:
			errorCode = InternalServerError
			statuscode = http.StatusBadRequest
			message = "Ooops something went wrong"
		case codes.Aborted:
			errorCode = BadRequest
			statuscode = http.StatusBadRequest
		case codes.InvalidArgument:
			errorCode = BadRequest
			statuscode = http.StatusBadRequest
		}

		h.log.Error(message, err, args)
		c.AbortWithStatusJSON(statuscode, models.StandardResponse{
			Status:  errorCode,
			Message: message,
		})
		return true
	}
	return false
}

// Handles response according to err arguments. If err is nil it returns false otherwise true
func (h *handlerV1) HandleResponse(c *gin.Context, err error, httpStatusCode int, status, message string, data any, args ...any) bool {
	if err != nil {
		if status != InternalServerError {
			c.AbortWithStatusJSON(httpStatusCode, models.StandardResponse{
				Status:  status,
				Message: message,
				Data:    data,
			})
		} else {
			h.log.Error(message, err, args)
			c.AbortWithStatusJSON(httpStatusCode, models.StandardResponse{
				Status:  status,
				Message: "Internal server error",
				Data:    data,
			})
		}
		return true
	} else if status == "success" {
		c.JSON(httpStatusCode, models.StandardResponse{
			Status:  status,
			Message: message,
			Data:    data,
		})
	}

	return false
}
