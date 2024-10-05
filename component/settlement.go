package component

import "github.com/yohamta/donburi"

type SettlementData struct {
	ID    int
	Level int
}

var Settlement = donburi.NewComponentType[SettlementData]()
