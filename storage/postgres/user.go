package postgres

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golanguzb70/go-gin-bearer-auth-postgres-monolithic-template/models"
)

func (r *postgresRepo) UserCreate(ctx context.Context, req *models.UserCreateReq) (*models.UserResponse, error) {
	res := &models.UserResponse{}
	query := r.Db.Builder.Insert("users").Columns(
		"id, user_name, email, hashed_password, refresh_token",
	).Values(req.Id, req.UserName, req.Email, req.Password, req.RefreshToken).Suffix(
		"RETURNING id, user_name, email, hashed_password, refresh_token, created_at, updated_at")

	err := query.RunWith(r.Db.Db).Scan(
		&res.Id, &res.UserName,
		&res.Email, &res.Password,
		&res.RefreshToken, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "(r *UserRepo) Create()")
	}
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) UserGet(ctx context.Context, req *models.UserGetReq) (*models.UserResponse, error) {
	query := r.Db.Builder.Select("id, user_name, created_at, updated_at").
		From("users").
		Where(squirrel.Eq{"id": req.Id})

	res := &models.UserResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.UserName,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "(r *UserRepo) Get()")
	}

	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) UserFind(ctx context.Context, req *models.UserFindReq) (*models.UserFindResponse, error) {
	var (
		res = &models.UserFindResponse{}
	)

	countQuery := r.Db.Builder.Select("count(1) as count").From("users").Where("deleted_at is null")
	err := countQuery.RunWith(r.Db.Db).QueryRow().Scan(&res.Count)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "(r *models.UserUserRepo) FindList()")

	}

	query := r.Db.Builder.Select("id, user_name, created_at, updated_at").
		From("users").Where("deleted_at is null").OrderBy("id").Limit(uint64(req.Limit)).Offset(uint64((req.Page - 1) * req.Limit))

	rows, err := query.RunWith(r.Db.Db).Query()
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "(r *models.UserUserRepo) FindList()")
	}
	defer rows.Close()

	for rows.Next() {
		temp := &models.UserResponse{}
		err := rows.Scan(
			&temp.Id, &temp.UserName,
			&CreatedAt, &UpdatedAt,
		)
		if err != nil {
			return res, HandleDatabaseError(err, r.Log, "(r *models.UserUserRepo) FindList()")
		}

		temp.CreatedAt = CreatedAt.Format(time.RFC1123)
		temp.UpdatedAt = UpdatedAt.Format(time.RFC1123)
		res.Users = append(res.Users, temp)
	}

	return res, nil
}

func (r *postgresRepo) UserUpdate(ctx context.Context, req *models.UserUpdateReq) (*models.UserResponse, error) {
	mp := make(map[string]interface{})
	mp["user_name"] = req.UserName
	mp["updated_at"] = time.Now()
	query := r.Db.Builder.Update("users").SetMap(mp).
		Where(squirrel.Eq{"id": req.Id}).
		Suffix("RETURNING id, user_name, created_at, updated_at")

	res := &models.UserResponse{}
	err := query.RunWith(r.Db.Db).QueryRow().Scan(
		&res.Id, &res.UserName,
		&CreatedAt, &UpdatedAt,
	)
	if err != nil {
		return res, HandleDatabaseError(err, r.Log, "(r *models.UserUserRepo) Update()")
	}
	res.CreatedAt = CreatedAt.Format(time.RFC1123)
	res.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return res, nil
}

func (r *postgresRepo) UserDelete(ctx context.Context, req *models.UserDeleteReq) error {
	query := r.Db.Builder.Delete("users").Where(squirrel.Eq{"id": req.Id})

	_, err := query.RunWith(r.Db.Db).Exec()
	return HandleDatabaseError(err, r.Log, "(r *models.UserUserRepo) Delete()")
}
