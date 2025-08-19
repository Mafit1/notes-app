package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func (u *User) CheckPassword(plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)) == nil
}
