package usecase

import (
	"context"
	"time"

	"main.go/internal/domain/entity"
)

// External services interfaces (Dependency Inversion)
// Implentations are into infraestructure layer
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type ApolloService interface {
	SearchContacts(ctx context.Context, companyID string) ([]entity.Contact, error)
	EnrichContact(ctx context.Context, email string) (*entity.Contact, error)
}

type GeocodingService interface {
	GeocodeAddress(ctx context.Context, address string) (lat, lng float64, err error)
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}
