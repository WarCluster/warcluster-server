package server

import (
	"time"
	"../vector"
)

func CalculateArrivalTime(start_point, end_point *vector.Vector, speed int) time.Duration {
	distance := end_point.Substitute(start_point)
	return time.Duration(time.Duration(distance.Length() / float64(speed)) * time.Second)
}


