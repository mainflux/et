package timescale

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/callhome"
	"github.com/pkg/errors"
)

var _ callhome.TelemetryRepo = (*repo)(nil)

type repo struct {
	db *sqlx.DB
}

// New returns new TimescaleSQL writer.
func New(db *sqlx.DB) callhome.TelemetryRepo {
	return &repo{db: db}
}

// RetrieveAll gets all records from repo.
func (r repo) RetrieveAll(ctx context.Context, pm callhome.PageMetadata) (callhome.TelemetryPage, error) {
	q := `
	WITH aggregated_data AS (
		SELECT ip_address, ARRAY_AGG(DISTINCT service) AS services
		FROM telemetry
		GROUP BY ip_address
	  )
	  SELECT ad.ip_address, ad.services, t.time, t.service_time, t.longitude, t.latitude, t.mf_version, t.country, t.city
	  FROM aggregated_data ad
	  INNER JOIN (
		  SELECT DISTINCT ON (ip_address) *
		  FROM telemetry
		  ORDER BY ip_address, time DESC
		  %s OFFSET :offset
	  ) t ON ad.ip_address = t.ip_address;	  
	`

	switch pm.Limit {
	case 0:
		q = fmt.Sprintf(q, "")
	default:
		q = fmt.Sprintf(q, "LIMIT :limit")
	}

	params := map[string]interface{}{
		"limit":  pm.Limit,
		"offset": pm.Offset,
	}

	rows, err := r.db.NamedQuery(q, params)
	if err != nil {
		return callhome.TelemetryPage{}, err
	}
	defer rows.Close()

	var results callhome.TelemetryPage

	for rows.Next() {
		var result callhome.Telemetry
		if err := rows.StructScan(&result); err != nil {
			return callhome.TelemetryPage{}, err
		}
		results.Telemetry = append(results.Telemetry, result)
	}

	q = `SELECT COUNT(*) FROM telemetry;`
	rows, err = r.db.NamedQuery(q, params)
	if err != nil {
		return callhome.TelemetryPage{}, err
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

// Save creates record in repo.
func (r repo) Save(ctx context.Context, t callhome.Telemetry) error {
	q := `INSERT INTO telemetry (ip_address, longitude, latitude,
		mf_version, service, time, country, city, service_time)
		VALUES (:ip_address, :longitude, :latitude,
			:mf_version, :service, :time, :country, :city, :service_time);`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(ErrSaveEvent, err.Error())
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = errors.Wrap(err, errors.Wrap(ErrTransRollback, txErr.Error()).Error())
			}
			return
		}

		if err = tx.Commit(); err != nil {
			err = errors.Wrap(ErrSaveEvent, err.Error())
		}
	}()

	if _, err := tx.NamedExec(q, t); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.InvalidTextRepresentation {
				return errors.Wrap(ErrSaveEvent, ErrInvalidEvent.Error())
			}
		}
		return errors.Wrap(ErrSaveEvent, err.Error())
	}
	return nil

}

// RetrieveDistinctIPsCountries retrieve distinct
func (r repo) RetrieveDistinctIPsCountries(ctx context.Context) (callhome.TelemetrySummary, error) {
	q := `select distinct ip_address from telemetry;`
	rows, err := r.db.Query(q)
	if err != nil {
		return callhome.TelemetrySummary{}, err
	}
	defer rows.Close()
	var summary callhome.TelemetrySummary
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return summary, err
		}
		summary.IpAddresses = append(summary.IpAddresses, val)
	}

	q = `select distinct country from telemetry;`
	rows, err = r.db.Query(q)
	if err != nil {
		return summary, err
	}
	defer rows.Close()
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return summary, err
		}
		summary.Countries = append(summary.Countries, val)
	}
	return summary, nil
}
