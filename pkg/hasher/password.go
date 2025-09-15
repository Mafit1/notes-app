package hasher

import "golang.org/x/crypto/bcrypt"

type bcryptHasher struct{}

func NewBcrypt() Hasher {
	return &bcryptHasher{}
}

func (h *bcryptHasher) Hash(value string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *bcryptHasher) Match(value, hashedValue string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value)); err != nil {
		return false
	}

	return true
}
