// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"bhs/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Auth -.
	Auth interface {
		Register(context.Context, entity.User) (string, error)
		Authenticate(context.Context, string, string) (string, error)
		HashPassword(string) (string, error)
		Logout(context.Context, string)
		Authorize(context.Context, string) (entity.User, error)
	}
	// UserRepo -.
	UserRepo interface {
		SaveUser(context.Context, entity.User) (*entity.User, error)
		GetUserById(context.Context, int64) (*entity.User, error)
		GetUser(context.Context, string) (*entity.User, error)
	}

	AssetRepo interface {
		GetUserAssetsPage(ctx context.Context, user entity.User, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error)
		GetAssetsPage(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error)
		AddUserAsset(context.Context, entity.User, int64) (bool, error)
	}

	Assets interface {
		GetAssetsPage(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error)
		GetUserAssets(context.Context, entity.User) ([]entity.Asset, error)
		GetUserAssetsPage(ctx context.Context, user entity.User, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error)
		AddUserAsset(context.Context, entity.User, int64) (bool, error)
	}
)
