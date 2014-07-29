package entities

import (
	"fmt"
	"time"

	"github.com/Vladimiroff/vec2d"
)

type SpyReport struct {
	Player     string
	Name       string
	Owner      string
	Position   *vec2d.Vector
	ShipCount  int32
	ValidUntil int64
	CreatedAt  int64
}

// Database key.
func (s *SpyReport) Key() string {
	return fmt.Sprintf("spy_report.%s_%d", s.Player, s.CreatedAt)
}

// It has to be there in order to implement Entity
func (s *SpyReport) AreaSet() string {
	return ""
}

func (s *SpyReport) IsValid() bool {
	return s.ValidUntil > time.Now().Unix()
}

func CreateSpyReport(target *Planet, mission *Mission) *SpyReport {
	now := time.Now()
	report := &SpyReport{
		Player:     mission.Player,
		Name:       target.Name,
		Owner:      target.Owner,
		Position:   target.Position,
		ShipCount:  target.ShipCount,
		CreatedAt:  now.Unix(),
		ValidUntil: now.Add(Settings.SpyReportValidity * time.Second).Unix(),
	}
	Save(report)
	return report
}
