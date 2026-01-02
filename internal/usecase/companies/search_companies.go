package company

import (
	"context"
	entity "main.go/internal/domain/entity"
	repository "main.go/internal/domain/repository"
)

type SearchCompaniesInput struct {
	Latitude  float64
	Longitude float64
	RadiusKm  float64
	Keyword   string
	Category  string
}

type SearchCompaniesUseCase struct {
	repository repository.CompanyRepository
}

func NewSearchCompaniesUseCase(repository repository.CompanyRepository) *SearchCompaniesUseCase {
	return &SearchCompaniesUseCase{repository: repository}
}

func (uc *SearchCompaniesUseCase) Execute(ctx context.Context, input SearchCompaniesInput) ([]*entity.Company, error) {
	return uc.repository.SearchWithFilters(ctx, input.Keyword, input.Category, input.Latitude, input.Longitude, input.RadiusKm)
}
