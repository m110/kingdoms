package domain

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/m110/kingdoms/assets"
)

type Terrain int

const (
	TerrainPlains Terrain = iota
	TerrainForest
	TerrainDesert
	TerrainMountains
	TerrainWater
	TerrainShore
	TerrainAnimals
)

type TerrainData struct {
	Terrain Terrain
	Name    string
	Sprites TerrainSprites

	CanBuildBuildings bool
	CanBuildRoads     bool
}

type TerrainSprites struct {
	Spring *ebiten.Image
	Summer *ebiten.Image
	Fall   *ebiten.Image
	Winter *ebiten.Image
}

func (s TerrainSprites) BySeason(season Season) *ebiten.Image {
	switch season {
	case SeasonSpring:
		return s.Spring
	case SeasonSummer:
		return s.Summer
	case SeasonFall:
		return s.Fall
	case SeasonWinter:
		return s.Winter
	default:
		panic("unknown season")
	}
}

var Terrains map[Terrain]TerrainData

func loadTerrains() {
	Terrains = map[Terrain]TerrainData{
		TerrainPlains: {
			Terrain: TerrainPlains,
			Name:    "Plains",
			Sprites: TerrainSprites(assets.TerrainPlains),
		},
		TerrainForest: {
			Terrain: TerrainForest,
			Name:    "Forest",
			Sprites: TerrainSprites{
				Spring: assets.TerrainForest,
				Summer: assets.TerrainForest,
				Fall:   assets.TerrainForest,
				Winter: assets.TerrainForest,
			},
		},
		TerrainDesert: {
			Terrain: TerrainDesert,
			Name:    "Desert",
			Sprites: TerrainSprites{
				Spring: assets.TerrainDesert,
				Summer: assets.TerrainDesert,
				Fall:   assets.TerrainDesert,
				Winter: assets.TerrainDesert,
			},
		},
		TerrainMountains: {
			Terrain: TerrainMountains,
			Name:    "Mountains",
			Sprites: TerrainSprites{
				Spring: assets.TerrainMountains,
				Summer: assets.TerrainMountains,
				Fall:   assets.TerrainMountains,
				Winter: assets.TerrainMountains,
			},
		},
		TerrainWater: {
			Terrain: TerrainWater,
			Name:    "Sea",
			Sprites: TerrainSprites(assets.TerrainWater),
		},
		TerrainShore: {
			Terrain: TerrainShore,
			Name:    "Shore",
			Sprites: TerrainSprites(assets.TerrainShore),
		},
	}
}

type Deposit int

const (
	DepositTrees Deposit = iota
	DepositRocks
	DepositIronOre
	DepositBerries
	DepositMushrooms
	DepositDeer
)

type DepositData struct {
	Deposit          Deposit
	Name             string
	Sprites          DepositSprites
	Resource         DepositResource
	HarvestModifiers DepositHarvestModifiers
	Discoverable     bool
}

type DepositResource struct {
	Resource   Resource
	AmountMin  int
	AmountMax  int
	AmountStep int
}

type DepositSprites struct {
	Spring SpritesByAmount
	Summer SpritesByAmount
	Fall   SpritesByAmount
	Winter SpritesByAmount
}

func (s DepositSprites) BySeasonAndAmountPercent(season Season, amountPercent float64) *ebiten.Image {
	var sprites SpritesByAmount
	switch season {
	case SeasonSpring:
		sprites = s.Spring
	case SeasonSummer:
		sprites = s.Summer
	case SeasonFall:
		sprites = s.Fall
	case SeasonWinter:
		sprites = s.Winter
	default:
		panic("unknown season")
	}

	if amountPercent >= 0.75 {
		return sprites.Full
	} else if amountPercent > 0 {
		return sprites.Half
	} else {
		return sprites.Empty
	}
}

type DepositHarvestModifiers struct {
	Spring float64
	Summer float64
	Fall   float64
	Winter float64
}

func (m DepositHarvestModifiers) BySeason(season Season) float64 {
	switch season {
	case SeasonSpring:
		return m.Spring
	case SeasonSummer:
		return m.Summer
	case SeasonFall:
		return m.Fall
	case SeasonWinter:
		return m.Winter
	default:
		panic("unknown season")
	}
}

type SpritesByAmount struct {
	Full  *ebiten.Image
	Half  *ebiten.Image
	Empty *ebiten.Image
}

var Deposits map[Deposit]DepositData

func loadDeposits() {
	Deposits = map[Deposit]DepositData{
		DepositTrees: {
			Deposit: DepositTrees,
			Name:    "Trees",
			Sprites: DepositSprites{
				Spring: SpritesByAmount(assets.DepositTrees.Spring),
				Summer: SpritesByAmount(assets.DepositTrees.Summer),
				Fall:   SpritesByAmount(assets.DepositTrees.Fall),
				Winter: SpritesByAmount(assets.DepositTrees.Winter),
			},
			Resource: DepositResource{
				Resource:   ResourceWood,
				AmountMin:  10,
				AmountMax:  20,
				AmountStep: 5,
			},
			HarvestModifiers: DepositHarvestModifiers{
				Spring: 1.0,
				Summer: 1.0,
				Fall:   1.0,
				Winter: 1.0,
			},
		},
		DepositRocks: {
			Deposit: DepositRocks,
			Name:    "Rocks",
			Sprites: DepositSprites{
				Spring: SpritesByAmount(assets.DepositRocks.Spring),
				Summer: SpritesByAmount(assets.DepositRocks.Summer),
				Fall:   SpritesByAmount(assets.DepositRocks.Fall),
				Winter: SpritesByAmount(assets.DepositRocks.Winter),
			},
			Resource: DepositResource{
				Resource:   ResourceStone,
				AmountMin:  10,
				AmountMax:  30,
				AmountStep: 5,
			},
			HarvestModifiers: DepositHarvestModifiers{
				Spring: 1.0,
				Summer: 1.0,
				Fall:   1.0,
				Winter: 1.0,
			},
		},
		DepositIronOre: {
			Deposit: DepositIronOre,
			Name:    "Iron Ore",
			Sprites: DepositSprites{
				Spring: SpritesByAmount(assets.DepositIronOre.Spring),
				Summer: SpritesByAmount(assets.DepositIronOre.Summer),
				Fall:   SpritesByAmount(assets.DepositIronOre.Fall),
				Winter: SpritesByAmount(assets.DepositIronOre.Winter),
			},
			Resource: DepositResource{
				Resource:   ResourceIron,
				AmountMin:  10,
				AmountMax:  30,
				AmountStep: 5,
			},
			HarvestModifiers: DepositHarvestModifiers{
				Spring: 1.0,
				Summer: 1.0,
				Fall:   1.0,
				Winter: 1.0,
			},
			Discoverable: true,
		},
		DepositBerries: {
			Deposit: DepositBerries,
			Name:    "Berries",
			Sprites: DepositSprites{
				Spring: SpritesByAmount(assets.DepositBerries.Spring),
				Summer: SpritesByAmount(assets.DepositBerries.Summer),
				Fall:   SpritesByAmount(assets.DepositBerries.Fall),
				Winter: SpritesByAmount(assets.DepositBerries.Winter),
			},
			Resource: DepositResource{
				Resource:   ResourceFood,
				AmountMin:  5,
				AmountMax:  10,
				AmountStep: 5,
			},
			HarvestModifiers: DepositHarvestModifiers{
				Spring: 1.0,
				Summer: 2.0,
				Fall:   2.0,
				Winter: 0.0,
			},
		},
		DepositMushrooms: {
			Deposit: DepositMushrooms,
			Name:    "Mushrooms",
			Sprites: DepositSprites{
				Spring: SpritesByAmount(assets.DepositMushrooms.Spring),
				Summer: SpritesByAmount(assets.DepositMushrooms.Summer),
				Fall:   SpritesByAmount(assets.DepositMushrooms.Fall),
				Winter: SpritesByAmount(assets.DepositMushrooms.Winter),
			},
			Resource: DepositResource{
				Resource:   ResourceFood,
				AmountMin:  5,
				AmountMax:  10,
				AmountStep: 5,
			},
			HarvestModifiers: DepositHarvestModifiers{
				Spring: 1.0,
				Summer: 2.0,
				Fall:   2.0,
				Winter: 0,
			},
		},
		DepositDeer: {
			Deposit: DepositDeer,
			Name:    "Deer",
			Sprites: DepositSprites{
				Spring: SpritesByAmount(assets.DepositDeer.Spring),
				Summer: SpritesByAmount(assets.DepositDeer.Summer),
				Fall:   SpritesByAmount(assets.DepositDeer.Fall),
				Winter: SpritesByAmount(assets.DepositDeer.Winter),
			},
			Resource: DepositResource{
				Resource:   ResourceFood,
				AmountMin:  10,
				AmountMax:  15,
				AmountStep: 5,
			},
			HarvestModifiers: DepositHarvestModifiers{
				Spring: 1.0,
				Summer: 1.0,
				Fall:   1.0,
				Winter: 1.0,
			},
			Discoverable: true,
		},
	}
}
