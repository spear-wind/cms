package site

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/events"
	"github.com/unrolled/render"
)

func InitRoutes(router *mux.Router, formatter *render.Render, siteRepository SiteRepository, eventPublisher events.EventPublisher) {
	router.HandleFunc("/site", createSiteHandler(formatter, siteRepository, eventPublisher)).Methods("POST")
	router.HandleFunc("/site", getSiteListHandler(formatter, siteRepository)).Methods("GET")
	router.HandleFunc("/site/{id}", getSiteHandler(formatter, siteRepository)).Methods("GET")
}

func createSiteHandler(formatter *render.Render, siteRepository SiteRepository, eventPublisher events.EventPublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var site Site

		err := json.Unmarshal(payload, &site)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse create site request")
			return
		}

		if result := site.validate(); result.HasErrors() {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": result.Errors,
			})
			return
		}

		if err := siteRepository.Add(&site); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"site":  site,
				"error": err.Error(),
			})
		} else {
			w.Header().Add("Location", fmt.Sprintf("/site/%d", site.ID))
			formatter.JSON(w, http.StatusCreated, site)
			//TODO newSiteCreatedEvent(user)
		}
	}
}

func getSiteListHandler(formatter *render.Render, siteRepository SiteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sites := siteRepository.List()

		formatter.JSON(w, http.StatusOK, map[string]interface{}{
			"sites": sites,
			"total": len(sites),
		})
	}
}

func getSiteHandler(formatter *render.Render, siteRepository SiteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		siteID := vars["id"]

		if site, err := siteRepository.GetByID(siteID); err != nil {
			formatter.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			formatter.JSON(w, http.StatusOK, site)
		}
	}
}
