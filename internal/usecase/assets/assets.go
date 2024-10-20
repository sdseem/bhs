package assets

import (
	"bhs/internal/entity"
	"bhs/internal/usecase"
	"context"
)

type Assets struct {
	repo usecase.AssetRepo
}

// NewAssets -.
func NewAssets(r usecase.AssetRepo) *Assets {
	return &Assets{r}
}

// GetAssetsPage -.
func (a *Assets) GetAssetsPage(ctx context.Context, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error) {
	return a.repo.GetAssetsPage(ctx, pageNumber, itemsPerPage)
}

// GetUserAssets -.
func (a *Assets) GetUserAssets(ctx context.Context, user entity.User) ([]entity.Asset, error) {
	pageNum, itemsPerPage := uint64(1), uint64(10)
	return a.repo.GetUserAssetsPage(ctx, user, pageNum, itemsPerPage)
}

// GetUserAssetsPage -.
func (a *Assets) GetUserAssetsPage(ctx context.Context, user entity.User, pageNumber uint64, itemsPerPage uint64) ([]entity.Asset, error) {
	return a.repo.GetUserAssetsPage(ctx, user, pageNumber, itemsPerPage)
}

// AddUserAsset -.
func (a *Assets) AddUserAsset(ctx context.Context, user entity.User, assetId int64) (bool, error) {
	// billing logic ???
	return a.repo.AddUserAsset(ctx, user, assetId)
}
