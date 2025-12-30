package company

import (
	"context"

	entity "main.go/internal/domain/entity"
	repository "main.go/internal/domain/repository"
)

type ListCompaniesInput struct {
	Page  int
	Limit int
}

type ListCompaniesUseCase struct {
	repository repository.CompanyRepository
}

func NewListCompaniesUseCase(repository repository.CompanyRepository) *ListCompaniesUseCase {
	return &ListCompaniesUseCase{repository: repository}
}

func (uc *ListCompaniesUseCase) Execute(ctx context.Context, input ListCompaniesInput) ([]*entity.Company, error) {
	// if input.Limit <= 0 {
	// 	input.Limit = 10
	// }
	// if input.Page <= 0 {
	// 	input.Page = 1
	// }
	// offset := (input.Page - 1) * input.Limit

	return uc.repository.Search(ctx, "")
}
