package archetype

import (
	"fmt"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
)

func ShowResourceInfoPanelIfRelevant(w donburi.World, tileEntry *donburi.Entry) *donburi.Entry {
	game := component.MustFindGame(w)

	// TODO For sure merge with unit panel
	_, ok := transform.FindChildWithComponent(tileEntry, component.Settlement)
	if ok {
		return nil
	}

	_, ok = transform.FindChildWithComponent(tileEntry, component.Settlers)
	if ok {
		return nil
	}

	tile := component.Tile.Get(tileEntry)
	hasDeposit := tile.Deposit != nil

	height := 130
	if !hasDeposit {
		height = 50
	}

	frame := NewFrameWithBorder(w, math.Vec2{X: 200, Y: float64(game.Settings.ScreenHeight - height)}, 350, 130)
	component.Layer.Get(frame).Layer = component.SpriteUILayerUI
	frame.AddComponent(component.UIPanel)
	frame.AddComponent(component.Script)
	component.UIPanel.SetValue(frame, component.UIPanelData{
		KeepSelection: true,
	})
	frame.AddComponent(component.ResourceInfoPanel)

	biomeName := domain.Terrains[tile.Terrain].Name

	New(w).
		WithParent(frame).
		WithPosition(math.Vec2{X: 10, Y: 30}).
		WithText(component.TextData{
			Text: biomeName,
		})

	if hasDeposit {
		deposit := domain.Deposits[*tile.Deposit]

		New(w).
			WithParent(frame).
			WithPosition(math.Vec2{X: 10, Y: 60}).
			WithText(component.TextData{
				Text: deposit.Name,
			}).
			Entry()

		icon := NewResourceIcon(w, tile.Resource.Resource, math.Vec2{X: 10, Y: 75})
		transform.GetTransform(icon).LocalScale = math.Vec2{X: 3, Y: 3}
		transform.AppendChild(frame, icon, false)

		amountText := New(w).
			WithParent(frame).
			WithPosition(math.Vec2{X: 65, Y: 95}).
			WithText(component.TextData{}).
			Entry()

		updateFunc := func(e *donburi.Entry) {
			component.Text.Get(amountText).Text = fmt.Sprintf("Amount: %v", tile.Resource.Amount)
		}

		component.Script.Get(frame).Update = updateFunc

		updateFunc(frame)
	} else {
		component.Script.Get(frame).Update = func(e *donburi.Entry) {}
	}

	return frame
}

func UpdateResourceInfoPanel(w donburi.World) {
	panel, ok := donburi.NewQuery(filter.Contains(component.ResourceInfoPanel)).First(w)
	if !ok {
		return
	}

	component.Script.Get(panel).Update(panel)
}

func HideResourceInfoPanel(w donburi.World) {
	panel, ok := donburi.NewQuery(filter.Contains(component.ResourceInfoPanel)).First(w)
	if !ok {
		return
	}

	component.Destroy(panel)
}
