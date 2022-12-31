package repositories

import (
	"context"
	"database/sql"
	"reflect"
	"social_network/internal/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProfileRepository(t *testing.T) {
	t.Parallel()

	database, err := sql.Open("mysql", "root:mauFJcuf5dhRMQrjj@tcp(127.0.0.1:3306)/social_network")
	require.NoError(t, err)

	defer func(database *sql.DB) {
		err = database.Close()
		require.NoError(t, err)
	}(database)

	profile := entities.Profile{
		ID:        100,
		Name:      "Ivan",
		Surname:   "Petrov",
		City:      "Moscow",
		Interests: []string{"football", "tennis"},
		Age:       25,
		Gender:    entities.MaleGender,
	}

	ctx := context.Background()
	repository := NewProfileRepository(database, zap.NewNop())

	err = repository.SaveProfile(context.Background(), &profile)
	assert.NoError(t, err)

	profiles, err := repository.GetProfiles(ctx)
	require.NoError(t, err)

	assert.Equal(t, 1, len(profiles))
	assert.Equal(t, profile.ID, profiles[0].ID)
	assert.Equal(t, profile.Name, profiles[0].Name)
	assert.Equal(t, profile.Surname, profiles[0].Surname)
	assert.Equal(t, profile.City, profiles[0].City)
	assert.Equal(t, profile.Age, profiles[0].Age)
	assert.Equal(t, profile.Gender, profiles[0].Gender)
	assert.True(t, reflect.DeepEqual(profile.Interests, profiles[0].Interests))
}
