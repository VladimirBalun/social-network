package services

import (
	"context"
	"reflect"
	"social_network/internal/entities"
	"social_network/internal/services/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSaveProfile(t *testing.T) {
	profile := entities.Profile{
		ID:        1,
		Name:      "Ivan",
		Surname:   "Petrov",
		City:      "Moscow",
		Interests: []string{"football", "piano"},
		Age:       25,
		Gender:    entities.MaleGender,
	}

	ctx := context.Background()
	tests := map[string]struct {
		stubs   func(store *mock.MockProfileRepository, profile *entities.Profile)
		profile *entities.Profile
		hasErr  bool
	}{
		"save profile with error from repository": {
			stubs: func(repository *mock.MockProfileRepository, profile *entities.Profile) {
				repository.EXPECT().
					SaveProfile(ctx, profile).
					Return(errors.New("error"))
			},
			profile: &profile,
			hasErr:  true,
		},
		"successfully save profile": {
			stubs: func(repository *mock.MockProfileRepository, profile *entities.Profile) {
				repository.EXPECT().
					SaveProfile(ctx, profile).
					Return(nil)
			},
			profile: &profile,
			hasErr:  false,
		},
	}

	ctl := gomock.NewController(t)
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repository := mock.NewMockProfileRepository(ctl)
			tc.stubs(repository, tc.profile)

			service := NewProfileService(repository)
			err := service.SaveProfile(ctx, tc.profile)
			assert.Equal(t, tc.hasErr, err != nil)
		})
	}
}

func TestGetProfiles(t *testing.T) {
	profiles := []entities.Profile{
		{
			ID:        1,
			Name:      "Ivan",
			Surname:   "Petrov",
			City:      "Moscow",
			Interests: []string{"football", "piano"},
			Age:       25,
			Gender:    entities.MaleGender,
		},
		{
			ID:        2,
			Name:      "Alexandra",
			Surname:   "Ivanova",
			City:      "Smolensk",
			Interests: []string{"books"},
			Age:       33,
			Gender:    entities.FemaleGender,
		},
	}

	ctx := context.Background()
	tests := map[string]struct {
		stubs    func(store *mock.MockProfileRepository)
		profiles []entities.Profile
		hasErr   bool
	}{
		"get profiles with error from repository": {
			stubs: func(repository *mock.MockProfileRepository) {
				repository.EXPECT().
					GetProfiles(ctx).
					Return(nil, errors.New("error"))
			},
			profiles: nil,
			hasErr:   true,
		},
		"successfully get profiles": {
			stubs: func(repository *mock.MockProfileRepository) {
				repository.EXPECT().
					GetProfiles(ctx).
					Return(profiles, nil)
			},
			profiles: profiles,
			hasErr:   false,
		},
	}

	ctl := gomock.NewController(t)
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repository := mock.NewMockProfileRepository(ctl)
			tc.stubs(repository)

			service := NewProfileService(repository)
			result, err := service.GetProfiles(ctx)
			assert.True(t, reflect.DeepEqual(tc.profiles, result))
			assert.Equal(t, tc.hasErr, err != nil)
		})
	}
}
