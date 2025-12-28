package company

import (
	"context"
	entity "main.go/internal/domain/entity"
	repository "main.go/internal/domain/repository"
)

type GetCompanyByIDUseCase struct {
	repo repository.CompanyRepository
}

func NewGetCompanyByIDUseCase(repo repository.CompanyRepository) *GetCompanyByIDUseCase {
	return &GetCompanyByIDUseCase{repo: repo}
}

func (uc *GetCompanyByIDUseCase) Execute(ctx context.Context, id string) (*entity.Company, error) {
	return uc.repo.GetByID(ctx, id)
}
