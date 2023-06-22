package v1

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseLimitQueryParam(c *gin.Context) (int, error) {
	return strconv.Atoi(c.DefaultQuery("limit", "10"))
}

func ParsePageQueryParam(c *gin.Context) (int, error) {
	return strconv.Atoi(c.DefaultQuery("page", "1"))
}

func StructToStruct(from, to any) error {
	body, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &to)
	return err
}
