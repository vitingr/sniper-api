package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	entity "main.go/internal/domain/entities"
	repository "main.go/internal/domain/repositories"
)

type googlePlaceRepo struct {
	apiKey string
	client *http.Client
}

var _ repository.CompanyRepository = (*googlePlaceRepo)(nil)

func NewGooglePlaceRepository(apiKey string) repository.CompanyRepository {
	return &googlePlaceRepo{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type googlePlaceResult struct {
	PlaceID          string   `json:"place_id"`
	Name             string   `json:"name"`
	Vicinity         string   `json:"vicinity"`
	FormattedAddress string   `json:"formatted_address"`
	BusinessStatus   string   `json:"business_status"`
	Types            []string `json:"types"`
	FormattedPhone   string   `json:"formatted_phone_number"`
	Website          string   `json:"website"`
	EditorialSummary struct {
		Overview string `json:"overview"`
	} `json:"editorial_summary"`
	Geometry struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

type googleListResponse struct {
	Results       []googlePlaceResult `json:"results"`
	Status        string              `json:"status"`
	NextPageToken string              `json:"next_page_token"`
}

type googleDetailResponse struct {
	Result googlePlaceResult `json:"result"`
	Status string            `json:"status"`
}

type NearbyResult struct {
	Companies     []*entity.Company
	NextPageToken string
}

func (r *googlePlaceRepo) Search(ctx context.Context, query string) ([]*entity.Company, error) {
	endpoint := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("query", query)
	q.Set("key", r.apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleResp googleListResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, err
	}

	var companies []*entity.Company
	for _, item := range googleResp.Results {
		addr := item.FormattedAddress
		if addr == "" {
			addr = item.Vicinity
		}

		companies = append(companies, &entity.Company{
			ID:             item.PlaceID,
			Name:           item.Name,
			Address:        addr,
			Latitude:       item.Geometry.Location.Lat,
			Longitude:      item.Geometry.Location.Lng,
			BusinessStatus: item.BusinessStatus,
			Types:          item.Types,
		})
	}

	return companies, nil
}

func (r *googlePlaceRepo) GetNearby(ctx context.Context, lat, lng, radiusKm float64) ([]*entity.Company, error) {
	result, err := r.GetNearbyWithPagination(ctx, lat, lng, radiusKm, "")
	if err != nil {
		return nil, err
	}
	return result.Companies, nil
}

func (r *googlePlaceRepo) GetNearbyWithPagination(ctx context.Context, lat, lng, radiusKm float64, pageToken string) (*NearbyResult, error) {
	endpoint := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	u, _ := url.Parse(endpoint)
	q := u.Query()

	if pageToken != "" {
		// Quando há pageToken, só precisamos dele e da key
		q.Set("pagetoken", pageToken)
	} else {
		// Primeira requisição com location e radius
		q.Set("location", fmt.Sprintf("%f,%f", lat, lng))
		q.Set("radius", fmt.Sprintf("%f", radiusKm*1000))
	}

	q.Set("key", r.apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var googleResp googleListResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, err
	}

	if googleResp.Status != "OK" && googleResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("google places api error: %s", googleResp.Status)
	}

	var companies []*entity.Company
	for _, item := range googleResp.Results {
		addr := item.Vicinity
		if addr == "" {
			addr = item.FormattedAddress
		}

		companies = append(companies, &entity.Company{
			ID:             item.PlaceID,
			Name:           item.Name,
			Address:        addr,
			Latitude:       item.Geometry.Location.Lat,
			Longitude:      item.Geometry.Location.Lng,
			BusinessStatus: item.BusinessStatus,
			Types:          item.Types,
		})
	}

	return &NearbyResult{
		Companies:     companies,
		NextPageToken: googleResp.NextPageToken,
	}, nil
}

// NOVO: Buscar todas as páginas automaticamente (até 60 resultados)
func (r *googlePlaceRepo) GetNearbyAll(ctx context.Context, lat, lng, radiusKm float64) ([]*entity.Company, error) {
	var allCompanies []*entity.Company
	var pageToken string

	// Google permite no máximo 3 páginas (60 resultados)
	maxPages := 3

	for page := 0; page < maxPages; page++ {
		result, err := r.GetNearbyWithPagination(ctx, lat, lng, radiusKm, pageToken)
		if err != nil {
			return nil, err
		}

		allCompanies = append(allCompanies, result.Companies...)

		if result.NextPageToken == "" {
			break
		}

		pageToken = result.NextPageToken

		// Google requer delay de ~2 segundos entre requisições de paginação
		// IMPORTANTE: Não fazer delay na última iteração
		if page < maxPages-1 && result.NextPageToken != "" {
			select {
			case <-ctx.Done():
				return allCompanies, ctx.Err()
			case <-time.After(2 * time.Second):
			}
		}
	}

	return allCompanies, nil
}

func (r *googlePlaceRepo) GetByID(ctx context.Context, id string) (*entity.Company, error) {
	endpoint := "https://maps.googleapis.com/maps/api/place/details/json"
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("place_id", id)
	fields := "place_id,name,formatted_address,geometry,formatted_phone_number,website,vicinity,business_status,types,editorial_summary"
	q.Set("fields", fields)
	q.Set("key", r.apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleResp googleDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, err
	}

	if googleResp.Status != "OK" {
		return nil, fmt.Errorf("place not found or api error: %s", googleResp.Status)
	}

	item := googleResp.Result

	return &entity.Company{
		ID:             item.PlaceID,
		Name:           item.Name,
		Address:        item.FormattedAddress,
		Latitude:       item.Geometry.Location.Lat,
		Longitude:      item.Geometry.Location.Lng,
		BusinessStatus: item.BusinessStatus,
		Types:          item.Types,
		Phone:          item.FormattedPhone,
		Website:        item.Website,
		Summary:        item.EditorialSummary.Overview,
	}, nil
}

func (r *googlePlaceRepo) List(ctx context.Context, limit, offset int) ([]*entity.Company, error) {
	return []*entity.Company{}, nil
}

func (r *googlePlaceRepo) Count(ctx context.Context) (int, error) {
	return 0, nil
}
