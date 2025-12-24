package entity

import (
	"errors"
	"time"
)

type Company struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	EmployeeCount string  `json:"employee_count,omitempty"`
	Revenue       string  `json:"revenue,omitempty"`
	Industry      string  `json:"industry,omitempty"`
	LinkedInURL   string  `json:"linkedin_url"`

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
