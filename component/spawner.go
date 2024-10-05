package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/kingdoms/engine"
)

type SpawnerData struct {
	Count     engine.IntRange
	Interval  engine.DurationRange
	SpawnFunc func(w donburi.World) *donburi.Entry

	Timer *engine.Timer
}

var Spawner = donburi.NewComponentType[SpawnerData]()
