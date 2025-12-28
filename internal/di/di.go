package di

import (
	"os"

	"main.go/internal/adapter/http/handler"
	"main.go/internal/infra/google"
	company "main.go/internal/usecase/companies"
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
