package entities

import "time"

// Area transfer points are the points between the areasets
type AreaTransferPoint struct {
	TravelTime     time.Duration
	Direction      int8
	CoordinateAxis rune
}

// Just a sorting interface of AreaTransferPoint
type AreaTransferPoints []*AreaTransferPoint

func (a AreaTransferPoints) Len() int {
	return len(a)
}

func (a AreaTransferPoints) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a AreaTransferPoints) Less(i, j int) bool {
	return a[i].TravelTime < a[j].TravelTime
}
