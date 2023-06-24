package postgres

import (
	"context"

	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
)

type PostgresI interface {
	// common
	UpdateSingleField(ctx context.Context, req *models.UpdateSingleFieldReq) error
	CheckIfExists(ctx context.Context, req *models.CheckIfExistsReq) (*models.CheckIfExistsRes, error)

	// User
	UserCreate(ctx context.Context, req *models.UserCreateReq) (*models.UserResponse, error)
	UserGet(ctx context.Context, req *models.UserGetReq) (*models.UserResponse, error)
	UserFind(ctx context.Context, req *models.UserFindReq) (*models.UserFindResponse, error)
	UserUpdate(ctx context.Context, req *models.UserUpdateReq) (*models.UserResponse, error)
	UserDelete(ctx context.Context, req *models.UserDeleteReq) error

	// Template 
	TemplateCreate(ctx context.Context, req *models.TemplateCreateReq) (*models.TemplateResponse, error)
	TemplateGet(ctx context.Context, req *models.TemplateGetReq) (*models.TemplateResponse, error)
	TemplateFind(ctx context.Context, req *models.TemplateFindReq) (*models.TemplateFindResponse, error)
	TemplateUpdate(ctx context.Context, req *models.TemplateUpdateReq) (*models.TemplateResponse, error)
	TemplateDelete(ctx context.Context, req *models.TemplateDeleteReq) error
	
	// Don't delete this line, it is used to modify the file automatically
}
