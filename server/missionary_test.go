package server

import (
	"testing"
	"time"
	"../vector"
)

func TestArrivalTime(t *testing.T) {
	start_point := vector.New(100, 200)
	end_point := vector.New(800, 150)
	speed := 5
	arrival_time := CalculateArrivalTime(start_point, end_point, speed)

	if arrival_time != time.Duration(140 * time.Second) {
		t.Error("Wrong arrival time:", arrival_time)
	}
}
