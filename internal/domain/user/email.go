package user

import (
	"fmt"
	"regexp"
)

type Email string

func NewEmail(s string) (Email, error) {
	const emailRegex = `^[a-z0-9]+(?:[._-][a-z0-9]+)*@[a-z0-9-]+\.[a-z0-9-.]+$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(s) {
		return "", fmt.Errorf("invalid email")
	}
	return Email(s), nil
}
