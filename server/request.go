package server

import (
	"encoding/json"
	"errors"
	"log"
)

type Request struct {
	Client		  *Client
	Command       string
	Position      []int
	Resolution    []int
	StartPlanet   string
	EndPlanet     string
	Fleet         int
}

func UnmarshalRequest(message []byte, client *Client) (*Request, error) {
	var request Request

	if err := json.Unmarshal(message, &request); err != nil {
		log.Println("Error in server.request.UnmarshalRequest:", err.Error())
		return nil, err
	}

	request.Client = client
	return &request, nil
}

func ParseRequest(request *Request) (func (*Request) error, error) {
	switch request.Command {
	case "start_mission":
		if len(request.StartPlanet) > 0 && len(request.EndPlanet) > 0 {
			if request.Fleet <= 0 {
				request.Fleet = 50
			}
			return parseAction, nil
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
