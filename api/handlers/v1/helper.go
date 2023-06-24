package v1

import (
	"encoding/json"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	t "github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/tokens"
	"github.com/spf13/cast"
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

func GetClaims(h handlerV1, c *gin.Context) (*t.CustomClaims, error) {
	var (
		claims = t.CustomClaims{}
	)
	strToken := c.GetHeader("Authorization")

	token, err := jwt.Parse(strToken, func(t *jwt.Token) (interface{}, error) { return []byte(h.cfg.SignInKey), nil })

	if err != nil {
		h.log.Error("invalid access token")
		return nil, err
	}
	rawClaims := token.Claims.(jwt.MapClaims)

	claims.Sub = rawClaims["sub"].(string)
	claims.Exp = rawClaims["exp"].(float64)
	aud := cast.ToStringSlice(rawClaims["aud"])
	claims.Aud = aud
	claims.Role = rawClaims["role"].(string)
	claims.Token = token
	return &claims, nil
}
