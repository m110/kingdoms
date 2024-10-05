package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"

	"github.com/m110/kingdoms/game"
)

func init() {
	mobile.SetGame(game.NewGame(game.Config{
		StartScene:   "",
		ScreenWidth:  1300,
		ScreenHeight: 600,
		SampleRate:   44100,
	}))
}

func Dummy() {}
