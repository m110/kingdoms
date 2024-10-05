package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/yohamta/donburi"
)

type TileResource struct {
	Resource  domain.Resource
	Amount    int
	MaxAmount int
}

func NewTileResource(resource domain.Resource, amount int) TileResource {
	return TileResource{
		Resource:  resource,
		Amount:    amount,
		MaxAmount: amount,
	}
}

func NewRandomTileResource(resource domain.Resource, minAmount int, maxAmount int, step int) TileResource {
	amount := engine.RandomIntRange(minAmount, maxAmount)
	amount -= amount % step
	return NewTileResource(resource, amount)
}

// TODO Probably could be a separate component
type Position struct {
	X, Y int
}

type NeighborTiles struct {
	N, E, S, W *donburi.Entry
}

type NeighborTile struct {
	Entry     *donburi.Entry
	Direction Direction
}

func (n NeighborTiles) All() []NeighborTile {
	var tiles []NeighborTile
	if n.N != nil {
		tiles = append(tiles, NeighborTile{
			Entry:     n.N,
			Direction: DirectionNorth,
		})
	}
	if n.E != nil {
		tiles = append(tiles, NeighborTile{
			Entry:     n.E,
			Direction: DirectionEast,
		})
	}
	if n.S != nil {
		tiles = append(tiles, NeighborTile{
			Entry:     n.S,
			Direction: DirectionSouth,
		})
	}
	if n.W != nil {
		tiles = append(tiles, NeighborTile{
			Entry:     n.W,
			Direction: DirectionWest,
		})
	}
	return tiles
}

type TileData struct {
	Position      Position
	Terrain       domain.Terrain
	Deposit       *domain.Deposit
	HiddenDeposit *domain.Deposit
	Resource      TileResource

	NeighborTiles NeighborTiles

	CanBeHarvested  bool
	InControlRange  bool
	InSightRange    bool
	HasRoad         bool
	HasBuilding     bool
	HasNeighborRoad bool

	CanBuildRoad      bool
	CanBuildBuildings bool
}

func (t *TileData) ResourceHarvestModifierBySeason(season domain.Season) float64 {
	if t.Deposit == nil {
		panic("no deposit")
	}
	deposit := domain.Deposits[*t.Deposit]

	return deposit.HarvestModifiers.BySeason(season)
}

func (t *TileData) RoadImage() *ebiten.Image {
	roadN := tileHasRoad(t.NeighborTiles.N)
	roadS := tileHasRoad(t.NeighborTiles.S)
	roadW := tileHasRoad(t.NeighborTiles.W)
	roadE := tileHasRoad(t.NeighborTiles.E)

	if roadN && roadS && roadW && roadE {
		return assets.RoadNSWE
	}
	if roadN && roadS && roadW {
		return assets.RoadNSW
	}
	if roadN && roadS && roadE {
		return assets.RoadNSE
	}
	if roadN && roadW && roadE {
		return assets.RoadNWE
	}
	if roadS && roadW && roadE {
		return assets.RoadSWE
	}
	if roadN && roadS {
		return assets.RoadNS
	}
	if roadW && roadE {
		return assets.RoadWE
	}
	if roadN && roadW {
		return assets.RoadNW
	}
	if roadN && roadE {
		return assets.RoadNE
	}
	if roadS && roadW {
		return assets.RoadSW
	}
	if roadS && roadE {
		return assets.RoadSE
	}
	if roadN {
		return assets.RoadN
	}
	if roadS {
		return assets.RoadS
	}
	if roadW {
		return assets.RoadW
	}
	if roadE {
		return assets.RoadE
	}

	return assets.Road
}

func (t *TileData) DepositImage(season domain.Season) *ebiten.Image {
	if t.Deposit == nil {
		return nil
	}

	depositData := domain.Deposits[*t.Deposit]
	percentFull := float64(t.Resource.Amount) / float64(t.Resource.MaxAmount)

	return depositData.Sprites.BySeasonAndAmountPercent(season, percentFull)
}

func (t *TileData) HasBorder() bool {
	return borderInDirection(t.NeighborTiles.N) || borderInDirection(t.NeighborTiles.S) || borderInDirection(t.NeighborTiles.W) || borderInDirection(t.NeighborTiles.E)
}

func (t *TileData) BorderImage() *ebiten.Image {
	borderN := borderInDirection(t.NeighborTiles.N)
	borderS := borderInDirection(t.NeighborTiles.S)
	borderW := borderInDirection(t.NeighborTiles.W)
	borderE := borderInDirection(t.NeighborTiles.E)

	bounds := assets.TileBorderN.Bounds()
	img := ebiten.NewImage(bounds.Dx(), bounds.Dy())

	if borderN {
		img.DrawImage(assets.TileBorderN, nil)
	}

	if borderE {
		img.DrawImage(assets.TileBorderE, nil)
	}

	if borderS {
		img.DrawImage(assets.TileBorderS, nil)
	}

	if borderW {
		img.DrawImage(assets.TileBorderW, nil)
	}

	return img
}

func borderInDirection(neighborTile *donburi.Entry) bool {
	if neighborTile == nil {
		return true
	}

	return !Tile.Get(neighborTile).InControlRange
}

func tileHasRoad(tileEntry *donburi.Entry) bool {
	if tileEntry == nil {
		return false
	}

	return Tile.Get(tileEntry).HasRoad
}

var Tile = donburi.NewComponentType[TileData]()
