package repository

import (
	"context"

	entity "main.go/internal/domain/entities"
)

type CompanyRepository interface {
	GetNearby(ctx context.Context, lat, lng, radiusKm float64) ([]*entity.Company, error)
	GetNearbyAll(ctx context.Context, lat, lng, radiusKm float64) ([]*entity.Company, error)
	Search(ctx context.Context, query string) ([]*entity.Company, error)
	GetByID(ctx context.Context, id string) (*entity.Company, error)
}
