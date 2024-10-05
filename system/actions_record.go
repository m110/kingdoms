package system

import (
	"github.com/yohamta/donburi"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
	"github.com/m110/kingdoms/save"
)

type ActionsRecord struct {
	actionsRecorder *component.ActionsRecorderData
	game            *component.GameData
}

func NewActionsRecord() *ActionsRecord {
	return &ActionsRecord{}
}

func (a *ActionsRecord) Init(w donburi.World) {
	a.actionsRecorder = engine.MustFindComponent[component.ActionsRecorderData](w, component.ActionsRecorder)
	a.game = engine.MustFindComponent[component.GameData](w, component.Game)

	events.NewRoadRequestedEvent.Subscribe(w, func(w donburi.World, event events.NewRoadRequested) {
		tile := component.Tile.Get(event.Tile)

		a.record(
			save.ActionType_ActionNewRoadRequested,
			&save.ActionPayload{
				Payload: &save.ActionPayload_NewRoadRequested{
					NewRoadRequested: &save.NewRoadRequested{
						Position: mapPositionToSave(tile.Position),
					},
				},
			},
		)
	})

	events.NewSettlementRequestedEvent.Subscribe(w, func(w donburi.World, event events.NewSettlementRequested) {
		tile := component.Tile.Get(event.Tile)
		a.record(
			save.ActionType_ActionNewSettlementRequested,
			&save.ActionPayload{
				Payload: &save.ActionPayload_NewSettlementRequested{
					NewSettlementRequested: &save.NewSettlementRequested{
						Position: mapPositionToSave(tile.Position),
					},
				},
			},
		)
	})

	events.SettlementUpgradeRequestedEvent.Subscribe(w, func(w donburi.World, event events.SettlementUpgradeRequested) {
		sett := component.Settlement.Get(event.Settlement)
		a.record(
			save.ActionType_ActionSettlementUpgradeRequested,
			&save.ActionPayload{
				Payload: &save.ActionPayload_SettlementUpgradeRequested{
					SettlementUpgradeRequested: &save.SettlementUpgradeRequested{
						SettlementId: int64(sett.ID),
					},
				},
			},
		)
	})

	events.CameraMovedEvent.Subscribe(w, func(w donburi.World, event events.CameraMoved) {
		a.record(
			save.ActionType_ActionCameraMoved,
			&save.ActionPayload{
				Payload: &save.ActionPayload_CameraMoved{
					CameraMoved: &save.CameraMoved{
						Delta: mapVectorToSave(event.Delta),
					},
				},
			},
		)
	})

	events.ResearchRequestedEvent.Subscribe(w, func(w donburi.World, event events.ResearchRequested) {
		a.record(
			save.ActionType_ActionResearchRequested,
			&save.ActionPayload{
				Payload: &save.ActionPayload_ResearchRequested{
					ResearchRequested: &save.ResearchRequested{
						Technology: mapTechnologyToSave(event.Technology),
					},
				},
			},
		)
	})
}

func (a *ActionsRecord) Update(w donburi.World) {
	if a.game.Paused {
		return
	}

	a.actionsRecorder.CurrentTick++
}

func (a *ActionsRecord) record(actionType save.ActionType, payload *save.ActionPayload) {
	a.actionsRecorder.Record(actionType, payload)
}
