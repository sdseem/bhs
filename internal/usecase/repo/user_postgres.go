package repo

import (
	"bhs/internal/entity"
	"bhs/pkg/postgres"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"strconv"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) GetUserById(ctx context.Context, userId int64) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, username, password_hash").
		From("users").
		Where(squirrel.Eq{"id": strconv.FormatInt(userId, 15)}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo error: %w", err)
	}
	userRow := r.Pool.QueryRow(ctx, sql, args...)
	user := entity.User{}
	err = userRow.Scan(&user.Id, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepo error: %w", err)
	}
	return &user, nil
}

func (r *UserRepo) GetUser(ctx context.Context, username string) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, username, password_hash").
		From("users").
		Where("username = ?", username).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo error: %w", err)
	}
	userRow := r.Pool.QueryRow(ctx, sql, args...)
	savedUser := entity.User{}
	err = userRow.Scan(&savedUser.Id, &savedUser.Username, &savedUser.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepo error: %w", err)
	}
	return &savedUser, nil
}

func (r *UserRepo) SaveUser(ctx context.Context, user entity.User) (*entity.User, error) {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("username, password_hash").
		Values(user.Username, user.PasswordHash).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("save user error: %w", err)
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("save user error: %w", err)
	}

	sql, args, err = r.Builder.
		Select("id, username, password_hash").
		From("users").
		Where("username = ?", user.Username).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo query builder error: %w", err)
	}
	userRow := r.Pool.QueryRow(ctx, sql, args...)
	savedUser := entity.User{}
	err = userRow.Scan(&savedUser.Id, &savedUser.Username, &savedUser.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepo query result error: %w", err)
	}
	return &savedUser, nil
}
