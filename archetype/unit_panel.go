package archetype

import (
	"fmt"

	"github.com/m110/kingdoms/assets"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/events"
)

// TODO Probably a good idea to deduplicate with ResourceInfoPanel
func ShowUnitInfoPanelIfRelevant(w donburi.World, tileEntry *donburi.Entry) *donburi.Entry {
	// TODO support more buildings
	// TODO should this logic be here?
	settlementEntry, hasSettlement := transform.FindChildWithComponent(tileEntry, component.Settlement)
	_, hasSettlers := transform.FindChildWithComponent(tileEntry, component.Settlers)

	if !hasSettlement && !hasSettlers {
		return nil
	}

	game := component.MustFindGame(w)

	frame := NewFrameWithBorder(w, math.Vec2{X: 200, Y: float64(game.Settings.ScreenHeight) - 130}, 350, 130)
	component.Layer.Get(frame).Layer = component.SpriteUILayerUI
	frame.AddComponent(component.UIPanel)
	frame.AddComponent(component.Script)
	component.UIPanel.SetValue(frame, component.UIPanelData{
		KeepSelection: true,
	})
	frame.AddComponent(component.UnitInfoPanel)

	title := New(w).
		WithParent(frame).
		WithPosition(math.Vec2{X: 10, Y: 30}).
		WithText(component.TextData{}).
		Entry()

	var icon *donburi.Entry
	iconPos := math.Vec2{X: 10, Y: 90}

	var newIcon func() *donburi.Entry

	if hasSettlement {
		settlement := component.Settlement.Get(settlementEntry)
		newIcon = func() *donburi.Entry {
			i := NewSettlementIcon(w, settlement.Level, iconPos)
			transform.AppendChild(frame, i, false)
			return i
		}
		icon = newIcon()
	} else {
		component.Text.Get(title).Text = "Settlers"
		icon = New(w).
			WithParent(frame).
			WithPosition(iconPos).
			WithScale(math.Vec2{X: 2, Y: 2}).
			WithSprite(component.SpriteData{
				Image: assets.Settlers,
			}).
			Entry()
	}

	updateFunc := func(e *donburi.Entry) {
		if !hasSettlement {
			return
		}
		component.Destroy(icon)
		icon = newIcon()

		settlement := component.Settlement.Get(settlementEntry)
		component.Text.Get(title).Text = fmt.Sprintf("%v (%v)", SettlementNameFromLevel(settlement.Level), settlement.ID)

		btn, ok := transform.FindChildWithComponent(frame, component.Button)
		if ok {
			component.Destroy(btn)
		}
		if settlement.Level < domain.MaxSettlementLevelWithCurrency {
			upgradeButton := NewToolButton(
				w,
				math.Vec2{X: 250, Y: 10},
				true,
				func(w donburi.World, e *donburi.Entry) {
					events.SettlementUpgradeRequestedEvent.Publish(w, events.SettlementUpgradeRequested{
						Settlement: settlementEntry,
					})
				},
				domain.ToolUpgrade,
			)
			component.Button.Get(upgradeButton).KeepSelection = true
			transform.GetTransform(upgradeButton).LocalScale = math.Vec2{X: 2, Y: 2}
			transform.AppendChild(frame, upgradeButton, false)
		}
	}

	component.Script.Get(frame).Update = updateFunc

	updateFunc(frame)

	return frame
}

func UpdateUnitInfoPanel(w donburi.World) {
	panel, ok := donburi.NewQuery(filter.Contains(component.UnitInfoPanel)).First(w)
	if !ok {
		return
	}

	component.Script.Get(panel).Update(panel)
}

func HideUnitInfoPanel(w donburi.World) {
	panel, ok := donburi.NewQuery(filter.Contains(component.UnitInfoPanel)).First(w)
	if !ok {
		return
	}

	component.Destroy(panel)
}
