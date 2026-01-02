package repository

import (
	"context"

	"main.go/internal/domain/entity"
)

type CompanyRepository interface {
	GetNearby(ctx context.Context, lat, lng, radiusKm float64) ([]*entity.Company, error)
	SearchWithFilters(ctx context.Context, keyword, category string, lat, lng, radiusKm float64) ([]*entity.Company, error)
	Search(ctx context.Context, query string) ([]*entity.Company, error)
	GetByID(ctx context.Context, id string) (*entity.Company, error)
}
