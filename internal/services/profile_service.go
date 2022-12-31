package services

import (
	"context"
	"social_network/internal/entities"

	"github.com/pkg/errors"
)

type ProfileRepository interface {
	SaveProfile(context.Context, *entities.Profile) error
	GetProfiles(context.Context) ([]entities.Profile, error)
}

type ProfileService struct {
	repository ProfileRepository
}

func NewProfileService(repository ProfileRepository) ProfileService {
	return ProfileService{
		repository: repository,
	}
}

func (s *ProfileService) SaveProfile(ctx context.Context, profile *entities.Profile) error {
	if err := s.repository.SaveProfile(ctx, profile); err != nil {
		return errors.Wrap(err, "failed to save profile")
	}

	return nil
}

func (s *ProfileService) GetProfiles(ctx context.Context) ([]entities.Profile, error) {
	profiles, err := s.repository.GetProfiles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get profiles")
	}

	return profiles, err
}
