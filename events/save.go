package events

import "github.com/yohamta/donburi/features/events"

type (
	QuitGameRequested struct{}
)

var (
	QuitGameRequestedEvent = events.NewEventType[QuitGameRequested]()
)
