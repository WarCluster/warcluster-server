package server

import (
	"../db_manager"
	e "../entities"
	"errors"
	"fmt"
)

const BEST_PING = 150
const WORST_PING = 1500
const STEPS = 10

func scopeOfView(position []int, resolution []int, lag int) ([]int, []int) {
	step := int(WORST_PING - BEST_PING/STEPS)
	multiply := 1.1 + float32((lag-BEST_PING)/step)*0.1
	end_resolution := []int{
		int(float32(resolution[0]) * multiply),
		int(float32(resolution[1]) * multiply),
	}

	top_left := []int{
		position[0] - int((end_resolution[0]-resolution[0])/2),
		position[1] - int((end_resolution[1]-resolution[1])/2),
	}

	bottom_right := []int{
		position[0] + resolution[0] + int((end_resolution[0]-resolution[0])/2),
		position[1] + resolution[1] + int((end_resolution[1]-resolution[1])/2),
	}
	return top_left, bottom_right
}

func actionParser(username, start_planet_key, end_planet_key string, fleet int) error {
	var err error = nil

	defer func() error {
		if panicked := recover(); panicked != nil {
			err = errors.New("Invalid action!")
		}
		return nil
	}()

	player_entity, err := db_manager.GetEntity(fmt.Sprint("player.", username))
	if err != nil {
		return errors.New("Player does not exist")
	}
	player := player_entity.(e.Player)

	start_planet, err := db_manager.GetEntity(fmt.Sprint("planet.", start_planet_key))
	if err != nil {
		return errors.New("Start planet does not exist")
	}

	end_planet, err := db_manager.GetEntity(fmt.Sprint("planet.", end_planet_key))
	if err != nil {
		return errors.New("End planet does not exist")
	}

	if start_planet.(e.Planet).Owner != username {
		err = errors.New("This is not your home!")
	}

	player.StartMission(start_planet.(e.Planet), end_planet.(e.Planet), fleet)
	// TODO: Write this mission in the DB
	// Insert this mission in the sorted list with missions

	return err
}
