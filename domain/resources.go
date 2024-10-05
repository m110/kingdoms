package domain

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/assets"
)

type Resource int

const (
	ResourceWood Resource = iota
	ResourceStone
	ResourceFood
	ResourceIron
	ResourceGold
	ResourceDiamonds
)

type ResourceData struct {
	Resource Resource
	Icon     *ebiten.Image
}

var Resources map[Resource]ResourceData

func loadResources() {
	Resources = map[Resource]ResourceData{
		ResourceWood: {
			Resource: ResourceWood,
			Icon:     assets.IconWood,
		},
		ResourceStone: {
			Resource: ResourceStone,
			Icon:     assets.IconStone,
		},
		ResourceFood: {
			Resource: ResourceFood,
			Icon:     assets.IconFood,
		},
		ResourceIron: {
			Resource: ResourceIron,
			Icon:     assets.IconIron,
		},
		ResourceGold: {
			Resource: ResourceGold,
			Icon:     assets.IconCurrency,
		},
		ResourceDiamonds: {
			Resource: ResourceDiamonds,
			Icon:     assets.IconDiamond,
		},
	}
}

type ResourceCost struct {
	Food  int
	Stone int
	Wood  int
	Iron  int
	Gold  int
}
