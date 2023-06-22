package v1

import (
	"github.com/gin-gonic/gin"
	t "github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/api/tokens"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/config"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/pkg/logger"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/storage"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/storage/redisrepo"
)

type HandlerV1I interface {
	UserCreate(c *gin.Context)
	UserGet(c *gin.Context)
	UserFind(c *gin.Context)
	UserUpdate(c *gin.Context)
	UserDelete(c *gin.Context)
	// Don't delete this line, it is used to modify the file automatically
}

type handlerV1 struct {
	log        *logger.Logger
	cfg        config.Config
	storage    storage.StorageI
	jwthandler t.JWTHandler
	redis      redisrepo.InMemoryStorageI
}

type HandlerV1Config struct {
	Logger     *logger.Logger
	Cfg        config.Config
	Postgres   storage.StorageI
	JWTHandler t.JWTHandler
	Redis      redisrepo.InMemoryStorageI
}

// New ...
func New(c *HandlerV1Config) HandlerV1I {
	return &handlerV1{
		log:        c.Logger,
		cfg:        c.Cfg,
		storage:    c.Postgres,
		jwthandler: c.JWTHandler,
		redis:      c.Redis,
	}
}
