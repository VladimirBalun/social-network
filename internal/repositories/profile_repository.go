package repositories

import (
	"context"
	"database/sql"
	"social_network/internal/entities"
	"strings"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ProfileRepository struct {
	database *sql.DB
	logger   *zap.Logger
}

func NewProfileRepository(database *sql.DB, logger *zap.Logger) ProfileRepository {
	return ProfileRepository{
		database: database,
		logger:   logger,
	}
}

func (r *ProfileRepository) SaveProfile(ctx context.Context, profile *entities.Profile) error {
	tx, err := r.database.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to create tx")
	}

	profileID, err := r.insertProfile(ctx, tx, profile)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			r.logger.Error("failed to rollback tx: " + txErr.Error())
		}

		return errors.Wrap(err, "failed to insert profile")
	}

	if err = r.insertProfileInterests(ctx, tx, profileID, profile.Interests); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			r.logger.Error("failed to rollback tx: " + txErr.Error())
		}

		return errors.Wrap(err, "failed to insert profile interests")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}

	return nil
}

func (r *ProfileRepository) insertProfile(ctx context.Context, tx *sql.Tx, profile *entities.Profile) (int64, error) {
	query := squirrel.StatementBuilder.
		Insert("profile").
		Columns("user_id", "name", "surname", "city", "age", "gender").
		Values(profile.ID, profile.Name, profile.Surname, profile.City, profile.Age, profile.Gender).
		RunWith(tx)

	if _, err := query.ExecContext(ctx); err != nil {
		return 0, errors.Wrap(err, "failed to exec insert profile query")
	}

	return profile.ID, nil
}

func (r *ProfileRepository) insertProfileInterests(ctx context.Context, tx *sql.Tx, profileID int64, interests []string) error {
	query := squirrel.StatementBuilder.
		Insert("profile_interest").
		Columns("user_id", "name").
		RunWith(tx)

	for _, interest := range interests {
		query = query.Values(profileID, interest)
	}

	if _, err := query.ExecContext(ctx); err != nil {
		return errors.Wrap(err, "failed to exec insert profile interests query")
	}

	return nil
}

func (r *ProfileRepository) GetProfiles(ctx context.Context) ([]entities.Profile, error) {
	query := squirrel.StatementBuilder.
		Select("p.user_id", "p.name", "p.surname", "p.city", "p.age", "p.gender", "group_concat(pi.name) AS interests").
		From("profile p").
		InnerJoin("profile_interest pi ON p.user_id = pi.user_id").
		GroupBy("p.user_id").
		RunWith(r.database)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec select profiles query")
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			r.logger.Error("failed to close rows")
		}
	}(rows)

	profiles := make([]entities.Profile, 0)

	for rows.Next() {
		var interests string
		var profile entities.Profile
		if err = rows.Scan(&profile.ID, &profile.Name, &profile.Surname, &profile.City, &profile.Age, &profile.Gender, &interests); err != nil {
			return nil, errors.Wrap(err, "failed to scan profile")
		}

		profile.Interests = strings.Split(interests, ",")
		profiles = append(profiles, profile)
	}

	return profiles, nil
}
