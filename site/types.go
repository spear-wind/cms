package site

import (
	"time"

	"github.com/spear-wind/cms/validator"

	"github.com/spear-wind/cms/user"
)

type SiteRepository interface {
	Add(site *Site) (err error)
	Update(site *Site) (err error)
	List() (sites []*Site)
	GetByID(id string) (site *Site, err error)
}

type Site struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	DomainName string     `json:"domain_name"`
	CreatedBy  *user.User `json:"created_by"`
	Created    time.Time  `json:"date_created"`
}

func NewSite(name string, domainName string, createdBy *user.User) *Site {
	return &Site{
		Name:       name,
		DomainName: domainName,
		CreatedBy:  createdBy,
		Created:    time.Now(),
	}
}

func (s *Site) validate() (result validator.ValidationResult) {
	result = validator.NewValidationResult()

	if len(s.Name) == 0 {
		result.AddError("name", "Name is required")
	}

	if len(s.DomainName) == 0 {
		result.AddError("domain_name", "Domain Name is required")
	}

	if s.CreatedBy == nil {
		result.AddError("created_by", "Created by is required")
	}

	if &s.Created == nil {
		result.AddError("date_created", "Date Created is required")
	}

	return result
}
