package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/docs" // docs
	v1 "github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/handlers/v1"
	t "github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/tokens"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/config"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/pkg/logger"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/storage"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/storage/redis"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option ...
type Option struct {
	Conf       config.Config
	Logger     *logger.Logger
	Postgres   storage.StorageI
	JWTHandler t.JWTHandler
	Redis      redis.InMemoryStorageI
}

// New ...
// @title           User project API Endpoints
// @version         1.0
// @description     Here QA can test and frontend or mobile developers can get information of API endpoints.

// @BasePath  /v1

// @securityDefinitions.basic BearerAuth
// @in header
// @name Authorization
func New(log *logger.Logger, cfg config.Config, strg storage.StorageI) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	h := v1.New(&v1.HandlerV1Config{
		Logger:   log,
		Cfg:      cfg,
		Postgres: strg,
	})

	api := router.Group("/v1")

	user := api.Group("/user")
	user.POST("", h.UserCreate)
	user.GET("/:id", h.UserGet)
	user.GET("/list", h.UserFind)
	user.PUT("", h.UserUpdate)
	user.DELETE(":id", h.UserDelete)

	// Don't delete this line, it is used to modify the file automatically

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
