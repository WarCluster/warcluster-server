package server

import (
	"errors"

	"github.com/Vladimiroff/vec2d"
)

// Request type hold player's requests data
type Request struct {
	Client       *Client         // Reference to the client who sent this. Populated by the server.
	Command      string          // Used to parse the request (required)
	Type         string          // Type of mission (possible values are: Attach, Supply, Spy)
	Position     *vec2d.Vector   // Position of the client when he sent the request
	Resolution   []uint64        // Clients resolution
	StartPlanets []string        // Planets from which to start a mission
	Path         []*vec2d.Vector // All intermidiate points (waypoints) that define the missions path
	EndPlanet    string          // Mission's destination
	Fleet        int32           // Percentge of ships to be sent in the start mission request
	Username     string          // Client's username needed while loggin in
	TwitterID    string          // Client's twitter id needed while logging in
	Race         uint8           // Race ID chosen during registration
	SunTextureId uint16          // Sun Texture ID chosen during registration
}

// ParseRequest is serving the purpouse of a request manager. Determines the
// type of the request and will return a function that will manage it.
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
