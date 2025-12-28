package company

import (
	"context"
	entity "main.go/internal/domain/entity"
	repository "main.go/internal/domain/repository"
)

type GetNearbyInput struct {
	Latitude  float64
	Longitude float64
	RadiusKm  float64
	FetchAll  bool
	PageToken string
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

	// Se o usuário quer buscar todas as páginas (até 60 resultados)
	if input.FetchAll {
		return uc.repo.GetNearbyAll(ctx, input.Latitude, input.Longitude, input.RadiusKm)
	}

	// Comportamento padrão: retorna apenas primeira página (20 resultados)
	return uc.repo.GetNearby(ctx, input.Latitude, input.Longitude, input.RadiusKm)
}
