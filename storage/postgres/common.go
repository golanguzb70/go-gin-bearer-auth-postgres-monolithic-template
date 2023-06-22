package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
)

func (p *postgresRepo) CheckIfExists(ctx context.Context, req *models.CheckIfExistsReq) (*models.CheckIfExistsRes, error) {
	var (
		res sql.NullBool
	)
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s='%s') AS value_exists", req.Table, req.Column, req.Value)
	err := p.Db.Db.QueryRow(query).Scan(&res)
	err = HandleDatabaseError(err, p.Log, "CheckIfExists")

	return &models.CheckIfExistsRes{
		Exists: res.Bool,
	}, err
}
