package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/m110/kingdoms/game"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screenWidth  = 1300
	screenHeight = 600

	sampleRate = 44100
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)

	cpuProfile := os.Getenv("CPU_PROFILE")
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			panic(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	startScene := os.Getenv("START_SCENE")

	err := ebiten.RunGame(game.NewGame(game.Config{
		StartScene:   startScene,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		SampleRate:   sampleRate,
	}))
	if err != nil {
		log.Fatal(err)
	}
}
