package timescale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/et/internal/homing"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	t.Run("failed to start transactions", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()

		assert.Nil(t, err)

		mock.ExpectBegin().WillReturnError(fmt.Errorf("eny error"))

		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		repo := New(sqlxDB)

		mockTelemetry := homing.Telemetry{
			ID:        uuid.NewString(),
			Services:  []string{},
			Service:   "mock service",
			Longitude: 1.2,
			Latitude:  30.2,
			IpAddress: "192.168.0.1",
			Version:   "0.13",
			LastSeen:  time.Now(),
			Country:   "someCountry",
			City:      "someCity",
		}

		err = repo.Save(context.Background(), mockTelemetry)
		assert.NotNil(t, err)
	})
	t.Run("failed exec", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		assert.Nil(t, err)

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO telemetry").WillReturnError(fmt.Errorf("failed save"))

		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		repo := New(sqlxDB)

		mockTelemetry := homing.Telemetry{
			ID:        uuid.NewString(),
			Services:  []string{},
			Service:   "mock service",
			Longitude: 1.2,
			Latitude:  30.2,
			IpAddress: "192.168.0.1",
			Version:   "0.13",
			LastSeen:  time.Now(),
			Country:   "someCountry",
			City:      "someCity",
		}

		err = repo.Save(context.Background(), mockTelemetry)
		assert.NotNil(t, err)
	})
	t.Run("invalid text representation", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		assert.Nil(t, err)

		mock.ExpectBegin()

		pgerr := pgconn.PgError{
			Code: pgerrcode.InvalidTextRepresentation,
		}

		mock.ExpectExec("INSERT INTO telemetry").WillReturnError(&pgerr)

		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		repo := New(sqlxDB)

		mockTelemetry := homing.Telemetry{
			ID:        uuid.NewString(),
			Services:  []string{},
			Service:   "mock service",
			Longitude: 1.2,
			Latitude:  30.2,
			IpAddress: "192.168.0.1",
			Version:   "0.13",
			LastSeen:  time.Now(),
			Country:   "someCountry",
			City:      "someCity",
		}

		err = repo.Save(context.Background(), mockTelemetry)
		assert.NotNil(t, err)
	})
	t.Run("successful save", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		assert.Nil(t, err)

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO telemetry").WillReturnResult(sqlmock.NewResult(0, 1))

		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		repo := New(sqlxDB)

		mockTelemetry := homing.Telemetry{
			ID:        uuid.NewString(),
			Services:  []string{},
			Service:   "mock service",
			Longitude: 1.2,
			Latitude:  30.2,
			IpAddress: "192.168.0.1",
			Version:   "0.13",
			LastSeen:  time.Now(),
			Country:   "someCountry",
			City:      "someCity",
		}

		err = repo.Save(context.Background(), mockTelemetry)
		assert.Nil(t, err)
	})

}
