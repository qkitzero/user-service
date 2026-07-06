package user

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const maxDisplayNameLength = 255

type DisplayName string

func (d DisplayName) String() string {
	return string(d)
}

func NewDisplayName(s string) (DisplayName, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return DisplayName(""), fmt.Errorf("invalid display name")
	}
	if utf8.RuneCountInString(s) > maxDisplayNameLength {
		return DisplayName(""), fmt.Errorf("display name is too long")
	}
	return DisplayName(s), nil
}
