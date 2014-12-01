package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Vladimiroff/vec2d"
)

// The Request struct is used to created to make the request manipulation easier by creating a template
// that can hold the information needed by all request types.
// 1. In client we keep the reference to the client struct to be able to edit data and return feedback
// 	  to the connection.
// 2. The command is the container for a key word used to swith between different input Request.
// 3. The Position field is used as a general position container for the different requests.
// 4. Resolution I think speaks for itself. Its the size of the screan.
// 5. Start and End plane are containers for general planet information(mostly used for mission requests).
// 6. Fleet contains the percent of ships to be sent in the start mission requests.
type Request struct {
	Client       *Client
	Command      string
	Type         string
	Position     *vec2d.Vector
	Resolution   []uint64
	StartPlanets []string
	EndPlanet    string
	Fleet        int32
	Username     string
	TwitterID    string
	Race         uint8
	SunTextureId uint16
}

// This function transfers the information of a new request from byte list
// to the special Request struct.
// Since the data is formated as json the function uses unmarshal to extract the data.
func UnmarshalRequest(message []byte, client *Client) (*Request, error) {
	var request Request

	if err := json.Unmarshal(message, &request); err != nil {
		log.Println("Error in server.request.UnmarshalRequest:", err.Error())
		return nil, err
	}

	request.Client = client
	return &request, nil
}

// ParseRequest is serving the purpouse of a request manager.
// After the request is parsed to the more usable Request struct
// ParseRequest will determine the type of the request and will
// return a function that will manage it.
func ParseRequest(request *Request) (func(*Request) error, error) {
	switch request.Command {
	case "start_mission":
		if len(request.StartPlanets) > 0 && len(request.EndPlanet) > 0 {
			if request.Fleet > 100 || request.Fleet <= 0 {
				request.Fleet = 100
			}
			return parseAction, nil
		} else {
			return nil, errors.New("Not enough arguments")
		}
	case "scope_of_view":
		if request.Position != nil && len(request.Resolution) > 0 {
			return scopeOfView, nil
		} else {
			return nil, errors.New("Not enough arguments")
		}
	case "voronoi_diagram":
		if request.Position != nil && len(request.Resolution) > 0 {
			return voronoiDiagram, nil
		} else {
			return nil, errors.New("Not enough arguments")
		}
	}
	return nil, errors.New("Unknown command")
}
