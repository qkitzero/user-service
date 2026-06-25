package group

import "errors"

var (
	ErrGroupNotFound     = errors.New("group not found")
	ErrCircularReference = errors.New("circular reference")
	ErrAlreadyChild      = errors.New("already a child")
)
