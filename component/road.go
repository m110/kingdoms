package component

import "github.com/yohamta/donburi"

type Direction int

const (
	DirectionNorth Direction = iota
	DirectionEast
	DirectionSouth
	DirectionWest
)

type RoadPlaceholderData struct {
	Direction Direction
}

var RoadPlaceholder = donburi.NewComponentType[RoadPlaceholderData]()
