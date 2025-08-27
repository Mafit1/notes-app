package jwtservice

import "context"

type Service interface {
	GeneratePair(ctx context.Context, in GenerateIn) (*GenerateOut, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(tokenString string) (bool, error)
}
