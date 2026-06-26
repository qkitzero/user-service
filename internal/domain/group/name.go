package group

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const maxGroupNameLength = 255

type GroupName string

func (n GroupName) String() string {
	return string(n)
}

func NewGroupName(s string) (GroupName, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return GroupName(""), fmt.Errorf("invalid group name")
	}
	if utf8.RuneCountInString(s) > maxGroupNameLength {
		return GroupName(""), fmt.Errorf("group name is too long")
	}
	return GroupName(s), nil
}
