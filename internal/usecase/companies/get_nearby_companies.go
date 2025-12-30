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
	repository repository.CompanyRepository
}

func NewGetNearbyCompaniesUseCase(repository repository.CompanyRepository) *GetNearbyCompaniesUseCase {
	return &GetNearbyCompaniesUseCase{repository: repository}
}

func (uc *GetNearbyCompaniesUseCase) Execute(ctx context.Context, input GetNearbyInput) ([]*entity.Company, error) {
	if input.RadiusKm <= 0 {
		input.RadiusKm = 10
	}

	// Se o usuário quer buscar todas as páginas (até 60 resultados)
	if input.FetchAll {
		return uc.repository.GetNearbyAll(ctx, input.Latitude, input.Longitude, input.RadiusKm)
	}

	// Comportamento padrão: retorna apenas primeira página (20 resultados)
	return uc.repository.GetNearby(ctx, input.Latitude, input.Longitude, input.RadiusKm)
}
