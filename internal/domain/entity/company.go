package entity

import (
	"errors"
	"time"
)

type Review struct {
	Author      string `json:"author"`
	Rating      int    `json:"rating"`
	Text        string `json:"text"`
	Time        string `json:"time"`
	PublishTime string `json:"publish_time"`
}

type ContactEmail struct {
	Value      string `json:"value"`
	Type       string `json:"type"`
	Department string `json:"department"`
	Seniority  string `json:"seniority"`
	Position   string `json:"position"`
	Confidence int    `json:"confidence"`
}

type Company struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Address        string   `json:"address"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	BusinessStatus string   `json:"business_status,omitempty"`
	Types          []string `json:"types,omitempty"`
	Website        string   `json:"website,omitempty"`
	Phone          string   `json:"phone_number,omitempty"`
	Summary        string   `json:"summary,omitempty"`

	// Enrichment by Hunter.IO
	EmployeeCount   string   `json:"employee_count,omitempty"`
	Revenue         string   `json:"revenue,omitempty"`
	Industry        string   `json:"industry,omitempty"`
	LinkedInURL     string   `json:"linkedin_url,omitempty"`
	Rating          float64  `json:"rating,omitempty"`
	UserRatingCount int      `json:"user_rating_count,omitempty"`
	GoogleMapsURI   string   `json:"google_maps_uri,omitempty"`
	Reviews         []Review `json:"reviews,omitempty"`

	PhotoURL            string `json:"photo_url,omitempty"`
	IconMaskBaseURI     string `json:"icon_mask_base_uri,omitempty"`
	IconBackgroundColor string `json:"icon_background_color,omitempty"`

	Contacts []ContactEmail `json:"contacts,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Company) TableName() string {
	return "companies"
}

func (c *Company) IsValid() error {
	if c.Name == "" {
		return errors.New("company name is required")
	}
	if c.Latitude == 0 || c.Longitude == 0 {
		return errors.New("invalid coordinates")
	}
	return nil
}
