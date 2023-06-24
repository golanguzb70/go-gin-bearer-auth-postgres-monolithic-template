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

func (p *postgresRepo) UpdateSingleField(ctx context.Context, req *models.UpdateSingleFieldReq) error {
	query := fmt.Sprintf("UPDATE %s SET %s=$1 where id=$2", req.Table, req.Column)

	_, err := p.Db.Db.Exec(query, req.NewValue, req.Id)
	return HandleDatabaseError(err, p.Log, "UpdateSingleField")
}
