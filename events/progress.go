package events

import (
	"github.com/yohamta/donburi/features/events"

	"github.com/m110/kingdoms/component"
)

type (
	GameOver        struct{}
	GameWon         struct{}
	ProgressUpdated struct {
		Progress *component.ProgressData
	}
	SeasonChanged struct {
		Progress *component.ProgressData
	}
)

var (
	GameOverEvent        = events.NewEventType[GameOver]()
	GameWonEvent         = events.NewEventType[GameWon]()
	ProgressUpdatedEvent = events.NewEventType[ProgressUpdated]()
	SeasonChangedEvent   = events.NewEventType[SeasonChanged]()
)
