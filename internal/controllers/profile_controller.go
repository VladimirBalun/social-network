package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"social_network/internal/entities"

	"go.uber.org/zap"
)

type ProfilesService interface {
	SaveProfile(context.Context, *entities.Profile) error
	GetProfiles(context.Context) ([]entities.Profile, error)
}

type ProfileController struct {
	service ProfilesService
	logger  *zap.Logger
}

func NewProfileController(service ProfilesService, logger *zap.Logger) ProfileController {
	return ProfileController{
		service: service,
		logger:  logger,
	}
}

type profile struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Surname   string   `json:"surname"`
	City      string   `json:"city"`
	Interests []string `json:"interests"`
	Age       int8     `json:"age"`
	Gender    int8     `json:"gender"`
}

type saveProfileRequest struct {
	profile
}

func (c *ProfileController) SaveProfile(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.Warn("failed to parse request body")
		return
	}

	var request saveProfileRequest
	if err = json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.logger.Warn("failed to unmarshal request body")
		return
	}

	profile := entities.Profile{}
	profile.ID = request.ID
	profile.Name = request.Name
	profile.Surname = request.Surname
	profile.City = request.City
	profile.Interests = request.Interests
	profile.Age = request.Age
	profile.Gender = request.Gender

	if err = c.service.SaveProfile(r.Context(), &profile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type getProfilesResponse struct {
	Profiles []profile `json:"profiles"`
}

func (c *ProfileController) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := c.service.GetProfiles(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.Error(err.Error())
		return
	}

	var response getProfilesResponse
	response.Profiles = make([]profile, len(profiles))
	for i := range profiles {
		response.Profiles[i] = profile{
			ID:        profiles[i].ID,
			Name:      profiles[i].Name,
			Surname:   profiles[i].Surname,
			City:      profiles[i].City,
			Interests: profiles[i].Interests,
			Age:       profiles[i].Age,
			Gender:    profiles[i].Gender,
		}
	}

	body, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.logger.Error("failed to marshal response body")
		return
	}

	if _, err = w.Write(body); err != nil {
		c.logger.Error("failed to write response body")
	}
}
