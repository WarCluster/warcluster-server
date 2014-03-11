package leaderboard

import (
	"testing"
)

func initLeaderboard() *Leaderboard {
	l := new(Leaderboard)
	l.board = []*Player{
		{Username: "0", Planets: 10},
		{Username: "1", Planets: 8},
		{Username: "2", Planets: 7},
		{Username: "3", Planets: 6},
		{Username: "4", Planets: 5},
	}
	l.places = map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
	}
	return l
}

func TestSimpleTransfer(t *testing.T) {
	l := initLeaderboard()
	l.Transfer("1", "0")
	if l.board[0].Planets != 11 && l.board[1].Planets != 7 {
		t.Errorf(
			"0 has %d and 1 has %d planets, instead of 11 and 7",
			l.board[0].Planets,
			l.board[1].Planets,
		)
	}

	if l.board[0].Username != "0" {
		t.Errorf("0 is %s instead of 0", l.board[0].Username)
	}

	if l.board[1].Username != "1" {
		t.Errorf("1 is %s instead of 0", l.board[1].Username)
	}
}

func TestMovingUp(t *testing.T) {
	l := initLeaderboard()
	l.board[4].Planets = 9

	l.moveUp("4")
	if l.board[1].Username != "4" {
		t.Errorf("4 is not in the 1 place, %s is there instead", l.board[1].Username)
	}

	if l.board[0].Username != "0" {
		t.Errorf("0 is not in the 0 place, %s is there instead", l.board[1].Username)
	}

	l.moveUp("0")
}

func TestMovingDown(t *testing.T) {
	l := initLeaderboard()
	l.board[1].Planets = 2

	l.moveDown("1")
	if l.board[4].Username != "1" {
		t.Errorf("1 is not last, %s is there instead", l.board[4].Username)
	}
	l.moveUp(string(len(l.board) - 1))
}
