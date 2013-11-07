package entities

import "testing"

func TestUsernameHash(t *testing.T) {
	hash := generateHash("Gopher")
	if len(hash) != 64 {
		t.Error("Wrong hash length")
	}
}

func TestExtractUsernameInitials(t *testing.T) {
	initials := [][]string{
		{"rob_pike", "ROB"},
		{"K_Vladimiroff", "KVL"},
		{"vitaliy_filipov", "VIT"},
		{"denytodorova", "DEN"},
		{"Mordevil", "MOR"},
		{"Lord_Voldemort7", "LOR"},
		{"@someecards", "SOM"},
		{"ShitGirlsSay", "SHI"},
		{"big_ben_clock", "BIG"},
		{"@Oatmeal", "OAT"},
		{"bppetrov", "BPP"},
		{"karadimov", "KAR"},
		{"iamdevloper", "IAM"},
		{"_rsc", "RSC"},
		{"golang", "GOL"},
		{"La_Chasc0na", "LAC"},
		{"if__fi", "IFF"},
	}

	for _, pair := range initials {
		initial := extractUsernameInitials(pair[0])
		if pair[1] != initial {
			t.Errorf("extractUsernameInitials(%q) = %q wants %q", pair[0], initial, pair[1])
		}
	}
}
