package servers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"social_network/internal/servers/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRESTServerRouting(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	profileController := mock.NewMockProfileController(ctl)
	profileController.EXPECT().
		SaveProfile(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("SaveProfile"))
			require.NoError(t, err)
		})
	profileController.EXPECT().
		GetProfiles(gomock.Any(), gomock.Any()).
		Do(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("GetProfiles"))
			require.NoError(t, err)
		})

	serverAddress := "0.0.0.0:5555"
	httpServer := NewRESTServer(serverAddress, nil, profileController, zap.NewNop())
	go func() {
		httpServer.ListenAndServe()
	}()

	httpClient := &http.Client{}

	saveProfileAddress := fmt.Sprintf("http://%s/profile", serverAddress)
	saveProfileRequest, err := http.NewRequest("POST", saveProfileAddress, nil)
	require.NoError(t, err)

	response, err := httpClient.Do(saveProfileRequest)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, []byte("SaveProfile"), body)

	getProfilesAddress := fmt.Sprintf("http://%s/profiles", serverAddress)
	getProfileRequest, err := http.NewRequest("GET", getProfilesAddress, nil)
	require.NoError(t, err)

	response, err = httpClient.Do(getProfileRequest)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err = io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, []byte("GetProfiles"), body)

	ctx := context.Background()
	httpServer.Shutdown(ctx)
}
