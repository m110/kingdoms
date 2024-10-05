package component

import (
	"github.com/yohamta/donburi"
)

type Size struct {
	Width  int
	Height int
}

type BoardData struct {
	Size     Size
	TileSize Size

	Tiles [][]*donburi.Entry

	StartPosition Position
}

func (b *BoardData) TileByPosition(pos Position) (*donburi.Entry, bool) {
	if pos.X < 0 || pos.X >= b.Size.Width || pos.Y < 0 || pos.Y >= b.Size.Height {
		return nil, false
	}

	return b.Tiles[pos.X][pos.Y], true
}

func (b *BoardData) MustTileByPosition(pos Position) *donburi.Entry {
	tile, ok := b.TileByPosition(pos)
	if !ok {
		panic("tile not found")
	}
	return tile
}

var Board = donburi.NewComponentType[BoardData]()
