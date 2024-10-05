package component

import "github.com/yohamta/donburi"

type TableData struct {
	Rows [][]*donburi.Entry
}

var Table = donburi.NewComponentType[TableData]()
