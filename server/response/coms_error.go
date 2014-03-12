package response

import (
	"warcluster/entities"
)

type ComsError struct {
	baseResponse
	Message string
}

func NewComsError(message string) *ComsError {
	c := new(ComsError)
	c.Command = "error"
	c.Message = message
	return c
}

func (m *ComsError) Sanitize(*entities.Player) {}
