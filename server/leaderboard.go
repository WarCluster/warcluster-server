package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"warcluster/config"
	"warcluster/entities"
	"warcluster/leaderboard"
)

type searchResult struct {
	Username string
	Page     int
}

func leaderboardPlayersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

func leaderboardRacesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	races, err := json.Marshal(leaderBoard.Races())
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	fmt.Fprintf(w, string(races))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	username, ok := r.URL.Query()["player"]
	if !ok || len(username[0]) < 3 {
		http.Error(w, "Bad Request", 400)
		return
	}

	players, err := entities.GetList(fmt.Sprintf("player.%s*", username[0]))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	result := make([]searchResult, 0)

	for _, player := range players {
		username := player[7:len(player)]
		page := math.Ceil(float64(leaderBoard.Place(username)+1) / 10)
		result = append(result, searchResult{username, int(page)})
	}

	marshalledResult, _ := json.Marshal(result)
	fmt.Fprintf(w, string(marshalledResult))
}

// Initialize the leaderboard
func InitLeaderboard(board *leaderboard.Leaderboard, cfg config.Config) {
	log.Println("Initializing the leaderboard...")
	allPlayers := make(map[string]*leaderboard.Player)
	planetEntities := entities.Find("planet.*")

	for key, value := range cfg.Race {
		board.AddRace(
			key,
			leaderboard.Color{
				value.Red,
				value.Green,
				value.Blue,
			},
		)
	}

	for _, entity := range planetEntities {
		planet, ok := entity.(*entities.Planet)
		if !planet.HasOwner() {
			continue
		}

		player, ok := allPlayers[planet.Owner]

		if !ok {
			player = &leaderboard.Player{
				Username: planet.Owner,
				Race:     planet.Color,
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
	board.RecountRacesPlanets()
	leaderBoard = board
}
