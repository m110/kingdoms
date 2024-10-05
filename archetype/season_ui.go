package archetype

import (
	"fmt"
	stdmath "math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"

	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
)

func NewSeasonUIPanel(w donburi.World) *donburi.Entry {
	g := component.MustFindGame(w)
	width := 250
	posX := g.Settings.ScreenWidth - width + 2

	panel := NewFrameWithBorder(w, math.Vec2{X: float64(posX), Y: -2}, width, 80)
	panel.AddComponent(component.UIPanel)
	panel.AddComponent(component.SeasonPanel)
	panel.AddComponent(component.Script)
	component.Layer.Get(panel).Layer = component.SpriteUILayerUI

	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)
	currentDay := progress.Day

	img := ebiten.NewImage(64, 64)
	drawSeasonWheel(img, rotationByDay(currentDay))

	seasonIcon := New(w).
		WithParent(panel).
		WithPosition(math.Vec2{X: 10, Y: 10}).
		WithSprite(component.SpriteData{
			Image: img,
		}).
		With(component.Animation).
		Entry()

	seasonText := New(w).
		WithParent(panel).
		WithPosition(math.Vec2{X: 86, Y: 28}).
		WithText(component.TextData{
			Size: component.TextSizeM,
		}).
		Entry()

	yearText := New(w).
		WithParent(panel).
		WithPosition(math.Vec2{X: 86, Y: 48}).
		WithText(component.TextData{
			Size: component.TextSizeM,
		}).
		Entry()

	updateFunc := func(e *donburi.Entry) {
		component.Text.Get(seasonText).Text = fmt.Sprintf("%v, Day %v", progress.Season.String(), domain.SeasonDay(progress.Day))
		component.Text.Get(yearText).Text = fmt.Sprintf("Year %v", progress.Year)

		if currentDay != progress.Day {
			currentDay = progress.Day
			anim := component.Animation.Get(seasonIcon)
			anim.Start()
		}
	}

	component.Animation.SetValue(seasonIcon, component.AnimationData{
		Active: false,
		Timer:  engine.NewTimer(200 * time.Millisecond),
		Update: func(e *donburi.Entry) {
			anim := component.Animation.Get(e)
			sprite := component.Sprite.Get(e)

			previousDay := progress.Day - 1
			if previousDay <= 0 {
				previousDay = domain.DaysPerYear
			}
			previousRotation := rotationByDay(previousDay)
			currentRotation := rotationByDay(progress.Day)

			if currentRotation == 0 {
				currentRotation = 1
			}

			targetRotation := previousRotation + (currentRotation-previousRotation)*anim.Timer.PercentDone()

			drawSeasonWheel(sprite.Image, targetRotation)

			if anim.Timer.IsReady() {
				anim.Stop()
			}
		},
	})

	component.Script.SetValue(panel, component.ScriptData{
		Update: updateFunc,
	})

	updateFunc(panel)

	return panel
}

func rotationByDay(day int) float64 {
	step := 1 / float64(domain.DaysPerYear)
	start := -step
	return start + step*float64(day)
}

func drawSeasonWheel(target *ebiten.Image, rotation float64) {
	radians := 2 * stdmath.Pi * rotation
	bounds := assets.SeasonsWheel.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
	op.GeoM.Rotate(-radians)
	op.GeoM.Translate(float64(bounds.Dx())/2, float64(bounds.Dy())/2)

	margin := bounds.Dx() / 8
	op.GeoM.Translate(float64(-margin*3), 0)

	target.DrawImage(assets.SeasonsWheel, op)
}
