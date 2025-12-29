package di

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"main.go/internal/adapter/http/handler"
	"main.go/internal/adapter/repository"
	"main.go/internal/infra/database/postgres"
	"main.go/internal/infra/google"
	company "main.go/internal/usecase/companies"
	user "main.go/internal/usecase/users"
	"main.go/pkg/logger"
)

type Container struct {
	CompanyHandler *handler.CompanyHandler
	UserHandler    *handler.UserHandler
	DBPool         *pgxpool.Pool
}

func NewContainer() *Container {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		logger.Fatal("GOOGLE_API_KEY is missing")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Fatal("DATABASE_URL is missing")
	}

	dbPool, err := postgres.NewConnection(databaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database")
	}

	company_repository := google.NewGooglePlaceRepository(apiKey)
	company_nearbyUC := company.NewGetNearbyCompaniesUseCase(company_repository)
	company_getByIDUC := company.NewGetCompanyByIDUseCase(company_repository)
	companyHandler := handler.NewCompanyHandler(company_getByIDUC, company_nearbyUC)

	user_repository := repository.NewUserRepository(dbPool)
	user_createUserUC := user.NewCreateUserUseCase(*user_repository)
	user_getUserByEmailUC := user.NewGetUserByEmailUseCase(*user_repository)
	user_listAllUsersUC := user.NewListAllUsersUseCase(*user_repository)
	user_updateUserUC := user.NewUpdateUserUseCase(*user_repository)
	user_deleteUserUC := user.NewDeleteUserUseCase(*user_repository)
	userHandler := handler.NewUserHandler(
		user_createUserUC,
		user_getUserByEmailUC,
		user_listAllUsersUC,
		user_updateUserUC,
		user_deleteUserUC,
	)

	return &Container{
		CompanyHandler: companyHandler,
		UserHandler:    userHandler,
		DBPool:         dbPool,
	}
}

func (c *Container) Close() {
	if c.DBPool != nil {
		c.DBPool.Close()
	}
}
