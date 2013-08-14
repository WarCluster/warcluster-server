package server

import (
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"time"
	"warcluster/db_manager"
	"warcluster/entities"
)

// The CalculateArrivalTime is used by the mission starter (StartMissionary) to calculate the mission duration.
func CalculateArrivalTime(start_point, end_point []int, speed int) time.Duration {
	start_vector := vec2d.New(float64(start_point[0]), float64(start_point[1]))
	end_vector := vec2d.New(float64(end_point[0]), float64(end_point[1]))
	distance := vec2d.Sub(end_vector, start_vector)
	return time.Duration(time.Duration(distance.Length()/float64(speed)) * time.Second)
}

// StartMissionary is used when a call to initiate a new mission is rescived.
// 1. The function gets the mission planet information from the DB and makes basic data checks.
// 2. Calls CalculateArrivalTime and sleeps the thread for the returned ammount of time.
// 3. When the delay ends the thread ends the mission calling EndMission
// 4. The end of the mission is bradcasted to all clients and the mission entry is erased from the DB.
func StartMissionary(mission *entities.Mission) {
	source_entity, err := db_manager.GetEntity(fmt.Sprintf("planet.%d_%d", mission.Source[0], mission.Source[1]))
	source := source_entity.(*entities.Planet)
	target_key := fmt.Sprintf("planet.%d_%d", mission.Target[0], mission.Target[1])
	target_entity, err := db_manager.GetEntity(target_key)
	target := target_entity.(*entities.Planet)

	speed := mission.GetSpeed()
	time.Sleep(CalculateArrivalTime(source.GetCoords(), target.GetCoords(), speed))

	// Fetch the end_planet again in order to know what has changed
	target_entity, err = db_manager.GetEntity(target_key)
	target = target_entity.(*entities.Planet)

	result := entities.EndMission(target, mission)
	key, serialized_planet, err := result.Serialize()
	if err == nil {
		db_manager.SetEntity(result)
		sessions.Broadcast([]byte(fmt.Sprintf("{\"Command\": \"state_change\", \"Planets\": {\"%s\": %s}}", key, serialized_planet)))
	}
	db_manager.DeleteEntity(mission.GetKey())
}
