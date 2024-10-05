package archetype

import (
	"github.com/yohamta/donburi"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
)

func MustFindBoard(w donburi.World) *component.BoardData {
	board := engine.MustFindWithComponent(w, component.Board)
	return component.Board.Get(board)
}
