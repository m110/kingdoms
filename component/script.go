package component

import "github.com/yohamta/donburi"

type ScriptData struct {
	Update func(e *donburi.Entry)
}

var Script = donburi.NewComponentType[ScriptData]()
