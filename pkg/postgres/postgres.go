package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultConnectionTimeout  = time.Second
	defaultConnectionAttempts = 10
)

type Postgres struct {
	connectionTimeout  time.Duration
	connectionAttempts int

	Pool *pgxpool.Pool
}

func New(url string) (*Postgres, error) {
	pg := &Postgres{
		connectionTimeout:  defaultConnectionTimeout,
		connectionAttempts: defaultConnectionAttempts,
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	attempt := pg.connectionAttempts

	for attempt > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)

		if err != nil {
			return nil, fmt.Errorf("pool creation error: %w", err)
		}
		if err = pg.Pool.Ping(context.Background()); err != nil {
			break
		}

		attempt--
		time.Sleep(pg.connectionTimeout)
	}

	if err != nil {
		return nil, fmt.Errorf("postgres connection error, out of connection attempts: %w", err)
	}

	return pg, nil
}

func (pg *Postgres) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
