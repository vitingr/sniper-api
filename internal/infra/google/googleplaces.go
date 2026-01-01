package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	entity "main.go/internal/domain/entity"
	repository "main.go/internal/domain/repository"
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

type circle struct {
	Center struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"center"`
	Radius float64 `json:"radius"`
}

type locationRestriction struct {
	Circle circle `json:"circle"`
}

type searchNearbyRequest struct {
	LocationRestriction locationRestriction `json:"locationRestriction"`
}

type newPlacePhoto struct {
	Name             string `json:"name"`
	WidthPx          int    `json:"widthPx"`
	HeightPx         int    `json:"heightPx"`
	AuthorAttributions []struct {
		DisplayName string `json:"displayName"`
		Uri         string `json:"uri"`
		PhotoUri    string `json:"photoUri"`
	} `json:"authorAttributions"`
}

type newPlaceResult struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	DisplayName struct {
		Text         string `json:"text"`
		LanguageCode string `json:"languageCode"`
	} `json:"displayName"`
	FormattedAddress    string   `json:"formattedAddress"`
	BusinessStatus      string   `json:"businessStatus"`
	Types               []string `json:"types"`
	NationalPhoneNumber string   `json:"nationalPhoneNumber"`
	WebsiteUri          string   `json:"websiteUri"`

	IconMaskBaseUri     string          `json:"iconMaskBaseUri"`
	IconBackgroundColor string          `json:"iconBackgroundColor"`
	Photos              []newPlacePhoto `json:"photos"`

	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
}

type newPlacesResponse struct {
	Places []newPlaceResult `json:"places"`
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

func (r *googlePlaceRepo) buildNewPhotoURL(photoName string, maxWidth int) string {
	if photoName == "" {
		return ""
	}

	return fmt.Sprintf("https://places.googleapis.com/v1/%s/media?maxHeightPx=%d&maxWidthPx=%d&key=%s", photoName, maxWidth, maxWidth, r.apiKey)
}

func (r *googlePlaceRepo) GetNearbyWithPagination(ctx context.Context, lat, lng, radiusKm float64, pageToken string) (*NearbyResult, error) {
	endpoint := "https://places.googleapis.com/v1/places:searchNearby"

	reqBody := searchNearbyRequest{
		LocationRestriction: locationRestriction{
			Circle: circle{
				Center: struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				}{
					Latitude:  lat,
					Longitude: lng,
				},
				Radius: radiusKm * 1000,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", r.apiKey)

	fieldMask := "places.id,places.displayName,places.formattedAddress,places.location,places.businessStatus,places.types,places.nationalPhoneNumber,places.websiteUri,places.iconMaskBaseUri,places.iconBackgroundColor,places.photos"
	req.Header.Set("X-Goog-FieldMask", fieldMask)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google places v1 error: status %d", resp.StatusCode)
	}

	var googleResp newPlacesResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, err
	}

	var companies []*entity.Company
	for _, item := range googleResp.Places {

		var photoURL string
		if len(item.Photos) > 0 {
			photoURL = r.buildNewPhotoURL(item.Photos[0].Name, 400)
		}

		companies = append(companies, &entity.Company{
			ID:                  item.ID,
			Name:                item.DisplayName.Text,
			Address:             item.FormattedAddress,
			Latitude:            item.Location.Latitude,
			Longitude:           item.Location.Longitude,
			BusinessStatus:      item.BusinessStatus,
			Types:               item.Types,
			Website:             item.WebsiteUri,
			Phone:               item.NationalPhoneNumber,

			PhotoURL:            photoURL,
			IconMaskBaseURI:     item.IconMaskBaseUri,
			IconBackgroundColor: item.IconBackgroundColor,
		})
	}

	return &NearbyResult{
		Companies:     companies,
		NextPageToken: "",
	}, nil
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
