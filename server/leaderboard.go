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
	defer func() {
		if panicked := recover(); panicked != nil {
			fmt.Fprintf(w, "HTTP/1.1 400 Bad Request\r\n\r\n")
			return
		}
	}()

	pageQuery := r.URL.Query()["page"]
	page, _ := strconv.ParseInt(pageQuery[0], 10, 0)
	result, _ := json.Marshal(leaderBoard.Page(page))
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
