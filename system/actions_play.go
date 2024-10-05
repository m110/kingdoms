package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
	"github.com/m110/kingdoms/save"
)

type ActionsPlay struct {
	game          *component.GameData
	actionsPlayer *component.ActionsPlayerData
}

func NewActionsPlay() *ActionsPlay {
	return &ActionsPlay{}
}

func (a *ActionsPlay) Init(w donburi.World) {
	a.game = engine.MustFindComponent[component.GameData](w, component.Game)
	a.actionsPlayer = engine.MustFindComponent[component.ActionsPlayerData](w, component.ActionsPlayer)
}

func (a *ActionsPlay) Update(w donburi.World) {
	if a.game.Paused {
		return
	}

	if a.actionsPlayer.Paused {
		return
	}

	a.actionsPlayer.CurrentTick += 1 * a.actionsPlayer.Speed

	actions := a.actionsPlayer.GetActionsForCurrentTick()

	for _, action := range actions {
		a.replayAction(w, action)
	}
}

func (a *ActionsPlay) replayAction(w donburi.World, action component.Action) {
	board := engine.MustFindComponent[component.BoardData](w, component.Board)

	switch action.ActionType {
	case save.ActionType_ActionNewRoadRequested:
		payload := action.Payload.GetNewRoadRequested()
		tile := board.MustTileByPosition(MapPositionFromSave(payload.Position))
		events.NewRoadRequestedEvent.Publish(w, events.NewRoadRequested{
			Tile: tile,
		})

	case save.ActionType_ActionNewSettlementRequested:
		payload := action.Payload.GetNewSettlementRequested()
		tile := board.MustTileByPosition(MapPositionFromSave(payload.Position))
		events.NewSettlementRequestedEvent.Publish(w, events.NewSettlementRequested{
			Tile: tile,
		})
	case save.ActionType_ActionSettlementUpgradeRequested:
		payload := action.Payload.GetSettlementUpgradeRequested()

		var settlement *donburi.Entry
		query.NewQuery(filter.Contains(component.Settlement)).Each(w, func(entry *donburi.Entry) {
			if settlement != nil {
				return
			}

			sett := component.Settlement.Get(entry)
			if int64(sett.ID) == payload.SettlementId {
				settlement = entry
			}
		})

		if settlement == nil {
			panic("settlement not found")
		}

		events.SettlementUpgradeRequestedEvent.Publish(w, events.SettlementUpgradeRequested{
			Settlement: settlement,
		})

	case save.ActionType_ActionCameraMoved:
		payload := action.Payload.GetCameraMoved()

		events.CameraMovedEvent.Publish(w, events.CameraMoved{
			Delta: MapVectorFromSave(payload.Delta),
		})

	case save.ActionType_ActionResearchRequested:
		payload := action.Payload.GetResearchRequested()

		events.ResearchRequestedEvent.Publish(w, events.ResearchRequested{
			Technology: MapTechnologyFromSave(payload.Technology),
		})
	}
}
