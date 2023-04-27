package timescale

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // required for SQL access
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

// Config defines the options that are used when connecting to a TimescaleSQL instance
type Config struct {
	Host        string `env:"TIMESCALE_HOST"            envDefault:"localhost"`
	Port        string `env:"TIMESCALE_PORT"            envDefault:"5432"`
	User        string `env:"TIMESCALE_USER"            envDefault:"postgres"`
	Pass        string `env:"TIMESCALE_PASSWORD"        envDefault:"password"`
	Name        string `env:"TIMESCALE_DB_NAME"         envDefault:"postgres"`
	SSLMode     string `env:"TIMESCALE_SSL_MODE"        envDefault:"disable"`
	SSLCert     string `env:"TIMESCALE_SSL_CERT"        envDefault:""`
	SSLKey      string `env:"TIMESCALE_SSL_KEY"         envDefault:""`
	SSLRootCert string `env:"TIMESCALE_SSL_ROOT_CERT"   envDefault:""`
}

// Connect creates a connection to the TimescaleSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

// Migration of telemetry service
func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "telemetry_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS telemetry (
						id			UUID,
						ip_address	TEXT	NOT	NULL,
					 	longitude 	FLOAT	NOT	NULL,
						latitude	FLOAT	NOT NULL,
						mf_version	TEXT,
						service		TEXT,
						last_seen	TIMESTAMPTZ,
						country 	TEXT,
						city 		TEXT,
						PRIMARY KEY (id)
					)`,
				},
				Down: []string{"DROP TABLE telemetry;"},
			},
		},
	}
	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	return err
}
