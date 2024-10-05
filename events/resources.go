package events

import (
	"github.com/yohamta/donburi/features/events"

	"github.com/m110/kingdoms/component"
)

type ResourcesDelta struct {
	Food     int
	Stone    int
	Wood     int
	Iron     int
	Gold     int
	Diamonds int
}

type ResourcesUpdated struct {
	Resources *component.ResourcesData
	Delta     ResourcesDelta
}

var ResourcesUpdatedEvent = events.NewEventType[ResourcesUpdated]()
