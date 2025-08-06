package postgres

import "time"

type Option func(*Postgres)

func ConnAttempts(a int) Option {
	return func(p *Postgres) {
		p.connectionAttempts = a
	}
}

func TimeOut(t time.Duration) Option {
	return func(p *Postgres) {
		p.connectionTimeout = t
	}
}
