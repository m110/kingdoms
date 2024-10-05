package archetype

import (
	"fmt"
	"github.com/m110/kingdoms/domain"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
)

type Character int

const (
	CharacterNone Character = iota
	CharacterSettler
	CharacterAdvisor
	CharacterHunter
)

type CharacterData struct {
	Name  string
	Image *ebiten.Image
}

// TODO Can't use globals because assets are not loaded yet
// How to work around this?
func characters() map[Character]CharacterData {
	return map[Character]CharacterData{
		CharacterSettler: {
			Name:  "Settler",
			Image: assets.CharacterSettler,
		},
		CharacterAdvisor: {
			Name:  "Advisor",
			Image: assets.CharacterAdvisor,
		},
		CharacterHunter: {
			Name:  "Hunter",
			Image: assets.CharacterHunter,
		},
	}
}

func ShowMessage(
	w donburi.World,
	from Character,
	content string,
) {
	character, ok := characters()[from]
	if !ok {
		panic(fmt.Sprint("unknown character: ", from))
	}

	player := engine.MustFindComponent[component.PlayerData](w, component.Player)
	addressedBy := domain.PlayerCharacters[player.Character].AddressedBy

	content = strings.Replace(content, "{player}", addressedBy, -1)

	frame := NewFrameWithBorder(w, math.Vec2{X: 200, Y: 100}, 600, 350)
	frame.AddComponent(component.Dialog)
	component.Layer.Get(frame).Layer = component.SpriteUILayerMessage

	iconImage := ebiten.NewImage(128, 128)
	iconImage.Fill(assets.CharacterFrameBackgroundColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(8, 8)

	iconImage.DrawImage(character.Image, op)
	iconImage.DrawImage(assets.CharacterFrame, op)

	icon := New(w).
		WithParent(frame).
		WithPosition(math.Vec2{X: 40, Y: 50}).
		WithSprite(component.SpriteData{
			Image: iconImage,
		}).Entry()

	New(w).
		WithParent(icon).
		WithPosition(math.Vec2{X: 0, Y: 160}).
		WithText(component.TextData{
			Text: character.Name,
		})

	New(w).
		WithPosition(math.Vec2{X: 180, Y: 80}).
		WithParent(frame).
		WithText(component.TextData{
			Text:           content,
			Size:           component.TextSizeM,
			Streaming:      true,
			StreamingTimer: engine.NewTimer(1 * time.Second),
		})

	closeButton := NewButton(w, "Close", math.Vec2{X: 250, Y: 290}, false, func(w donburi.World, e *donburi.Entry) {
		component.Destroy(frame)
	})

	transform.AppendChild(frame, closeButton, false)
}
