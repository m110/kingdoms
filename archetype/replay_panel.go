package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func NewReplayPanel(w donburi.World, rewindFunc func()) {
	game := component.MustFindGame(w)

	frame := NewFrameWithBorder(w, math.Vec2{X: 750, Y: float64(game.Settings.ScreenHeight) - 130}, 350, 130)
	component.Layer.Get(frame).Layer = component.SpriteUILayerUI
	frame.AddComponent(component.UIPanel)

	rewindButton := NewButton(w, "|<<", math.Vec2{X: 30, Y: 60}, false, func(w donburi.World, e *donburi.Entry) {
		rewindFunc()
	})
	transform.AppendChild(frame, rewindButton, false)

	pauseButton := NewButton(w, "||", math.Vec2{X: 130, Y: 60}, false, func(w donburi.World, e *donburi.Entry) {})
	pauseIndicator := activeButtonIndicator(w, pauseButton)
	transform.AppendChild(frame, pauseButton, false)

	playButton := NewButton(w, ">", math.Vec2{X: 200, Y: 60}, false, func(w donburi.World, e *donburi.Entry) {})
	playIndicator := activeButtonIndicator(w, playButton)
	transform.AppendChild(frame, playButton, false)

	doubleSpeedButton := NewButton(w, ">>", math.Vec2{X: 270, Y: 60}, false, func(w donburi.World, e *donburi.Entry) {})
	doubleSpeedIndicator := activeButtonIndicator(w, doubleSpeedButton)
	transform.AppendChild(frame, doubleSpeedButton, false)

	playIndicator.Active = true

	player := engine.MustFindComponent[component.ActionsPlayerData](w, component.ActionsPlayer)

	component.Button.Get(pauseButton).OnClick = func(w donburi.World, e *donburi.Entry) {
		player.Pause()

		pauseIndicator.Active = true
		playIndicator.Active = false
		doubleSpeedIndicator.Active = false
	}

	component.Button.Get(playButton).OnClick = func(w donburi.World, e *donburi.Entry) {
		player.Speed = 1.0
		player.Play()

		pauseIndicator.Active = false
		playIndicator.Active = true
		doubleSpeedIndicator.Active = false
	}

	component.Button.Get(doubleSpeedButton).OnClick = func(w donburi.World, e *donburi.Entry) {
		player.Speed = 10.0
		player.Play()

		pauseIndicator.Active = false
		playIndicator.Active = false
		doubleSpeedIndicator.Active = true
	}

	var maxTick int
	if len(player.Actions) == 0 {
		maxTick = 1
	} else {
		maxTick = player.Actions[len(player.Actions)-1].Tick
	}

	progressWidth := 300

	progressImage := ebiten.NewImage(progressWidth, 5)
	progressImage.Fill(assets.IconFrameLightColor)

	markImage := ebiten.NewImage(5, 20)
	markImage.Fill(assets.IconFrameLightColor)

	progress := New(w).
		WithParent(frame).
		WithPosition(math.Vec2{X: 30, Y: 20}).
		WithSprite(component.SpriteData{
			Image: progressImage,
		}).
		Entry()

	progressMark := New(w).
		WithPosition(math.Vec2{X: 0, Y: -8}).
		WithParent(progress).
		With(component.Animation).
		WithSprite(component.SpriteData{
			Image: markImage,
		}).
		Entry()

	component.Animation.SetValue(progressMark, component.AnimationData{
		Active: true,
		Update: func(e *donburi.Entry) {
			t := transform.GetTransform(progressMark)
			pos := t.LocalPosition
			pos.X = float64(int(player.CurrentTick)) / float64(maxTick) * float64(progressWidth-5)
			t.LocalPosition = pos
		},
	})
}

func activeButtonIndicator(w donburi.World, button *donburi.Entry) *component.ActiveData {
	component.Layer.Get(button).Layer = component.SpriteUILayerButtons
	collider := component.Collider.Get(button)

	img := ebiten.NewImage(int(collider.Width), 20)
	img.Fill(assets.IconFrameLightColor)

	indicator := New(w).
		WithParent(button).
		With(component.Active).
		WithPosition(math.Vec2{X: 0, Y: 35}).
		WithSprite(component.SpriteData{
			Image: img,
		}).
		Entry()

	// TODO Not sure if hack?
	component.Layer.Get(indicator).Layer = component.SpriteUILayerUI + 1

	return component.Active.Get(indicator)
}
