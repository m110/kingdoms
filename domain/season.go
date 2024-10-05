package domain

import (
	"github.com/m110/kingdoms/engine"
)

const (
	DaysPerMonth = 20
	Seasons      = 4
	DaysPerYear  = DaysPerMonth * Seasons
)

type Season int

const (
	SeasonSpring Season = iota
	SeasonSummer
	SeasonFall
	SeasonWinter
)

func (s Season) NextSeason() Season {
	return (s + 1) % 4
}

func (s Season) PreviousSeason() Season {
	return (s - 1 + 4) % 4
}

func (s Season) String() string {
	switch s {
	case SeasonSpring:
		return "Spring"
	case SeasonSummer:
		return "Summer"
	case SeasonFall:
		return "Fall"
	case SeasonWinter:
		return "Winter"
	default:
		return ""
	}
}

func RandomSeason() Season {
	return Season(engine.RandomIntRange(0, 3))
}

func SeasonByDay(day int) Season {
	return Season((day - 1) / DaysPerMonth)
}

func SeasonDay(day int) int {
	return day - DaysPerMonth*int(SeasonByDay(day))
}
