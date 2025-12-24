package di

import (
	"os"

	"main.go/internal/infra/google"
	handler "main.go/internal/infra/http/handlers"
	company "main.go/internal/use-cases/companies"
)

type Container struct {
	CompanyHandler *handler.CompanyHandler
}

func NewContainer() *Container {
	apiKey := os.Getenv("GOOGLE_API_KEY")
    if apiKey == "" {
        panic("GOOGLE_API_KEY is missing")
    }

	repo := google.NewGooglePlaceRepository(apiKey)
	nearbyUC := company.NewGetNearbyCompaniesUseCase(repo)
	getByIDUC := company.NewGetCompanyByIDUseCase(repo)

	companyHandler := handler.NewCompanyHandler(getByIDUC, nearbyUC)

	return &Container{
		CompanyHandler: companyHandler,
	}
}
