// Package leaderboard
package leaderboard

import (
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
type Leaderboard []*Player

var (
	Places map[string]int
	Board  Leaderboard
)

func New() Leaderboard {
	if Board != nil {
		return Board
	}

	Places = make(map[string]int)
	Board = make(Leaderboard, 0)
	return Board
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

func (l Leaderboard) Transfer(from_username, to_username string) {
	from, hasOwner := Places[from_username]
	to := Places[to_username]

	if hasOwner {
		l[from].Planets--
		l.moveDown(from_username)
	}

	l[to].Planets++
	l.moveUp(to_username)
}

func (l Leaderboard) move(username string, modificator int) bool {
	firstBlood := true
	isMoved := false
	i := Places[username]

	for isMoved || firstBlood {
		firstBlood = false
		if modificator < 0 && i == 0 || modificator > 0 && i == len(l)-1 {
			return false
		}

		if modificator < 0 {
			isMoved = l.Less(i, i+modificator)
		} else {
			isMoved = l.Less(i+modificator, i)
		}

		if isMoved {
			l.Swap(i, i+modificator)
			Places[l[i].Username], Places[l[i+modificator].Username] = i+modificator, i
		}
		i += modificator
	}
	return isMoved
}

func (l Leaderboard) moveUp(username string) bool {
	return l.move(username, -1)
}

func (l Leaderboard) moveDown(username string) bool {
	return l.move(username, 1)
}
