package component

import (
	"github.com/yohamta/donburi"

	"github.com/m110/kingdoms/domain"
)

type ProgressData struct {
	Season domain.Season
	Day    int
	Year   int
}

var Progress = donburi.NewComponentType[ProgressData]()
