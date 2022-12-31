package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"social_network/internal/controllers/mock"
	"social_network/internal/entities"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSaveProfile(t *testing.T) {
	tests := map[string]struct {
		stubs func(store *mock.MockProfilesService)
		code  int
	}{
		"save profile with error from service": {
			stubs: func(service *mock.MockProfilesService) {
				service.EXPECT().
					SaveProfile(gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			code: http.StatusInternalServerError,
		},
		"successfully save profile": {
			stubs: func(service *mock.MockProfilesService) {
				service.EXPECT().
					SaveProfile(gomock.Any(), &entities.Profile{
						ID:        100,
						Name:      "Ivan",
						Surname:   "Vetrov",
						City:      "Rostov",
						Interests: []string{"games", "poker"},
						Age:       30,
						Gender:    entities.MaleGender,
					}).
					Return(nil)
			},
			code: http.StatusCreated,
		},
	}

	ctl := gomock.NewController(t)
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			service := mock.NewMockProfilesService(ctl)
			tc.stubs(service)

			requestData := profile{
				ID:        100,
				Name:      "Ivan",
				Surname:   "Vetrov",
				City:      "Rostov",
				Interests: []string{"games", "poker"},
				Age:       30,
				Gender:    entities.MaleGender,
			}

			requestBody, err := json.Marshal(requestData)
			require.NoError(t, err)

			writer := httptest.NewRecorder()
			requestReader := strings.NewReader(string(requestBody))
			request := httptest.NewRequest(http.MethodPost, "/profile", requestReader)

			controller := NewProfileController(service, zap.NewNop())
			controller.SaveProfile(writer, request)
			assert.Equal(t, tc.code, writer.Code)
		})
	}
}

func TestGetProfiles(t *testing.T) {
	tests := map[string]struct {
		stubs    func(store *mock.MockProfilesService)
		response getProfilesResponse
		code     int
	}{
		"get profiles with error from service": {
			stubs: func(service *mock.MockProfilesService) {
				service.EXPECT().
					GetProfiles(gomock.Any()).
					Return(nil, errors.New("error"))
			},
			response: getProfilesResponse{},
			code:     http.StatusInternalServerError,
		},
		"successfully get profiles": {
			stubs: func(service *mock.MockProfilesService) {
				service.EXPECT().
					GetProfiles(gomock.Any()).
					Return([]entities.Profile{
						{
							ID:        100,
							Name:      "Ivan",
							Surname:   "Vetrov",
							City:      "Rostov",
							Interests: []string{"games", "poker"},
							Age:       30,
							Gender:    entities.MaleGender,
						},
						{
							ID:        200,
							Name:      "Vladimir",
							Surname:   "Ivanov",
							City:      "Moscow",
							Interests: []string{"music"},
							Age:       20,
							Gender:    entities.MaleGender,
						},
					}, nil)
			},
			response: getProfilesResponse{
				Profiles: []profile{
					{
						ID:        100,
						Name:      "Ivan",
						Surname:   "Vetrov",
						City:      "Rostov",
						Interests: []string{"games", "poker"},
						Age:       30,
						Gender:    entities.MaleGender,
					},
					{
						ID:        200,
						Name:      "Vladimir",
						Surname:   "Ivanov",
						City:      "Moscow",
						Interests: []string{"music"},
						Age:       20,
						Gender:    entities.MaleGender,
					},
				},
			},
			code: http.StatusOK,
		},
	}

	ctl := gomock.NewController(t)
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			service := mock.NewMockProfilesService(ctl)
			tc.stubs(service)

			writer := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/profiles", nil)

			controller := NewProfileController(service, zap.NewNop())
			controller.GetProfiles(writer, request)
			assert.Equal(t, tc.code, writer.Code)

			if tc.code != http.StatusInternalServerError {
				result := writer.Result()
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					require.NoError(t, err)
				}(result.Body)

				responseBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				var response getProfilesResponse
				err = json.Unmarshal(responseBody, &response)
				require.NoError(t, err)
				assert.Equal(t, tc.response, response)
			}
		})
	}
}
