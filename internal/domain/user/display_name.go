package user

import (
	"fmt"
	"strings"
)

type DisplayName string

func (d DisplayName) String() string {
	return string(d)
}

func NewDisplayName(s string) (DisplayName, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("invalid display name")
	}
	return DisplayName(s), nil
}
