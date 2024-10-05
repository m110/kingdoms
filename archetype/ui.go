package archetype

import (
	"fmt"
	"github.com/m110/kingdoms/domain"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/events"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"golang.org/x/image/font"
)

func NewToolButton(
	w donburi.World,
	pos math.Vec2,
	isActionButton bool,
	onClick func(w donburi.World, e *donburi.Entry),
	tool domain.Tool,
) *donburi.Entry {
	buttonImage := ebiten.NewImage(40, 40)
	buttonImage.Fill(color.RGBA{
		R: 193,
		G: 130,
		B: 82,
		A: 255,
	})

	// TODO should not be hardcoded here
	var costIcon *ebiten.Image
	var cost int
	var icon *ebiten.Image
	var key string

	// TODO cost should not be hardcoded
	switch tool {
	case domain.ToolSelect:
		icon = assets.Crosshair
		key = "Q"
	case domain.ToolRoad:
		icon = assets.RoadNE
		costIcon = assets.IconStone
		cost = 1
		key = "W"
	case domain.ToolSettlement:
		icon = assets.SettlementLevel1
		costIcon = assets.IconWood
		cost = 5
		key = "E"
	case domain.ToolUpgrade:
		icon = assets.IconUpgrade
		costIcon = assets.IconFood
		cost = 1
		key = "R"
	default:
		panic("unknown tool")
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(4, 4)
	buttonImage.DrawImage(icon, op)

	button := New(w).
		WithPosition(pos).
		WithLayer(component.SpriteUILayerUI).
		WithSprite(component.SpriteData{
			Image: buttonImage,
		}).
		With(component.Collider).
		With(component.Button).
		With(component.Active).
		Entry()

	component.Collider.SetValue(button, component.ColliderData{
		Width:  40,
		Height: 40,
		Layer:  component.CollisionLayerButtons,
	})
	component.Button.SetValue(button, component.ButtonData{
		OnClick: func(w donburi.World, e *donburi.Entry) {
			events.ButtonClickedEvent.Publish(w, events.ButtonClicked{
				Button: component.Button.Get(e),
			})
			onClick(w, e)
		},
		Sound: component.ButtonSoundNone,
	})
	setButtonActive(w, button, isActionButton)

	selectedIndicatorImage := ebiten.NewImage(44, 44)
	vector.StrokeRect(selectedIndicatorImage, 0, 0, 44, 44, 4, color.RGBA{
		R: 224,
		G: 186,
		B: 139,
		A: 255,
	}, true)

	New(w).
		WithParent(button).
		WithPosition(math.Vec2{X: -2, Y: -2}).
		WithSprite(component.SpriteData{
			Image:  selectedIndicatorImage,
			Hidden: true,
		}).
		Entry()

	if costIcon != nil {
		costIconEntry := New(w).
			WithParent(button).
			WithPosition(math.Vec2{X: 32, Y: 80}).
			WithSprite(component.SpriteData{
				Image: costIcon,
			}).
			Entry()

		costText := strconv.Itoa(cost)
		width := font.MeasureString(assets.SmallSquareFont, costText).Round()

		New(w).
			WithParent(costIconEntry).
			WithPosition(math.Vec2{X: -float64(width * 2), Y: 25}).
			WithText(component.TextData{
				Text: costText,
				Size: component.TextSizeM,
			})
	}

	// TODO Temporarily disabled: rethink if needed
	_ = key
	/*
		New(w).
			WithParent(button).
			WithPosition(math.Vec2{X: 17, Y: -4}).
			WithText(component.TextData{
				Text: key,
				Size: component.TextSizeS,
			})
	*/

	return button
}

func NewResourceIcon(
	w donburi.World,
	resource domain.Resource,
	pos math.Vec2,
) *donburi.Entry {
	icon := New(w).
		WithPosition(pos).
		WithScale(math.Vec2{X: 2, Y: 2}).
		WithSprite(component.SpriteData{
			Image: domain.Resources[resource].Icon,
		}).entry

	return icon
}

func NewResourceIconWithValue(
	w donburi.World,
	resource domain.Resource,
	value int,
	pos math.Vec2,
) *donburi.Entry {
	icon := NewResourceIcon(w, resource, pos)

	New(w).
		WithParent(icon).
		WithPosition(math.Vec2{X: 38, Y: 24}).
		WithText(component.TextData{
			Text: strconv.Itoa(value),
		})

	return icon
}

func NewSettlementIcon(
	w donburi.World,
	level int,
	pos math.Vec2,
) *donburi.Entry {
	image := SettlementImageFromLevel(level)

	icon := New(w).
		WithPosition(pos).
		WithScale(math.Vec2{X: 2, Y: 2}).
		WithLayer(component.SpriteUILayerUI).
		WithSprite(component.SpriteData{
			Image: image,
		}).
		Entry()

	New(w).
		WithParent(icon).
		WithPosition(math.Vec2{X: 38, Y: 24}).
		WithText(component.TextData{
			Text: fmt.Sprintf("Level %v", level),
		})

	return icon
}

func SettlementImageFromLevel(level int) *ebiten.Image {
	switch level {
	case 1:
		return assets.SettlementLevel1
	case 2:
		return assets.SettlementLevel2
	case 3:
		return assets.SettlementLevel3
	case 4:
		return assets.SettlementLevel4
	case 5:
		return assets.SettlementLevel5
	default:
		panic("unknown level")
	}
}

func SettlementNameFromLevel(level int) string {
	switch level {
	case 1:
		return "Settlement"
	case 2:
		return "Village"
	case 3:
		return "Town"
	case 4:
		return "City"
	case 5:
		return "Castle"
	default:
		panic("unknown level")
	}
}

func NewButtonWithWidth(
	w donburi.World,
	name string,
	pos math.Vec2,
	// TODO deduplicate all button creation
	isActionButton bool,
	onClick func(w donburi.World, e *donburi.Entry),
	width int,
) *donburi.Entry {
	marginLeft := 16
	marginTop := 32
	textHeight := 48

	totalWidth := width + marginLeft*2
	totalHeight := textHeight

	buttonImage := ebiten.NewImage(totalWidth, totalHeight)
	buttonImage.Fill(assets.ButtonColor)

	button := New(w).
		WithPosition(pos).
		WithSprite(component.SpriteData{
			Image: buttonImage,
		}).
		With(component.Collider).
		With(component.Button).
		With(component.Active).
		Entry()

	component.Collider.SetValue(button, component.ColliderData{
		Width:  float64(totalWidth),
		Height: float64(totalHeight),
		Layer:  component.CollisionLayerButtons,
	})
	component.Button.SetValue(button, component.ButtonData{
		OnClick: func(w donburi.World, e *donburi.Entry) {
			events.ButtonClickedEvent.Publish(w, events.ButtonClicked{
				Button: component.Button.Get(e),
			})
			onClick(w, e)
		},
	})
	setButtonActive(w, button, isActionButton)

	New(w).
		WithParent(button).
		WithPosition(math.Vec2{X: float64(marginLeft), Y: float64(marginTop)}).
		WithText(component.TextData{
			Text: name,
		})

	return button
}

func NewButton(
	w donburi.World,
	name string,
	pos math.Vec2,
	// TODO deduplicate all button creation
	isActionButton bool,
	onClick func(w donburi.World, e *donburi.Entry),
) *donburi.Entry {
	face := assets.LargeSquareFont
	textWidth := font.MeasureString(face, name).Round()

	return NewButtonWithWidth(w, name, pos, isActionButton, onClick, textWidth)
}

func NewIconButton(
	w donburi.World,
	icon *ebiten.Image,
	pos math.Vec2,
	// TODO deduplicate all button creation
	isActionButton bool,
	onClick func(w donburi.World, e *donburi.Entry),
) *donburi.Entry {
	width := 64
	height := 64

	buttonImage := ebiten.NewImage(width, height)
	buttonImage.Fill(assets.ButtonColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	buttonImage.DrawImage(icon, op)

	button := New(w).
		WithPosition(pos).
		WithLayer(component.SpriteUILayerButtons).
		WithSprite(component.SpriteData{
			Image: buttonImage,
		}).
		With(component.Collider).
		With(component.Button).
		With(component.Active).
		Entry()

	component.Collider.SetValue(button, component.ColliderData{
		Width:  float64(width),
		Height: float64(height),
		Layer:  component.CollisionLayerButtons,
	})
	component.Button.SetValue(button, component.ButtonData{
		OnClick: func(w donburi.World, e *donburi.Entry) {
			events.ButtonClickedEvent.Publish(w, events.ButtonClicked{
				Button: component.Button.Get(e),
			})
			onClick(w, e)
		},
	})
	setButtonActive(w, button, isActionButton)

	return button
}

func setButtonActive(w donburi.World, button *donburi.Entry, isActionButton bool) {
	game := component.MustFindGame(w)

	component.Active.SetValue(button, component.ActiveData{
		Active: !game.ActionsDisabled || !isActionButton,
	})
}
