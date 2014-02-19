package leaderboard

import (
	"testing"
)

func initLeaderboard() Leaderboard {
	board = Leaderboard{
		{Username: "first", Planets: 10},
		{Username: "second", Planets: 8},
		{Username: "third", Planets: 7},
		{Username: "forth", Planets: 6},
		{Username: "fifth", Planets: 5},
	}
	return board
}

func TestSimpleTransfer(t *testing.T) {
	board = initLeaderboard()
	board.Transfer(1, 0)
	if board[0].Planets != 11 && board[1].Planets != 7 {
		t.Errorf(
			"first has %d and second has %d planets, instead of 11 and 7",
			board[0].Planets,
			board[1].Planets,
		)
	}

	if board[0].Username != "first" {
		t.Errorf("0 is %s instead of first", board[0].Username)
	}

	if board[1].Username != "second" {
		t.Errorf("1 is %s instead of first", board[1].Username)
	}
}

func TestMovingUp(t *testing.T) {
	board = initLeaderboard()
	board[4].Planets = 9

	board.moveUp(4)
	if board[1].Username != "fifth" {
		t.Errorf("fifth is not in the second place, %s is there instead", board[1].Username)
	}

	if board[0].Username != "first" {
		t.Errorf("first is not in the first place, %s is there instead", board[1].Username)
	}

	board.moveUp(0)
}

func TestMovingDown(t *testing.T) {
	board = initLeaderboard()
	board[1].Planets = 2

	board.moveDown(1)
	if board[4].Username != "second" {
		t.Errorf("second is not last, %s is there instead", board[4].Username)
	}
	board.moveUp(len(board) - 1)
}
