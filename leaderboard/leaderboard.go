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
	from := Places[from_username]
	to := Places[to_username]

	if from != -1 {
		l[from].Planets--
		l.moveDown(from)
	}

	l[to].Planets++
	l.moveUp(to)
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
