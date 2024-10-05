package events

import (
	"github.com/m110/kingdoms/domain"
	"github.com/yohamta/donburi/features/events"
)

type ResearchRequested struct {
	Technology domain.Technology
}

var ResearchRequestedEvent = events.NewEventType[ResearchRequested]()

type TechnologyResearched struct {
	Technology domain.Technology
}

var TechnologyResearchedEvent = events.NewEventType[TechnologyResearched]()
