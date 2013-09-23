package server

import "testing"

func TestResponse(t *testing.T) {
	const BEST_PING = 150
	const WORST_PING = 1500
	const STEPS = 10

	top_left, bottom_right := calculateCanvasSize([]int{20, 50}, []int{800, 600}, 200)
	if top_left[0] != -20 || top_left[1] != 20 ||
		bottom_right[0] != 860 || bottom_right[1] != 680 {
		t.Error("scopeOfView([20 50], [800 600], 200) gives", top_left, bottom_right)
	}
}
