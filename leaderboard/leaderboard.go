// Package leaderboard
package leaderboard

import (
	"errors"
	"sort"
)

type Color struct {
	R float32
	G float32
	B float32
}

type Player struct {
	Username string
	Race     struct {
		R float32
		G float32
		B float32
	}
	HomePlanet string
	Planets    uint32
}

type Race struct {
	Name    string
	Color   Color
	Players uint32
	Planets uint32
}

type Races []*Race

type Leaderboard struct {
	places  map[string]int
	board   []*Player
	races   Races
	Channel chan [2]string
}

func New() *Leaderboard {
	l := new(Leaderboard)
	l.places = make(map[string]int)
	l.board = make([]*Player, 0)
	l.races = make([]*Race, 0)
	l.Channel = make(chan [2]string)

	go func(l *Leaderboard) {
		var transfer [2]string
		for {
			transfer = <-l.Channel
			l.Transfer(transfer[0], transfer[1])
		}
	}(l)

	return l
}

func (l *Leaderboard) Add(player *Player) {
	l.board = append(l.board, player)
	l.places[player.Username] = len(l.board) - 1

	for _, race := range l.races {
		if player.Race == race.Color {
			race.Players++
			race.Planets += player.Planets
			return
		}
	}
}

func (l *Leaderboard) FindRace(color Color) *Race {
	for _, race := range l.races {
		if race.Color == color {
			return race
		}
	}
	return nil
}

func (l *Leaderboard) Len() int {
	return len(l.board)
}

func (l *Leaderboard) Swap(i, j int) {
	l.board[i], l.board[j] = l.board[j], l.board[i]
}

func (l *Leaderboard) Less(i, j int) bool {
	return l.board[i].Planets > l.board[j].Planets
}

func (l *Leaderboard) Sort() {
	sort.Sort(l)
	for index, player := range l.board {
		l.places[player.Username] = index
	}
	l.races.Sort()
}

func (l *Leaderboard) RecountRacesPlanets() {
	for _, race := range l.races {
		race.Planets = 0
	}

	for _, player := range l.board {
		race := l.FindRace(player.Race)
		if race != nil {
			race.Planets += player.Planets
		}
	}
}

func (l *Leaderboard) Transfer(from_username, to_username string) {
	from, hasOwner := l.places[from_username]
	to := l.places[to_username]

	if hasOwner {
		l.board[from].Planets--
	}
	l.board[to].Planets++

	if hasOwner {
		race := l.FindRace(l.board[from].Race)
		if race != nil {
			race.Planets--
		}
		l.moveDown(from_username)
	}

	race := l.FindRace(l.board[to].Race)
	if race != nil {
		race.Planets++
	}
	l.moveUp(to_username)
	l.races.Sort()
}

func (l *Leaderboard) Page(page int64) ([]*Player, error) {
	if page <= 0 {
		return []*Player{}, errors.New("No such page")
	}
	from := (page - 1) * 10
	to := ((page - 1) * 10) + 10
	if len(l.board) < int(from) {
		return []*Player{}, errors.New("No such page")
	}

	if len(l.board) < int(to) {
		to = int64(len(l.board))
	}
	return l.board[from:to], nil
}

func (l *Leaderboard) Place(username string) int {
	place, ok := l.places[username]
	if !ok {
		return 0
	}

	return place
}

func (l *Leaderboard) Races() []*Race {
	return l.races

}

func (l *Leaderboard) move(username string, modificator int) bool {
	firstBlood := true
	isMoved := false
	i := l.places[username]

	for isMoved || firstBlood {
		firstBlood = false
		if modificator < 0 && i == 0 || modificator > 0 && i == len(l.board)-1 {
			return false
		}

		if modificator < 0 {
			isMoved = l.Less(i, i+modificator)
		} else {
			isMoved = l.Less(i+modificator, i)
		}

		if isMoved {
			l.Swap(i, i+modificator)
			l.places[l.board[i].Username] = i
			l.places[l.board[i+modificator].Username] = i + modificator
		}
		i += modificator
	}
	return isMoved
}

func (l *Leaderboard) moveUp(username string) bool {
	return l.move(username, -1)
}

func (l *Leaderboard) moveDown(username string) bool {
	return l.move(username, 1)
}

func (l *Leaderboard) AddRace(name string, color Color) {
	l.races = append(l.races, &Race{
		Name:    name,
		Color:   color,
		Players: 0,
		Planets: 0,
	})
}

func (r *Races) Len() int {
	return len(*r)
}

func (r *Races) Swap(i, j int) {
	(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
}

func (r *Races) Less(i, j int) bool {
	return (*r)[i].Planets > (*r)[j].Planets
}

func (r *Races) Sort() {
	sort.Sort(r)
}
