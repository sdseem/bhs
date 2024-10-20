package repo

import (
	"bhs/internal/entity"
	"bhs/pkg/postgres"
	"context"
	"fmt"
)

const _defaultEntityCapacity = 100

type AssetsRepo struct {
	*postgres.Postgres
}

func NewAssetsRepo(pg *postgres.Postgres) *AssetsRepo {
	return &AssetsRepo{pg}
}

func (r *AssetsRepo) GetUserAssetsPage(ctx context.Context, user entity.User, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error) {
	offset := (pageNumber - 1) * itemsPerPage
	b := r.Builder.
		Select("a.id, a.name, a.description, a.price").
		From("assets as a").
		InnerJoin("assets_library as al on a.id = al.asset_id and al.buyer_id = ?", user.Id).
		OrderBy("a.name").
		Offset(offset).
		Limit(itemsPerPage)

	sql, args, err := b.ToSql()
	if err != nil {
		return nil, fmt.Errorf("AssetsRepo error: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AssetsRepo - GetUserAssets - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.Asset, 0, itemsPerPage)
	for rows.Next() {
		e := entity.Asset{}
		err = rows.Scan(&e.Id, &e.Name, &e.Description, &e.Price)
		if err != nil {
			return nil, fmt.Errorf("AssetsRepo - GetUserAssets - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}
	return entities, nil
}

func (r *AssetsRepo) GetAssetsPage(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error) {
	offset := (pageNumber - 1) * itemsPerPage
	b := r.Builder.
		Select("a.id, a.name, a.description, a.price").
		From("assets as a").
		OrderBy("a.name").
		Limit(itemsPerPage).
		Offset(offset)
	sql, args, err := b.ToSql()

	if err != nil {
		return nil, fmt.Errorf("AssetsRepo error: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AssetsRepo - GetAssetsPage - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	entities := make([]entity.Asset, 0, itemsPerPage)
	for rows.Next() {
		e := entity.Asset{}

		err = rows.Scan(&e.Id, &e.Name, &e.Description, &e.Price)
		if err != nil {
			return nil, fmt.Errorf("AssetsRepo - GetUserAssets - rows.Scan: %w", err)
		}

		entities = append(entities, e)
	}

	return entities, nil
}

func (r *AssetsRepo) AddUserAsset(ctx context.Context, user entity.User, assetId int64) (bool, error) {
	sql, args, err := r.Builder.
		Insert("assets_library").
		Columns("buyer_id, asset_id").
		Values(user.Id, assetId).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("save user error: %w", err)
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("save user error: %w", err)
	}
	return true, nil
}
