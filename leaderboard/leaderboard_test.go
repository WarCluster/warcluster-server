package leaderboard

import (
	"testing"
)

func initLeaderboard() Leaderboard {
	board := Leaderboard{
		{Username: "0", Planets: 10},
		{Username: "1", Planets: 8},
		{Username: "2", Planets: 7},
		{Username: "3", Planets: 6},
		{Username: "4", Planets: 5},
	}
	Places = map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
	}
	return board
}

func TestSimpleTransfer(t *testing.T) {
	board := initLeaderboard()
	board.Transfer("1", "0")
	if board[0].Planets != 11 && board[1].Planets != 7 {
		t.Errorf(
			"0 has %d and 1 has %d planets, instead of 11 and 7",
			board[0].Planets,
			board[1].Planets,
		)
	}

	if board[0].Username != "0" {
		t.Errorf("0 is %s instead of 0", board[0].Username)
	}

	if board[1].Username != "1" {
		t.Errorf("1 is %s instead of 0", board[1].Username)
	}
}

func TestMovingUp(t *testing.T) {
	board := initLeaderboard()
	board[4].Planets = 9

	board.moveUp("4")
	if board[1].Username != "4" {
		t.Errorf("4 is not in the 1 place, %s is there instead", board[1].Username)
	}

	if board[0].Username != "0" {
		t.Errorf("0 is not in the 0 place, %s is there instead", board[1].Username)
	}

	board.moveUp("0")
}

func TestMovingDown(t *testing.T) {
	board := initLeaderboard()
	board[1].Planets = 2

	board.moveDown("1")
	if board[4].Username != "1" {
		t.Errorf("1 is not last, %s is there instead", board[4].Username)
	}
	board.moveUp(string(len(board) - 1))
}
