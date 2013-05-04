package server

import (
	"warcluster/db_manager"
	"warcluster/entities"
	"github.com/Vladimiroff/vec2d"
	"fmt"
	"time"
)

func CalculateArrivalTime(start_point, end_point []int, speed int) time.Duration {
	start_vector := vec2d.New(float64(start_point[0]), float64(start_point[1]))
	end_vector := vec2d.New(float64(end_point[0]), float64(end_point[1]))
	distance := end_vector.Substitute(start_vector)
	return time.Duration(time.Duration(distance.Length()/float64(speed)) * time.Second)
}

func StartMissionary(mission entities.Mission) {
	start_entity, err := db_manager.GetEntity(mission.GetStartPlanet())
	end_entity, err := db_manager.GetEntity(mission.EndPlanet)
	start_planet := start_entity.(entities.Planet)
	end_planet := end_entity.(entities.Planet)

	speed := mission.GetSpeed()
	time.Sleep(CalculateArrivalTime(start_planet.GetCoords(), end_planet.GetCoords(), speed))

	result := entities.EndMission(end_planet, mission)
	key, serialized_planet, err := result.Serialize()
	if err == nil {
		db_manager.SetEntity(result)
		sessions.Broadcast([]byte(fmt.Sprintf("{%s: %s}", key, serialized_planet)))
	}
	db_manager.DeleteEntity(mission.GetKey())
}
