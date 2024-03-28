package encryption

import "golang.org/x/crypto/bcrypt"

func HashString(input string) (string, error) {
	salt, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(salt), nil
}

func ValidateString(input, retrieved string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(input), []byte(retrieved))
	return err == nil
}
