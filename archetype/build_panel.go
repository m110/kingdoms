package archetype

import (
	"github.com/m110/kingdoms/domain"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/events"
)

func ShowBuildPanelIfRelevant(w donburi.World, tileEntry *donburi.Entry) *donburi.Entry {
	tile := component.Tile.Get(tileEntry)

	neighborCanBuildRoad := false
	neighbors := tile.NeighborTiles.All()
	for _, neighbor := range neighbors {
		t := component.Tile.Get(neighbor.Entry)
		if t.CanBuildRoad && !t.HasRoad {
			neighborCanBuildRoad = true
		}
	}

	showBuildRoad := tile.CanBuildRoad || (tile.HasRoad && neighborCanBuildRoad)
	showBuildSettlement := tile.CanBuildBuildings

	if !showBuildRoad && !showBuildSettlement {
		return nil
	}

	game := component.MustFindGame(w)

	frame := NewFrameWithBorder(w, math.Vec2{X: 750, Y: float64(game.Settings.ScreenHeight) - 130}, 350, 130)
	component.Layer.Get(frame).Layer = component.SpriteUILayerUI
	frame.AddComponent(component.UIPanel)
	component.UIPanel.SetValue(frame, component.UIPanelData{
		KeepSelection: true,
	})
	frame.AddComponent(component.BuildPanel)

	if showBuildRoad {
		roadButton := NewToolButton(w, math.Vec2{X: 10, Y: 10}, true, func(w donburi.World, e *donburi.Entry) {
			events.NewRoadRequestedEvent.Publish(w, events.NewRoadRequested{
				Tile: tileEntry,
			})
		}, domain.ToolRoad)
		component.Button.Get(roadButton).KeepSelection = true
		transform.GetTransform(roadButton).LocalScale = math.Vec2{X: 2, Y: 2}
		transform.AppendChild(frame, roadButton, false)
	}

	if showBuildSettlement {
		settlementButton := NewToolButton(w, math.Vec2{X: 120, Y: 10}, true, func(w donburi.World, e *donburi.Entry) {
			events.NewSettlementRequestedEvent.Publish(w, events.NewSettlementRequested{
				Tile: tileEntry,
			})
		}, domain.ToolSettlement)
		component.Button.Get(settlementButton).KeepSelection = true
		transform.GetTransform(settlementButton).LocalScale = math.Vec2{X: 2, Y: 2}
		transform.AppendChild(frame, settlementButton, false)
	}

	return frame
}

func HideBuildPanel(w donburi.World) bool {
	panel, ok := query.NewQuery(filter.Contains(component.BuildPanel)).First(w)
	if !ok {
		return false
	}

	component.Destroy(panel)

	return true
}
