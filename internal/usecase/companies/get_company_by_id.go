package company

import (
	"context"
	entity "main.go/internal/domain/entity"
	repository "main.go/internal/domain/repository"
)

type GetCompanyByIDUseCase struct {
	repository repository.CompanyRepository
}

func NewGetCompanyByIDUseCase(repository repository.CompanyRepository) *GetCompanyByIDUseCase {
	return &GetCompanyByIDUseCase{repository: repository}
}

func (uc *GetCompanyByIDUseCase) Execute(ctx context.Context, id string) (*entity.Company, error) {
	return uc.repository.GetByID(ctx, id)
}
