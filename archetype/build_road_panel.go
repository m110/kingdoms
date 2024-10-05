package archetype

import (
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type roadPanelButton struct {
	Position math.Vec2
	Text     string
}

const roadPanelButtonPadding = 64

var roadPanelButtons = map[component.Direction]roadPanelButton{
	component.DirectionNorth: {Position: math.Vec2{X: 0, Y: -roadPanelButtonPadding}, Text: "^"},
	component.DirectionEast:  {Position: math.Vec2{X: roadPanelButtonPadding, Y: 0}, Text: ">"},
	component.DirectionSouth: {Position: math.Vec2{X: 0, Y: roadPanelButtonPadding}, Text: "v"},
	component.DirectionWest:  {Position: math.Vec2{X: -roadPanelButtonPadding, Y: 0}, Text: "<"},
}

// Experiment: tag defined here instead of components
var buildRoadPanel = donburi.NewTag()

func ShowBuildRoadPanel(w donburi.World) {
	game := component.MustFindGame(w)

	frame := NewFrameWithBorder(w, math.Vec2{X: 850, Y: float64(game.Settings.ScreenHeight) - 350}, 200, 200)
	component.Layer.Get(frame).Layer = component.SpriteUILayerUI
	frame.AddComponent(component.UIPanel)
	component.UIPanel.SetValue(frame, component.UIPanelData{
		KeepSelection: true,
	})
	frame.AddComponent(buildRoadPanel)

	buttonsContainer := New(w).
		WithParent(frame).
		WithPosition(math.Vec2{X: 72, Y: 72}).
		Entry()

	q := query.NewQuery(
		filter.And(
			filter.Contains(component.RoadPlaceholder),
			filter.Not(filter.Contains(component.Destroyed)),
		),
	)

	q.Each(w, func(entry *donburi.Entry) {
		direction := component.RoadPlaceholder.Get(entry).Direction
		buttonData := roadPanelButtons[direction]

		tile := engine.MustGetParent(entry)

		button := NewButtonWithWidth(w, buttonData.Text, buttonData.Position, true, func(w donburi.World, e *donburi.Entry) {
			events.NewRoadRequestedEvent.Publish(w, events.NewRoadRequested{
				Tile: tile,
			})
		}, 16)

		scale := 1.3
		component.Button.Get(button).KeepSelection = true
		transform.GetTransform(button).LocalScale = math.Vec2{X: scale, Y: scale}

		text := engine.MustFindChildWithComponent(button, component.Text)
		textPos := transform.GetTransform(text).LocalPosition
		transform.GetTransform(text).LocalPosition = textPos.MulScalar(scale)

		transform.AppendChild(buttonsContainer, button, false)
	})
}

func HideBuildRoadPanel(w donburi.World) {
	panel, ok := engine.FindWithComponent(w, buildRoadPanel)
	if !ok {
		return
	}

	component.Destroy(panel)
}
