package server

import "testing"

func TestResponse(t *testing.T) {
	const BEST_PING = 150
	const WORST_PING = 1500
	const STEPS = 10

	topLeft, bottomRight := calculateCanvasSize([]int{20, 50}, []int{800, 600}, 200)
	if topLeft[0] != -20 || topLeft[1] != 20 ||
		bottomRight[0] != 860 || bottomRight[1] != 680 {
		t.Error("scopeOfView([20 50], [800 600], 200) gives", topLeft, bottomRight)
	}
}
