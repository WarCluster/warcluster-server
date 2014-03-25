package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"warcluster/entities"
	"warcluster/leaderboard"
)

func leaderboardPlayersHandler(w http.ResponseWriter, r *http.Request) {
	pageQuery, ok := r.URL.Query()["page"]
	if !ok {
		http.Error(w, "Bad Request", 400)
		return
	}

	page, intErr := strconv.ParseInt(pageQuery[0], 10, 0)
	if intErr != nil {
		http.Error(w, "Page Not Found", 404)
		return
	}

	boardPage, err := leaderBoard.Page(page)
	if err != nil {
		http.Error(w, "Page Not Found", 404)
		return
	}

	result, _ := json.Marshal(boardPage)
	fmt.Fprintf(w, string(result))
}

func leaderboardTeamsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%#v\n", r.URL.Query())
}

// Initialize the leaderboard
func InitLeaderboard(board *leaderboard.Leaderboard) {
	log.Println("Initializing the leaderboard...")
	allPlayers := make(map[string]*leaderboard.Player)
	planetEntities := entities.Find("planet.*")

	for _, entity := range planetEntities {
		planet, ok := entity.(*entities.Planet)
		if !planet.HasOwner() {
			continue
		}

		player, ok := allPlayers[planet.Owner]

		if !ok {
			player = &leaderboard.Player{
				Username: planet.Owner,
				Team:     planet.Color,
				Planets:  0,
			}
			allPlayers[planet.Owner] = player
			board.Add(player)
		}

		if planet.IsHome {
			player.HomePlanet = planet.Name
		}

		player.Planets++
	}
	board.Sort()
	leaderBoard = board
}
