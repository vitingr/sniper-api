package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

type locationBias struct {
	Circle circle `json:"circle"`
}

type searchNearbyRequest struct {
	LocationRestriction locationRestriction `json:"locationRestriction"`
}

type searchTextRequest struct {
	TextQuery    string        `json:"textQuery"`
	LocationBias *locationBias `json:"locationBias,omitempty"`
}

type newPlacePhoto struct {
	Name               string `json:"name"`
	WidthPx            int    `json:"widthPx"`
	HeightPx           int    `json:"heightPx"`
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

type newPlaceReview struct {
	Name                           string `json:"name"`
	RelativePublishTimeDescription string `json:"relativePublishTimeDescription"`
	Rating                         int    `json:"rating"`
	Text                           struct {
		Text         string `json:"text"`
		LanguageCode string `json:"languageCode"`
	} `json:"text"`
	OriginalText struct {
		Text         string `json:"text"`
		LanguageCode string `json:"languageCode"`
	} `json:"originalText"`
	AuthorAttribution struct {
		DisplayName string `json:"displayName"`
		PhotoUri    string `json:"photoUri"`
	} `json:"authorAttribution"`
	PublishTime string `json:"publishTime"`
}

type newPlaceGenerativeSummary struct {
	Overview struct {
		Text         string `json:"text"`
		LanguageCode string `json:"languageCode"`
	} `json:"overview"`
	Description struct {
		Text         string `json:"text"`
		LanguageCode string `json:"languageCode"`
	} `json:"description"`
}

type newPlaceDetailResult struct {
	ID          string `json:"id"`
	DisplayName struct {
		Text string `json:"text"`
	} `json:"displayName"`
	FormattedAddress string `json:"formattedAddress"`
	Location         struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	BusinessStatus      string   `json:"businessStatus"`
	Types               []string `json:"types"`
	NationalPhoneNumber string   `json:"nationalPhoneNumber"`
	WebsiteUri          string   `json:"websiteUri"`
	Rating              float64  `json:"rating"`
	UserRatingCount     int      `json:"userRatingCount"`
	GoogleMapsUri       string   `json:"googleMapsUri"`
	EditorialSummary    struct {
		Text string `json:"text"`
	} `json:"editorialSummary"`
	GenerativeSummary newPlaceGenerativeSummary `json:"generativeSummary"`

	IconMaskBaseUri     string          `json:"iconMaskBaseUri"`
	IconBackgroundColor string          `json:"iconBackgroundColor"`
	Photos              []newPlacePhoto `json:"photos"`

	Reviews []newPlaceReview `json:"reviews"`
}

type newPlacesResponse struct {
	Places []newPlaceResult `json:"places"`
}

func (r *googlePlaceRepo) Search(ctx context.Context, query string) ([]*entity.Company, error) {
	endpoint := "https://places.googleapis.com/v1/places:searchText"

	reqBody := searchTextRequest{
		TextQuery: query,
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
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.formattedAddress,places.location,places.businessStatus,places.types,places.nationalPhoneNumber,places.websiteUri,places.iconMaskBaseUri,places.iconBackgroundColor,places.photos")

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
			IconMaskBaseURI:     item.IconMaskBaseUri + ".svg",
			IconBackgroundColor: item.IconBackgroundColor,
		})
	}

	return companies, nil
}

func (r *googlePlaceRepo) SearchWithFilters(ctx context.Context, keyword, category string, lat, lng, radiusKm float64) ([]*entity.Company, error) {
	endpoint := "https://places.googleapis.com/v1/places:searchText"

	finalQuery := keyword
	if category != "" {
		finalQuery = fmt.Sprintf("%s %s", finalQuery, category)
	}
	if finalQuery == "" {
		finalQuery = "point of interest"
	}

	reqBody := searchTextRequest{
		TextQuery: finalQuery,
		LocationBias: &locationBias{
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
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.formattedAddress,places.location,places.businessStatus,places.types,places.nationalPhoneNumber,places.websiteUri,places.iconMaskBaseUri,places.iconBackgroundColor,places.photos")

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
			IconMaskBaseURI:     item.IconMaskBaseUri + ".svg",
			IconBackgroundColor: item.IconBackgroundColor,
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
			ID:             item.ID,
			Name:           item.DisplayName.Text,
			Address:        item.FormattedAddress,
			Latitude:       item.Location.Latitude,
			Longitude:      item.Location.Longitude,
			BusinessStatus: item.BusinessStatus,
			Types:          item.Types,
			Website:        item.WebsiteUri,
			Phone:          item.NationalPhoneNumber,

			PhotoURL:            photoURL,
			IconMaskBaseURI:     item.IconMaskBaseUri + ".svg",
			IconBackgroundColor: item.IconBackgroundColor,
		})
	}

	return &NearbyResult{
		Companies:     companies,
		NextPageToken: "",
	}, nil
}

func (r *googlePlaceRepo) GetByID(ctx context.Context, id string) (*entity.Company, error) {
	endpoint := fmt.Sprintf("https://places.googleapis.com/v1/places/%s", id)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", r.apiKey)
	req.Header.Set("Accept-Language", "pt-BR,en;q=0.9")

	fields := []string{
		"id", "displayName", "formattedAddress", "location", "businessStatus", "types",
		"nationalPhoneNumber", "websiteUri",
		"iconMaskBaseUri", "iconBackgroundColor", "photos",
		"rating", "userRatingCount", "googleMapsUri",
		"editorialSummary", "reviews", "generativeSummary",
	}

	fieldMask := ""
	for i, f := range fields {
		if i > 0 {
			fieldMask += ","
		}
		fieldMask += f
	}
	req.Header.Set("X-Goog-FieldMask", fieldMask)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google places v1 error: status %d", resp.StatusCode)
	}

	var item newPlaceDetailResult
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	var photoURL string
	if len(item.Photos) > 0 {
		photoURL = r.buildNewPhotoURL(item.Photos[0].Name, 800)
	}

	iconURL := item.IconMaskBaseUri
	if iconURL != "" {
		iconURL = iconURL + ".svg"
	}

	var reviews []entity.Review
	for _, rev := range item.Reviews {
		reviews = append(reviews, entity.Review{
			Author:      rev.AuthorAttribution.DisplayName,
			Rating:      rev.Rating,
			Text:        rev.OriginalText.Text, // Ou rev.Text.Text
			Time:        rev.RelativePublishTimeDescription,
			PublishTime: rev.PublishTime,
		})
	}

	finalSummary := item.EditorialSummary.Text

	if finalSummary == "" {
		if item.GenerativeSummary.Overview.Text != "" {
			finalSummary = item.GenerativeSummary.Overview.Text
		} else if item.GenerativeSummary.Description.Text != "" {
			finalSummary = item.GenerativeSummary.Description.Text
		}
	}

	return &entity.Company{
		ID:                  item.ID,
		Name:                item.DisplayName.Text,
		Address:             item.FormattedAddress,
		Latitude:            item.Location.Latitude,
		Longitude:           item.Location.Longitude,
		BusinessStatus:      item.BusinessStatus,
		Types:               item.Types,
		Phone:               item.NationalPhoneNumber,
		Website:             item.WebsiteUri,
		Summary:             finalSummary,
		Rating:              item.Rating,
		UserRatingCount:     item.UserRatingCount,
		GoogleMapsURI:       item.GoogleMapsUri,
		Reviews:             reviews,
		PhotoURL:            photoURL,
		IconMaskBaseURI:     iconURL,
		IconBackgroundColor: item.IconBackgroundColor,
	}, nil
}

func (r *googlePlaceRepo) List(ctx context.Context, limit, offset int) ([]*entity.Company, error) {
	return []*entity.Company{}, nil
}

func (r *googlePlaceRepo) Count(ctx context.Context) (int, error) {
	return 0, nil
}
