package auth

import "context"

type Service interface {
	Register(ctx context.Context, in RegisterIn) (*RegisterOut, error)
	Login(ctx context.Context, in LoginIn) (*LoginOut, error)
}
