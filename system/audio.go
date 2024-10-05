package system

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/yohamta/donburi"

	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/events"
)

type Audio struct {
	audioContext *audio.Context

	click1Player *audio.Player
	click2Player *audio.Player
	click3Player *audio.Player

	winPlayer  *audio.Player
	losePlayer *audio.Player

	chime1Player *audio.Player
	chime2Player *audio.Player
	chime3Player *audio.Player
	chime4Player *audio.Player
}

func NewAudio() *Audio {
	ctx := audio.CurrentContext()

	return &Audio{
		audioContext: ctx,

		click1Player: ctx.NewPlayerFromBytes(assets.Click1Sound),
		click2Player: ctx.NewPlayerFromBytes(assets.Click2Sound),
		click3Player: ctx.NewPlayerFromBytes(assets.Click3Sound),

		winPlayer:  ctx.NewPlayerFromBytes(assets.WinSound),
		losePlayer: ctx.NewPlayerFromBytes(assets.LoseSound),

		chime1Player: ctx.NewPlayerFromBytes(assets.Chime1Sound),
		chime2Player: ctx.NewPlayerFromBytes(assets.Chime2Sound),
		chime3Player: ctx.NewPlayerFromBytes(assets.Chime3Sound),
		chime4Player: ctx.NewPlayerFromBytes(assets.Chime4Sound),
	}
}

func (a *Audio) Init(w donburi.World) {
	events.RoadSpawnedEvent.Subscribe(w, a.onRoadSpawned)
	events.OriginSettlementSpawnedEvent.Subscribe(w, a.onOriginSettlementSpawned)
	events.SettlementSpawnedEvent.Subscribe(w, a.onSettlementSpawned)
	events.SettlementUpgradedEvent.Subscribe(w, a.onSettlementUpgraded)
	events.ButtonClickedEvent.Subscribe(w, a.onButtonClicked)
	events.GameOverEvent.Subscribe(w, a.onGameOver)
	events.GameWonEvent.Subscribe(w, a.onGameWon)
}

func (a *Audio) Update(w donburi.World) {}

func (a *Audio) onRoadSpawned(w donburi.World, event events.RoadSpawned) {
	_ = a.click2Player.Rewind()
	a.click2Player.Play()
}

func (a *Audio) onOriginSettlementSpawned(w donburi.World, event events.OriginSettlementSpawned) {
	// TODO Another sound
	_ = a.click3Player.Rewind()
	a.click3Player.Play()
}

func (a *Audio) onSettlementSpawned(w donburi.World, event events.SettlementSpawned) {
	_ = a.click3Player.Rewind()
	a.click3Player.Play()
}

func (a *Audio) onButtonClicked(w donburi.World, event events.ButtonClicked) {
	switch event.Button.Sound {
	case component.ButtonSoundNone:
	case component.ButtonSoundDefault:
		_ = a.click1Player.Rewind()
		a.click1Player.Play()
	}
}

func (a *Audio) onSettlementUpgraded(w donburi.World, event events.SettlementUpgraded) {
	sett := component.Settlement.Get(event.Settlement)
	switch sett.Level {
	case 2:
		_ = a.chime1Player.Rewind()
		a.chime1Player.Play()
	case 3:
		_ = a.chime2Player.Rewind()
		a.chime2Player.Play()
	case 4:
		_ = a.chime3Player.Rewind()
		a.chime3Player.Play()
	case 5:
		_ = a.chime4Player.Rewind()
		a.chime4Player.Play()
	}
}

func (a *Audio) onGameOver(w donburi.World, event events.GameOver) {
	_ = a.losePlayer.Rewind()
	a.losePlayer.Play()
}

func (a *Audio) onGameWon(w donburi.World, event events.GameWon) {
	_ = a.winPlayer.Rewind()
	a.winPlayer.Play()
}
