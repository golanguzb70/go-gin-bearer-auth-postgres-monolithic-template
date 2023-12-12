package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
)

func (r *postgresRepo) TemplateCreate(ctx context.Context, req *models.TemplateCreateReq) (*models.TemplateResponse, error) {
	res := &models.TemplateResponse{}
	query := r.Db.Builder.Insert("templates").Columns(
		"template_name",
	).Values(req.TemplateName).Suffix(
		"RETURNING id, template_name, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.Id, &res.TemplateName,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "TemplateCreate: query.RunWith(r.Db.Db).Scan()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) TemplateGet(ctx context.Context, req *models.TemplateGetReq) (*models.TemplateResponse, error) {
	query := r.Db.Builder.Select("id, template_name, created_at, updated_at").
		From("templates")

	if req.Id != "" {
		query = query.Where(squirrel.Eq{"id": req.Id})
	} else {
		return &models.TemplateResponse{}, fmt.Errorf("at least one filter should be exists")
	}

	res := &models.TemplateResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.TemplateName,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "TemplateGet:query.RunWith(r.Db.Db).QueryRow()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) TemplateFind(ctx context.Context, req *models.TemplateFindReq) (*models.TemplateFindResponse, error) {
	var (
		res            = &models.TemplateFindResponse{}
		whereCondition = squirrel.And{}
		orderBy        = []string{}
	)

	if strings.TrimSpace(req.Search) != "" {
		whereCondition = append(whereCondition, squirrel.ILike{"template_name": req.Search + "%"})
	}

	if req.OrderByCreatedAt != 0 {
		if req.OrderByCreatedAt > 0 {
			orderBy = append(orderBy, "created_at DESC")
		} else {
			orderBy = append(orderBy, "created_at ASC")
		}
	}

	countQuery := r.Db.Builder.Select("count(1) as count").From("templates").Where("deleted_at is null").Where(whereCondition)
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "TemplateFind: countQuery.RunWith(r.Db.Db).QueryRow().Scan()")
	}

	query := r.Db.Builder.Select("id, template_name, created_at, updated_at").
		From("templates").Where("deleted_at is null").Where(whereCondition)

	if len(orderBy) > 0 {
		query = query.OrderBy(strings.Join(orderBy, ", "))
	}

	query = query.Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "TemplateFind: query.RunWith(r.Db.Db).Query()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.TemplateResponse{}
		err := rows.Scan(
			&temp.Id, &temp.TemplateName,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "TemplateFind: rows.Scan()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)
		res.Templates = append(res.Templates, temp)
	}

	return res, nil
}

func (r *postgresRepo) TemplateUpdate(ctx context.Context, req *models.TemplateUpdateReq) (*models.TemplateResponse, error) {
	var (
		mp             = make(map[string]interface{})
		whereCondition = squirrel.And{squirrel.Eq{"id": req.Id}}
	)
	mp["template_name"] = req.TemplateName
	mp["updated_at"] = time.Now()

	query := r.Db.Builder.Update("templates").SetMap(mp).
		Where(whereCondition).
		Suffix("RETURNING id, template_name, created_at, updated_at")

	res := &models.TemplateResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.TemplateName,
		&CreatedAt, &UpdatedAt,
	)

	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "TemplateUpdate: query.RunWith(r.Db.Db).QueryRow().Scan()")
	}
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) TemplateDelete(ctx context.Context, req *models.TemplateDeleteReq) error {
	whereCondition := squirrel.And{squirrel.Eq{"id": req.Id}}

	query := r.Db.Builder.Delete("templates").Where(whereCondition)

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "TemplateDelete: query.RunWith(r.Db.Db).Exec()")
}
