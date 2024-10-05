package events

import (
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
)

type CameraMoved struct {
	Delta math.Vec2
}

var CameraMovedEvent = events.NewEventType[CameraMoved]()
