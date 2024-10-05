package component

import (
	"github.com/m110/kingdoms/domain"
	"github.com/yohamta/donburi"
)

type ResearchData struct {
	KnownTechnologies map[domain.Technology]bool
}

var Research = donburi.NewComponentType[ResearchData]()
