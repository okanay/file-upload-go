package upload

import (
	"database/sql"
	"github.com/okanay/file-upload-go/types"
)

type Repository struct {
	db *sql.DB
}

type IRepository interface {
	CreateAssetRecord(req types.CreateAssetReq) (types.Assets, error)
	GetAllAssets() ([]types.Assets, error)
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAssetRecord(req types.CreateAssetReq) (types.Assets, error) {
	var asset types.Assets

	// SQL sorgusunu hazırla
	query := `INSERT INTO assets (creator, name, type, filename, description, size) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, creator, name, type, filename, description, size, created_at, updated_at`

	// SQL sorgusunu çalıştır
	err := r.db.QueryRow(query, req.Creator, req.Name, req.Type, req.Filename, req.Description, req.Size).Scan(&asset.ID, &asset.Creator, &asset.Name, &asset.Type, &asset.Filename, &asset.Description, &asset.Size, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		return asset, err
	}

	return asset, nil
}

func (r *Repository) GetAllAssets() ([]types.Assets, error) {
	var assets []types.Assets

	// SQL sorgusunu hazırla
	query := `SELECT id, creator, name, type, filename, description, size, created_at, updated_at FROM assets`

	// SQL sorgusunu çalıştır
	rows, err := r.db.Query(query)
	if err != nil {
		return assets, err
	}
	defer rows.Close()

	// SQL sorgusundan dönen verileri diziye çevir
	for rows.Next() {
		var asset types.Assets
		if err := rows.Scan(&asset.ID, &asset.Creator, &asset.Name, &asset.Type, &asset.Filename, &asset.Description, &asset.Size, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return assets, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}
