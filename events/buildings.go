package events

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type NewRoadRequested struct {
	Tile *donburi.Entry
}

var NewRoadRequestedEvent = events.NewEventType[NewRoadRequested]()

type NewSettlementRequested struct {
	Tile *donburi.Entry
}

var NewSettlementRequestedEvent = events.NewEventType[NewSettlementRequested]()

type SettlementUpgradeRequested struct {
	Settlement *donburi.Entry
}

var SettlementUpgradeRequestedEvent = events.NewEventType[SettlementUpgradeRequested]()

type TileResourcesHarvested struct {
	Tile  *donburi.Entry
	Delta ResourcesDelta
}

var TileResourcesHarvestedEvent = events.NewEventType[TileResourcesHarvested]()

type RoadSpawned struct {
	Road *donburi.Entry
	Tile *donburi.Entry
}

var RoadSpawnedEvent = events.NewEventType[RoadSpawned]()

type OriginSettlementSpawned struct {
	Settlement *donburi.Entry
	Tile       *donburi.Entry
}

var OriginSettlementSpawnedEvent = events.NewEventType[OriginSettlementSpawned]()

type SettlementSpawned struct {
	Settlement *donburi.Entry
	Tile       *donburi.Entry
}

var SettlementSpawnedEvent = events.NewEventType[SettlementSpawned]()

type SettlementUpgraded struct {
	Settlement *donburi.Entry
	Tile       *donburi.Entry
}

var SettlementUpgradedEvent = events.NewEventType[SettlementUpgraded]()
