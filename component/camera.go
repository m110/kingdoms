package component

import (
	"github.com/m110/kingdoms/engine"
	"github.com/yohamta/donburi"
)

type CameraData struct {
	Zoom engine.FloatRange
}

var Camera = donburi.NewComponentType[CameraData]()
