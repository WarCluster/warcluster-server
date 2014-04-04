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
	Team     struct {
		R float32
		G float32
		B float32
	}
	HomePlanet string
	Planets    uint32
}

type Team struct {
	Name    string
	Color   Color
	Players uint32
	Planets uint32
}

type Teams []*Team

type Leaderboard struct {
	places map[string]int
	board  []*Player
	teams  Teams
}

func New() *Leaderboard {
	l := new(Leaderboard)
	l.places = make(map[string]int)
	l.board = make([]*Player, 0)
	l.teams = make([]*Team, 0)
	return l
}

func (l *Leaderboard) Add(player *Player) {
	l.board = append(l.board, player)
	l.places[player.Username] = len(l.board) - 1

	for _, team := range l.teams {
		if player.Team == team.Color {
			team.Players++
			team.Planets += player.Planets
			return
		}
	}

	l.teams = append(l.teams, &Team{
		Name:    "Red Panda",
		Color:   player.Team,
		Players: 1,
		Planets: player.Planets,
	})
}

func (l *Leaderboard) FindTeam(color Color) *Team {
	for _, team := range l.teams {
		if team.Color == color {
			return team
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
}

func (l *Leaderboard) RecountTeamsPlanets() {
	for _, team := range l.teams {
		team.Planets = 0
	}

	for _, player := range l.board {
		team := l.FindTeam(player.Team)
		if team != nil {
			team.Planets += player.Planets
		}
	}
}

func (l *Leaderboard) Transfer(from_username, to_username string) {
	from, hasOwner := l.places[from_username]
	to := l.places[to_username]

	if hasOwner {
		l.board[from].Planets--
		l.moveDown(from_username)
		team := l.FindTeam(l.board[from].Team)
		if team != nil {
			team.Planets++
		}
	}

	l.board[to].Planets++
	l.moveUp(to_username)
	team := l.FindTeam(l.board[to].Team)
	if team != nil {
		team.Planets++
	}
	l.teams.Sort()
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

func (l *Leaderboard) Teams() []*Team {
	return l.teams
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
			l.places[l.board[i].Username], l.places[l.board[i+modificator].Username] = i+modificator, i
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

func (t *Teams) Len() int {
	return len(*t)
}

func (t *Teams) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

func (t *Teams) Less(i, j int) bool {
	return (*t)[i].Planets > (*t)[j].Planets
}

func (t *Teams) Sort() {
	sort.Sort(t)
}
