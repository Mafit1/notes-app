package models

import "github.com/google/uuid"

type RoleType string

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
	Role     RoleType
}

const (
	RoleTypeAdmin RoleType = "admin"
	RoleTypeUser  RoleType = "user"
)
