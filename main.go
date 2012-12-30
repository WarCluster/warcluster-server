package main

import (
	"./db_manager"
	"./entities"
	net "./server"
	"log"
	"os"
	"os/signal"
)

var server *net.Server = &net.Server{}

func main() {
	go final()
	defer func() {
		final()
	}()

	username := "gophie"
	sun_position := []int{500, 300}
	hash := entities.GenerateHash(username)
	_, home_planet := entities.GeneratePlanets(hash, sun_position)
	player := entities.CreatePlayer(username, hash, home_planet)

	log.Println("Created player:", player)
	db_manager.SetEntity(player)
	if new_player := db_manager.GetEntity(player.GetKey()); new_player != nil {
		log.Println("Fetched player from the db:", new_player)
	}
	log.Println("------------------------------")
	server.Start("localhost", 7000)
	return
}

func final() {
	sigtermchan := make(chan os.Signal, 1)
	signal.Notify(sigtermchan, os.Interrupt)
	<- sigtermchan

	server.Stop()
	db_manager.Finalize()
}
