package users

type CreateUser struct {
	Name           string
	Email          string
	HashedPassword string
}
