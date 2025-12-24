package company

import (
	"context"
	entity "main.go/internal/domain/entities"
	repository "main.go/internal/domain/repositories"
)

type GetNearbyInput struct {
	Latitude  float64
	Longitude float64
	RadiusKm  float64
}

type GetNearbyCompaniesUseCase struct {
	repo repository.CompanyRepository
}

func NewGetNearbyCompaniesUseCase(repo repository.CompanyRepository) *GetNearbyCompaniesUseCase {
	return &GetNearbyCompaniesUseCase{repo: repo}
}

func (uc *GetNearbyCompaniesUseCase) Execute(ctx context.Context, input GetNearbyInput) ([]*entity.Company, error) {
	if input.RadiusKm <= 0 {
		input.RadiusKm = 10
	}
	return uc.repo.GetNearby(ctx, input.Latitude, input.Longitude, input.RadiusKm)
}
