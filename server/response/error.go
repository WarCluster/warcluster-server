package response

import "warcluster/entities"

type Error struct {
	baseResponse
	Message string
}

func NewError(message string) *Error {
	c := new(Error)
	c.Command = "error"
	c.Message = message
	return c
}

func (m *Error) Sanitize(*entities.Player) {}
