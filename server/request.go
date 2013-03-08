package server

import (
	"encoding/json"
	"errors"
	"../entities"
	"net"
)

type Request struct {
	Command       string
	Position      []int
	Resolution    []int
	StartPlanet   string
	EndPlanet     string
	Fleet         int
}

func UnmarshalRequest(input string) (*Request, error) {
	var request *Request
	if err := json.Unmarshal([]byte(input), &request); err != nil {
		return nil, err
	}
	return request, nil
}

func ParseRequest(request *Request) (func (chan<- string, net.Conn, *entities.Player, *Request) error , error) {
	switch request.Command {
	case "start_mission":
		if len(request.StartPlanet) > 0 && len(request.EndPlanet) > 0 {
			if request.Fleet <= 0 {
				request.Fleet = 50
			}
			return actionParser, nil
		} else {
			return nil, errors.New("Not enough arguments")
		}
	case "scope_of_view":
		if len(request.Position) > 0 && len(request.Resolution) > 0 {
			return scopeOfView, nil
		} else {
			return nil, errors.New("Not enough arguments")
		}
	}
	return nil, nil
}
