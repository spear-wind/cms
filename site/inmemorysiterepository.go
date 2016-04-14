package site

import (
	"errors"
	"fmt"
)

type inMemoryRepository struct {
	sites map[string]*Site
}

func NewInMemoryRepository() *inMemoryRepository {
	repo := &inMemoryRepository{}
	repo.sites = make(map[string]*Site)
	return repo
}

func (repo *inMemoryRepository) Add(site *Site) (err error) {
	site.ID = fmt.Sprintf("%d", len(repo.sites)+1)
	repo.sites[site.ID] = site
	return err
}

func (repo *inMemoryRepository) Update(site *Site) (err error) {
	repo.sites[site.ID] = site
	return err
}

func (repo *inMemoryRepository) List() (sites []*Site) {
	for _, site := range repo.sites {
		sites = append(sites, site)
	}

	return sites
}

func (repo *inMemoryRepository) GetByID(siteID string) (site *Site, err error) {
	found := false

	for _, target := range repo.sites {
		if siteID == target.ID {
			site = target
			found = true
		}
	}
	if !found {
		err = errors.New("Could not find site in repository")
	}
	return site, err
}
