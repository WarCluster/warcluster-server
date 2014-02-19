// Package leaderboard
package leaderboard

import (
	"log"
	"sort"

	"warcluster/entities"
)

type Player struct {
	Username   string
	Team       *entities.Color
	HomePlanet string
	Planets    uint32
}
type Leaderboard []*Player

var (
	board Leaderboard
)

func New() Leaderboard {
	if board != nil {
		return board
	}

	board = make(Leaderboard, 0)
	return board
}

func (l Leaderboard) Init() {
	log.Println("Initializing the leaderboard...")
	allPlayers := make(map[string]*Player)
	planetEntities := entities.Find("planet.*")

	for _, entity := range planetEntities {
		planet, ok := entity.(*entities.Planet)
		if !planet.HasOwner() {
			continue
		}

		player, ok := allPlayers[planet.Owner]

		if !ok {
			player = &Player{
				Username: planet.Owner,
				Team:     &planet.Color,
				Planets:  0,
			}
			allPlayers[planet.Owner] = player
			l = append(l, player)
		}

		if planet.IsHome {
			player.HomePlanet = planet.Name
		}

		player.Planets++
	}
	l.Sort()
}

func (l Leaderboard) Len() int {
	return len(l)
}

func (l Leaderboard) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Leaderboard) Less(i, j int) bool {
	return l[i].Planets > l[j].Planets
}

func (l Leaderboard) Sort() {
	sort.Sort(l)
}

func (l Leaderboard) Transfer(from, to int) {
	if from != -1 {
		l[from].Planets--
	}

	l[to].Planets++
}

func (l Leaderboard) move(index, modificator int) bool {
	firstBlood := true
	isMoved := false

	for isMoved || firstBlood {
		firstBlood = false
		if modificator < 0 && index == 0 || modificator > 0 && index == len(l)-1 {
			return false
		}

		if modificator < 0 {
			isMoved = l.Less(index, index+modificator)
		} else {
			isMoved = l.Less(index+modificator, index)
		}

		if isMoved {
			l.Swap(index, index+modificator)
		}
		index += modificator
	}
	return isMoved
}

func (l Leaderboard) moveUp(i int) bool {
	return l.move(i, -1)
}

func (l Leaderboard) moveDown(i int) bool {
	return l.move(i, 1)
}
