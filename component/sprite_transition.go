package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SpriteTransitionData struct {
	Active bool
	From   *ebiten.Image
	To     *ebiten.Image
}

var SpriteTransition = donburi.NewComponentType[SpriteTransitionData]()
