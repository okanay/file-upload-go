package asset

import (
	"database/sql"
	"github.com/okanay/file-upload-go/types"
)

type Repository struct {
	db *sql.DB
}

type IRepository interface {
	GetAllAssets() ([]types.Assets, error)
	GetAssetWithFilename(filename string) (types.Assets, error)
	DeleteAsset(filename string) error
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllAssets() ([]types.Assets, error) {
	var assets []types.Assets

	query := `SELECT id, creator, name, type, filename, description, size, created_at, updated_at FROM assets`

	rows, err := r.db.Query(query)
	if err != nil {
		return assets, err
	}
	defer rows.Close()

	for rows.Next() {
		var asset types.Assets
		if err := rows.Scan(&asset.ID, &asset.Creator, &asset.Name, &asset.Type, &asset.Filename, &asset.Description, &asset.Size, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return assets, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *Repository) DeleteAsset(filename string) error {
	query := `DELETE FROM assets WHERE filename = $1`

	_, err := r.db.Exec(query, filename)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAssetWithFilename(filename string) (types.Assets, error) {
	var asset types.Assets

	query := `SELECT id, creator, name, type, filename, description, size, created_at, updated_at FROM assets WHERE filename = $1`

	err := r.db.QueryRow(query, filename).Scan(&asset.ID, &asset.Creator, &asset.Name, &asset.Type, &asset.Filename, &asset.Description, &asset.Size, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		return asset, err
	}

	return asset, nil
}
