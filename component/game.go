package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type GameData struct {
	Settings Settings

	Paused          bool
	ActionsDisabled bool
}

type Settings struct {
	ScreenWidth  int
	ScreenHeight int
}

var Game = donburi.NewComponentType[GameData]()

func MustFindGame(w donburi.World) *GameData {
	game, ok := query.NewQuery(filter.Contains(Game)).FirstEntity(w)
	if !ok {
		panic("game not found")
	}
	return Game.Get(game)
}
