package hasher

type Hasher interface {
	Hash(password string) (string, error)
	Match(password, hashedPassword string) bool
}
