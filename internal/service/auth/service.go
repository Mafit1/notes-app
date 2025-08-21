package auth

import (
	jwt_service "github.com/Mafit1/notes-app/internal/service/jwtservice"
	users_service "github.com/Mafit1/notes-app/internal/service/users"
)

type service struct {
	usersService users_service.Service
	jwtService   jwt_service.Service
}

func New(usersService users_service.Service, jwtService jwt_service.Service) Service {
	return &service{
		usersService: usersService,
		jwtService:   jwtService,
	}
}

func (s *service) Register(in RegisterIn) (out RegisterOut, err error) {
	panic("unimplemented")
}

func (s *service) Login(in LoginIn) (out LoginOut, err error) {
	panic("unimplemented")
}
