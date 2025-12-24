package company

import (
	"context"
	entity "main.go/internal/domain/entities"
	repository "main.go/internal/domain/repositories"
)

type ListCompaniesInput struct {
	Page  int
	Limit int
}

type ListCompaniesUseCase struct {
	repo repository.CompanyRepository
}

func NewListCompaniesUseCase(repo repository.CompanyRepository) *ListCompaniesUseCase {
	return &ListCompaniesUseCase{repo: repo}
}

func (uc *ListCompaniesUseCase) Execute(ctx context.Context, input ListCompaniesInput) ([]*entity.Company, error) {
	// if input.Limit <= 0 {
	// 	input.Limit = 10
	// }
	// if input.Page <= 0 {
	// 	input.Page = 1
	// }
	// offset := (input.Page - 1) * input.Limit

	return uc.repo.Search(ctx, "")
}
