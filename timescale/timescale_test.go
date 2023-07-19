package timescale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/mainflux/callhome"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	ctx := context.TODO()
	mockTelemetry := callhome.Telemetry{
		Services:    []string{},
		Service:     "mock service",
		Longitude:   1.2,
		Latitude:    30.2,
		IpAddress:   "192.168.0.1",
		Version:     "0.13",
		LastSeen:    time.Now(),
		Country:     "someCountry",
		City:        "someCity",
		ServiceTime: time.Now(),
	}
	t.Run("failed to start transactions", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()

		assert.Nil(t, err)

		mock.ExpectBegin().WillReturnError(fmt.Errorf("eny error"))

		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		repo := New(sqlxDB)

		err = repo.Save(ctx, mockTelemetry)
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

		err = repo.Save(ctx, mockTelemetry)
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

		err = repo.Save(ctx, mockTelemetry)
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

		err = repo.Save(ctx, mockTelemetry)
		assert.Nil(t, err)
	})
}

func TestRetrieveAll(t *testing.T) {
	ctx := context.TODO()
	now := time.Now()
	mTel := callhome.Telemetry{
		Service:     "mock service",
		Longitude:   1.2,
		Latitude:    30.2,
		IpAddress:   "192.168.0.1",
		Version:     "0.13",
		LastSeen:    now,
		Country:     "someCountry",
		City:        "someCity",
		ServiceTime: now,
	}
	t.Run("error performing select", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		assert.Nil(t, err)

		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		repo := New(sqlxDB)

		mock.ExpectQuery("SELECT(.*)").WillReturnError(fmt.Errorf("any error"))

		_, err = repo.RetrieveAll(ctx, callhome.PageMetadata{Limit: 10, Offset: 0}, callhome.TelemetryFilters{})
		assert.NotNil(t, err)
	})
	t.Run("successful", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		assert.Nil(t, err)

		defer sqlDB.Close()
		sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

		repo := New(sqlxDB)

		rows := sqlmock.NewRows(
			[]string{"ip_address", "longitude", "latitude", "mf_version", "service", "time", "country", "city", "service_time"},
		).AddRow(mTel.IpAddress, mTel.Longitude, mTel.Latitude, mTel.Version, mTel.Service, mTel.LastSeen, mTel.Country, mTel.City, mTel.ServiceTime)

		rows2 := sqlmock.NewRows(
			[]string{"count"},
		).AddRow(1)

		mock.ExpectQuery("SELECT(.*)").WillReturnRows(rows)
		mock.ExpectQuery("SELECT COUNT(.*) FROM telemetry").WillReturnRows(rows2)

		tp, err := repo.RetrieveAll(ctx, callhome.PageMetadata{Limit: 10, Offset: 0}, callhome.TelemetryFilters{})
		assert.Nil(t, err)
		assert.Equal(t, mTel, tp.Telemetry[0])
	})
}
