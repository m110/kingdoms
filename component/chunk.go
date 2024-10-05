package component

import "github.com/yohamta/donburi"

type ChunkData struct {
	Size Size
}

var Chunk = donburi.NewComponentType[ChunkData]()
