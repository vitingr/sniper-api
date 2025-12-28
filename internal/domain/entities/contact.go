package entity

import (
	"errors"
	"strings"
	"time"
)

type Contact struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"company_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`

	Title      string `json:"title"`
	Seniority  string `json:"seniority"`
	Department string `json:"department"`

	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DirectPhone string `json:"direct_phone"`

	LinkedInURL string `json:"linkedin_url"`
	TwitterURL  string `json:"twitter_url"`
	PhotoURL    string `json:"photo_url"`

	EmailStatus string `json:"email_status"` // Ex: "verified", "risky", "unknown" (Vem do enriquecimento)
	Source      string `json:"source"`       // Ex: "apollo", "clearbit", "manual"

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Contact) TableName() string {
	return "contacts"
}

func (c *Contact) IsValid() error {
	if c.CompanyID == 0 {
		return errors.New("company_id is required")
	}

	if c.FirstName == "" {
		return errors.New("First name is required")
	}

	if c.Email == "" && !strings.Contains(c.Email, "@") {
		return errors.New("Invalid email format")
	}

	return nil
}
