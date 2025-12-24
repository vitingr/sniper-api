package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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
	PlaceID          string `json:"place_id"`
	Name             string `json:"name"`
	Vicinity         string `json:"vicinity"`
	FormattedAddress string `json:"formatted_address"`
	FormattedPhone   string `json:"formatted_phone_number"`
	Website          string `json:"website"`
	Geometry         struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

type googleListResponse struct {
	Results []googlePlaceResult `json:"results"`
	Status  string              `json:"status"`
}

type googleDetailResponse struct {
	Result googlePlaceResult `json:"result"`
	Status string            `json:"status"`
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
			ID:        item.PlaceID,
			Name:      item.Name,
			Address:   addr,
			Latitude:  item.Geometry.Location.Lat,
			Longitude: item.Geometry.Location.Lng,
		})
	}

	return companies, nil
}

func (r *googlePlaceRepo) GetNearby(ctx context.Context, lat, lng, radiusKm float64) ([]*entity.Company, error) {
	endpoint := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("location", fmt.Sprintf("%f,%f", lat, lng))
	q.Set("radius", fmt.Sprintf("%f", radiusKm*1000))
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
		addr := item.Vicinity
		if addr == "" {
			addr = item.FormattedAddress
		}

		companies = append(companies, &entity.Company{
			ID:        item.PlaceID,
			Name:      item.Name,
			Address:   addr,
			Latitude:  item.Geometry.Location.Lat,
			Longitude: item.Geometry.Location.Lng,
		})
	}

	return companies, nil
}

func (r *googlePlaceRepo) GetByID(ctx context.Context, id string) (*entity.Company, error) {
	endpoint := "https://maps.googleapis.com/maps/api/place/details/json"
	u, _ := url.Parse(endpoint)
	q := u.Query()
	q.Set("place_id", id)
	q.Set("fields", "place_id,name,formatted_address,geometry,formatted_phone_number,website,vicinity")
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
		ID:        item.PlaceID,
		Name:      item.Name,
		Address:   item.FormattedAddress,
		Latitude:  item.Geometry.Location.Lat,
		Longitude: item.Geometry.Location.Lng,
		// Phone: item.FormattedPhone,
		// Website: item.Website,
	}, nil
}

func (r *googlePlaceRepo) List(ctx context.Context, limit, offset int) ([]*entity.Company, error) {
	return []*entity.Company{}, nil
}

func (r *googlePlaceRepo) Count(ctx context.Context) (int, error) {
	return 0, nil
}
