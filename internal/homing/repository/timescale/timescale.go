package timescale

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/et/internal/homing"
	"github.com/mainflux/et/internal/homing/repository"
	"github.com/mainflux/mainflux/readers"
	"github.com/pkg/errors"
)

var _ homing.TelemetryRepo = (*repo)(nil)

type repo struct {
	db *sqlx.DB
}

// New returns new TimescaleSQL writer.
func New(db *sqlx.DB) homing.TelemetryRepo {
	return &repo{db: db}
}

// RetrieveAll implements homing.TelemetryRepo
func (r repo) RetrieveAll(ctx context.Context, pm homing.PageMetadata) (homing.TelemetryPage, error) {
	q := `SELECT * FROM telemetry LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := r.db.NamedQuery(q, params)
	if err != nil {
		return homing.TelemetryPage{}, errors.Wrap(readers.ErrReadMessages, err.Error())
	}
	defer rows.Close()

	var results homing.TelemetryPage

	for rows.Next() {
		var result homing.Telemetry
		if err := rows.StructScan(&result); err != nil {
			return homing.TelemetryPage{}, errors.Wrap(readers.ErrReadMessages, err.Error())
		}

		results.Telemetry = append(results.Telemetry, result)
	}

	q = `SELECT COUNT(*) FROM telemetry;`
	rows, err = r.db.NamedQuery(q, params)
	if err != nil {
		return homing.TelemetryPage{}, errors.Wrap(readers.ErrReadMessages, err.Error())
	}
	defer rows.Close()

	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return results, err
		}
	}
	results.Total = total

	return results, nil
}

// RetrieveByIP implements homing.TelemetryRepo
func (repo) RetrieveByIP(ctx context.Context, email string) (*homing.Telemetry, error) {
	return &homing.Telemetry{}, repository.ErrRecordNotFound
}

// Save implements homing.TelemetryRepo
func (r repo) Save(ctx context.Context, t homing.Telemetry) error {
	q := `INSERT INTO telemetry (id, ip_address, longitude, latitude,
		mf_version, service, last_seen, country, city)
		VALUES (:id, :ip_address, :longitude, :latitude,
			:mf_version, :service, :last_seen, :country, :city);`

	tx, err := r.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(repository.ErrSaveEvent, err.Error())
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = errors.Wrap(err, errors.Wrap(repository.ErrTransRollback, txErr.Error()).Error())
			}
			return
		}

		if err = tx.Commit(); err != nil {
			err = errors.Wrap(repository.ErrSaveEvent, err.Error())
		}
	}()

	if _, err := tx.NamedExec(q, t); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.InvalidTextRepresentation {
				return errors.Wrap(repository.ErrSaveEvent, repository.ErrInvalidEvent.Error())
			}
		}
		return errors.Wrap(repository.ErrSaveEvent, err.Error())
	}
	return nil

}

// UpdateTelemetry implements homing.TelemetryRepo
func (repo) UpdateTelemetry(ctx context.Context, u homing.Telemetry) error {
	return nil
}
