package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/scene"
)

type SandboxConfig struct {
	ScreenWidth  int
	ScreenHeight int
}

type Sandbox struct {
	config SandboxConfig
	scene  Scene
}

func NewSandbox(config SandboxConfig) *Sandbox {
	assets.MustLoadAssets()
	domain.LoadDomain()

	context := scene.Context{
		ScreenWidth:  config.ScreenWidth,
		ScreenHeight: config.ScreenHeight,
	}
	return &Sandbox{
		scene: scene.NewSandbox(context),
	}
}

func (s *Sandbox) Update() error {
	s.scene.Update()
	return nil
}

func (s *Sandbox) Draw(screen *ebiten.Image) {
	s.scene.Draw(screen)
}

func (s *Sandbox) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
