package leaderboard

import (
	"testing"
)

func initLeaderboard() *Leaderboard {
	l := New()
	l.board = []*Player{
		{Username: "0", Planets: 10, Team: Color{13, 11, 92}},
		{Username: "1", Planets: 8, Team: Color{16, 5, 90}},
		{Username: "2", Planets: 7, Team: Color{16, 5, 90}},
		{Username: "3", Planets: 6, Team: Color{13, 11, 92}},
		{Username: "4", Planets: 5, Team: Color{13, 11, 92}},
		{Username: "5", Planets: 0, Team: Color{20, 6, 90}},
		{Username: "6", Planets: 0, Team: Color{13, 11, 92}},
		{Username: "7", Planets: 0, Team: Color{20, 6, 90}},
		{Username: "8", Planets: 0, Team: Color{20, 6, 90}},
		{Username: "9", Planets: 0, Team: Color{20, 6, 90}},
		{Username: "10", Planets: 0, Team: Color{20, 6, 90}},
		{Username: "11", Planets: 0, Team: Color{13, 11, 92}},
		{Username: "12", Planets: 0, Team: Color{16, 5, 90}},
		{Username: "13", Planets: 0, Team: Color{13, 11, 92}},
		{Username: "14", Planets: 0, Team: Color{20, 6, 90}},
		{Username: "15", Planets: 0, Team: Color{16, 5, 90}},
		{Username: "16", Planets: 0, Team: Color{16, 5, 90}},
		{Username: "17", Planets: 0, Team: Color{16, 5, 90}},
		{Username: "18", Planets: 0, Team: Color{20, 6, 90}},
	}

	l.places = map[string]int{
		"0":  0,
		"1":  1,
		"2":  2,
		"3":  3,
		"4":  4,
		"5":  5,
		"6":  6,
		"7":  7,
		"8":  8,
		"9":  9,
		"10": 10,
		"11": 11,
		"12": 12,
		"13": 13,
		"14": 14,
		"15": 15,
		"16": 16,
		"17": 17,
		"18": 18,
	}

	l.teams = Teams{
		{Color: Color{13, 11, 92}},
		{Color: Color{20, 6, 90}},
		{Color: Color{16, 5, 90}},
	}

	l.RecountTeamsPlanets()

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

func TestPlace(t *testing.T) {
	l := initLeaderboard()
	results := []struct {
		in  string
		out int
	}{
		{"0", 0},
		{"4", 4},
		{"2", 2},
		{"42", 0},
	}

	for _, result := range results {
		place := l.Place(result.in)
		if place != result.out {
			t.Errorf("Player %#v is in place %d, expected %d", result.in, place, result.out)
		}
	}
}

func TestPage(t *testing.T) {
	l := initLeaderboard()
	results := []struct {
		in  int64
		out int
		err bool
	}{
		{-1, 0, true},
		{0, 0, true},
		{1, 10, false},
		{2, 9, false},
		{3, 0, true},
	}

	for _, result := range results {
		page, err := l.Page(result.in)
		if len(page) != result.out || result.err != (err != nil) {
			t.Errorf(
				"Page %d returned %d (err: %t) players, expected %d (err: %t)",
				result.in, len(page), err != nil, result.out, result.err)
		}
	}
}

func TestAdd(t *testing.T) {
	l := initLeaderboard()
	boardLengthBefore := l.Len()
	placesLengthBefore := len(l.places)
	teamsLengthBefore := len(l.teams)
	l.Add(&Player{Username: "panda", Planets: 42, Team: Color{13, 11, 92}})
	l.FindTeam(Color{21, 6, 90})
	l.Add(&Player{Username: "gophie", Planets: 42, Team: Color{21, 6, 90}})
	if boardLengthBefore+2 != l.Len() {
		t.Error("Board size did not changed after adding a player")
	}

	if placesLengthBefore+2 != len(l.places) {
		t.Error("Places map size did not changed after adding a player")
	}

	if teamsLengthBefore+1 != len(l.Teams()) {
		t.Error("Places map size did not changed after adding a player")
	}
}

func TestSort(t *testing.T) {
	l := initLeaderboard()
	l.places["0"] = 20
	l.board[1].Planets = 128
	l.RecountTeamsPlanets()
	l.Sort()

	if l.board[0].Username != "1" {
		t.Error("Leaderboard.Sort() did not sorted the board")
	}

	if l.places["0"] != 1 {
		t.Error("Leaderboard.Sort() did not changed the places")
	}

	if l.teams[0] != l.FindTeam(Color{16, 5, 90}) {
		t.Error("Leaderboard.Sort() did not sort teams")
	}
}

func TestChangingPlacesAndPlanets(t *testing.T) {
	l := initLeaderboard()
	l.Transfer("3", "4")
	l.Transfer("0", "4")

	if l.places["4"] > l.places["3"] {
		t.Errorf(
			"Player 4 is on place %d, Player 3 - %d",
			l.places["4"],
			l.places["3"],
		)
	}

	planets := l.board[l.places["4"]].Planets
	if planets != 7 {
		t.Errorf("Player 4 is has %d planets, expected 7", planets)
	}
}

func TestTakePlanetsWithoutOwner(t *testing.T) {
	l := initLeaderboard()
	l.Transfer("", "4")
	l.Transfer("", "4")

	if l.places["4"] > l.places["3"] {
		t.Errorf(
			"Player 4 is on place %d, Player 3 - %d",
			l.places["4"],
			l.places["3"],
		)
	}
}

func TestTeamPlanetsTransfer(t *testing.T) {
	d := New()
	d.board = []*Player{
		{Username: "0", Planets: 10, Team: Color{13, 11, 92}},
		{Username: "1", Planets: 10, Team: Color{16, 5, 90}},
		{Username: "2", Planets: 9, Team: Color{16, 5, 90}},
		{Username: "3", Planets: 8, Team: Color{13, 11, 92}},
		{Username: "4", Planets: 7, Team: Color{6, 9, 90}},
	}

	d.places = map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
	}

	d.teams = Teams{
		{Color: Color{13, 11, 92}},
		{Color: Color{6, 9, 90}},
		{Color: Color{16, 5, 90}},
	}

	d.RecountTeamsPlanets()

	d.Transfer("", "3")
	d.Transfer("", "3")
	d.Transfer("", "2")

	planets := d.FindTeam(Color{13, 11, 92}).Planets
	if planets != 20 {
		t.Errorf("This team has %d", planets)
	}
}
