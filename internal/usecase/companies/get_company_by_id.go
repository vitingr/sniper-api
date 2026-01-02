package company

import (
	"context"

	entity "main.go/internal/domain/entity"
	repository "main.go/internal/domain/repository"
	"main.go/internal/infra/hunter"
)

type GetCompanyByIDUseCase struct {
	repository   repository.CompanyRepository
	hunterClient *hunter.HunterClient
}

func NewGetCompanyByIDUseCase(repository repository.CompanyRepository, hunterClient *hunter.HunterClient) *GetCompanyByIDUseCase {
	return &GetCompanyByIDUseCase{repository: repository, hunterClient: hunterClient}
}

func (uc *GetCompanyByIDUseCase) Execute(ctx context.Context, id string) (*entity.Company, error) {
	company, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if company.Website != "" && uc.hunterClient != nil {

		domain := uc.hunterClient.ExtractDomain(company.Website)

		if domain != "" {
			// Busca contatos executivos/gerentes (limite 5)
			hunterEmails, err := uc.hunterClient.GetRelevantContacts(domain, 5)

			if err == nil && len(hunterEmails) > 0 {
				for _, he := range hunterEmails {
					company.Contacts = append(company.Contacts, entity.ContactEmail{
						Value:      he.Value,
						Type:       he.Type,
						Position:   he.Position,
						Seniority:  he.Seniority,
						Department: he.Department,
						Confidence: he.Confidence,
					})
				}
			}
			// Se der erro no Hunter, apenas ignoramos e retornamos os dados do Google
			// (Não queremos falhar o request todo se o Hunter cair ou acabar cota)
		}
	}

	return company, nil
}
