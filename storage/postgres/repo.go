package postgres

import (
	"context"

	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
)

type PostgresI interface {
	// common
	UpdateSingleField(ctx context.Context, req *models.UpdateSingleFieldReq) error

	CheckIfExists(ctx context.Context, req *models.CheckIfExistsReq) (*models.CheckIfExistsRes, error)
	UserCreate(ctx context.Context, req *models.UserCreateReq) (*models.UserResponse, error)
	UserGet(ctx context.Context, req *models.UserGetReq) (*models.UserResponse, error)
	UserFind(ctx context.Context, req *models.UserFindReq) (*models.UserFindResponse, error)
	UserUpdate(ctx context.Context, req *models.UserUpdateReq) (*models.UserResponse, error)
	UserDelete(ctx context.Context, req *models.UserDeleteReq) error
	// Don't delete this line, it is used to modify the file automatically
}
