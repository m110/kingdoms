package component

import "github.com/yohamta/donburi"

type SelectableData struct {
	Selected bool
}

var Selectable = donburi.NewComponentType[SelectableData]()
