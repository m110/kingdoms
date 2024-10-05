package events

import (
	"github.com/m110/kingdoms/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type TileUnselected struct {
	Tile *donburi.Entry
}

var TileUnselectedEvent = events.NewEventType[TileUnselected]()

type TileSelected struct {
	Tile *donburi.Entry
}

var TileSelectedEvent = events.NewEventType[TileSelected]()

type ButtonClicked struct {
	Button *component.ButtonData
}

var ButtonClickedEvent = events.NewEventType[ButtonClicked]()
