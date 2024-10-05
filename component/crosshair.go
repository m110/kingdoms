package component

import "github.com/yohamta/donburi"

type CrosshairData struct {
	Position Position
}

var Crosshair = donburi.NewComponentType[CrosshairData]()
