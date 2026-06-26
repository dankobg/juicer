package game

import "errors"

var (
	ErrGameNotFound      = errors.New("game not found")
	ErrGameAlreadyExists = errors.New("game already exists")
)
