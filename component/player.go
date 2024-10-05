package component

import (
	"github.com/m110/kingdoms/domain"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	Character domain.PlayerCharacterType
}

var Player = donburi.NewComponentType[PlayerData]()
