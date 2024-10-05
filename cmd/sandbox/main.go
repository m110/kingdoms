package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/game"
	"log"
)

var (
	screenWidth  = 1300
	screenHeight = 800
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)

	err := ebiten.RunGame(game.NewSandbox(game.SandboxConfig{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}))
	if err != nil {
		log.Fatal(err)
	}
}
