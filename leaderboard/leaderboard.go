// Package leaderboard
package leaderboard

import (
	"errors"
	"sort"
)

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

type Leaderboard struct {
	board  []*Player
	places map[string]int
}

func New() *Leaderboard {
	l := new(Leaderboard)
	l.places = make(map[string]int)
	l.board = make([]*Player, 0)
	return l
}

func (l *Leaderboard) Add(player *Player) {
	l.board = append(l.board, player)
	l.places[player.Username] = len(l.board) - 1
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
}

func (l *Leaderboard) Transfer(from_username, to_username string) {
	from, hasOwner := l.places[from_username]
	to := l.places[to_username]

	if hasOwner {
		l.board[from].Planets--
		l.moveDown(from_username)
	}

	l.board[to].Planets++
	l.moveUp(to_username)
}

func (l *Leaderboard) Page(page int64) ([]*Player, error) {
	from := (page - 1) * 10
	to := ((page - 1) * 10) + 10
	if len(l.board) < int(from) {
		return []*Player{}, errors.New("Not such page")
	}

	if len(l.board) < int(to) {
		to = int64(len(l.board))
	}
	return l.board[from:to], nil
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
