package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type WanderData struct {
	Min math.Vec2
	Max math.Vec2
}

var Wander = donburi.NewComponentType[WanderData]()
