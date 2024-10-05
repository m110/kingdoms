package archetype

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"golang.org/x/image/font"

	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
)

func NewAmountTextImage(amount int) *ebiten.Image {
	amountText := strconv.Itoa(amount)
	image := ebiten.NewImage(16, 16)

	width := font.MeasureString(assets.SmallSquareFont, amountText)

	text.Draw(image, amountText, assets.SmallSquareFont, 15-width.Round(), 15, assets.TextColor)

	return image
}

func NewTileFloatingAnimation(w donburi.World, parent *donburi.Entry, image *ebiten.Image) {
	animation := New(w).
		WithParent(parent).
		With(component.Velocity).
		With(component.TimeToLive).
		With(component.Animation).
		WithLayer(component.SpriteLayerForeground).
		WithSprite(component.SpriteData{
			Image: image,
		}).Entry()

	component.Velocity.SetValue(animation, component.VelocityData{
		Velocity: math.Vec2{
			X: 0,
			Y: -1,
		},
	})
	component.TimeToLive.SetValue(animation, component.TimeToLiveData{
		Timer: engine.NewTimer(1 * time.Second),
	})
	component.Animation.SetValue(animation, component.AnimationData{
		Active: true,
		Timer:  engine.NewTimer(1 * time.Second),
		Update: func(e *donburi.Entry) {
			anim := component.Animation.Get(e)
			sprite := component.Sprite.Get(e)

			if sprite.AlphaOverride == nil {
				sprite.AlphaOverride = &component.AlphaOverride{}
			}

			sprite.AlphaOverride.A = 1.0 - anim.Timer.PercentDone()
		},
	})
}
