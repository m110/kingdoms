package system

import (
	stdmath "math"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
	"github.com/m110/kingdoms/save"
)

const (
	TileSize  = 32
	ChunkSize = 50
)

type Map struct{}

func NewMap() *Map {
	return &Map{}
}

func (m *Map) Init(w donburi.World) {
	events.TechnologyResearchedEvent.Subscribe(w, m.onTechnologyResearched)
}

func (m *Map) Update(w donburi.World) {}

func harvestTiles(w donburi.World) {
	query.NewQuery(filter.Contains(component.Tile)).Each(w, func(entry *donburi.Entry) {
		tile := component.Tile.Get(entry)
		if !tile.CanBeHarvested || tile.Deposit == nil || tile.Resource.Amount == 0 {
			return
		}

		amount := tile.Resource.Amount

		research := component.Research.Get(engine.MustFindWithComponent(w, component.Research))
		bonusFood := 0
		if research.KnownTechnologies[domain.TechnologyGranary] {
			bonusFood = domain.GranaryBonusFood
		}

		bonusStone := 0
		if research.KnownTechnologies[domain.TechnologyMining] {
			bonusStone = domain.MiningBonusStone
		}

		delta := events.ResourcesDelta{}
		var updated bool
		switch tile.Resource.Resource {
		case domain.ResourceFood:
			delta.Food += amount + bonusFood
			updated = true
		case domain.ResourceStone:
			delta.Stone += amount + bonusStone
			updated = true
		case domain.ResourceWood:
			delta.Wood += amount
			updated = true
		}

		if updated {
			UpdateResources(w, delta)

			tile.Resource.Amount -= amount
			if tile.Resource.Amount < 0 {
				panic("negative resource amount")
			}

			tile.CanBeHarvested = false

			updateTileResourcesDisplay(w, entry)

			events.TileResourcesHarvestedEvent.Publish(w, events.TileResourcesHarvested{
				Tile:  entry,
				Delta: delta,
			})
		}

		archetype.UpdateResourceInfoPanel(w)
	})
}

func updateTileResourcesDisplay(w donburi.World, tileEntry *donburi.Entry) {
	_, ok := transform.FindChildWithComponent(tileEntry, component.ResourceAmountText)
	if !ok {
		archetype.New(w).
			WithParent(tileEntry).
			WithLayer(component.SpriteLayerBoardText).
			WithSprite(component.SpriteData{}).
			With(component.ResourceAmountText).
			Entry()
	}

	tile := component.Tile.Get(tileEntry)

	text := engine.MustFindChildWithComponent(tileEntry, component.ResourceAmountText)
	component.Sprite.Get(text).Image = archetype.NewAmountTextImage(tile.Resource.Amount)
	if tile.Resource.Amount <= 0 {
		component.Sprite.Get(text).Hidden = true
	}
}

func (m *Map) onTechnologyResearched(w donburi.World, event events.TechnologyResearched) {
	if event.Technology == domain.TechnologyHunting {
		m.showAnimalsBiomes(w)
	}
}

func (m *Map) showAnimalsBiomes(w donburi.World) {
	query.NewQuery(filter.Contains(component.Tile)).Each(w, func(entry *donburi.Entry) {
		tile := component.Tile.Get(entry)

		if tile.HiddenDeposit == nil || *tile.HiddenDeposit != domain.DepositDeer {
			return
		}

		if tile.HasRoad {
			return
		}

		tile.Deposit = tile.HiddenDeposit
		tile.HiddenDeposit = nil

		depositData := domain.Deposits[*tile.Deposit]

		resource := component.NewRandomTileResource(
			depositData.Resource.Resource,
			depositData.Resource.AmountMin,
			depositData.Resource.AmountMax,
			depositData.Resource.AmountStep,
		)
		tile.Resource = resource

		newDeposit(w, entry)
		updateTileResourcesDisplay(w, entry)
	})

	UpdateAllTilesBuildPermissions(w)
}

func CreateTiles(
	w donburi.World,
	chunks [][]*donburi.Entry,
	terrains [][]domain.Terrain,
	deposits [][]*domain.Deposit,
	width int,
	height int,
	showLabels bool,
	startInSight bool,
) [][]*donburi.Entry {
	tiles := make([][]*donburi.Entry, width)

	for i := 0; i < width; i++ {
		tiles[i] = make([]*donburi.Entry, height)

		for j := 0; j < height; j++ {
			terrain := terrains[i][j]
			deposit := deposits[i][j]

			tile := newTile(w, terrain, deposit, i, j, showLabels, startInSight, nil)

			chunkX := i / ChunkSize
			chunkY := j / ChunkSize
			chunk := chunks[chunkX][chunkY]

			transform.AppendChild(chunk, tile, false)
			tiles[i][j] = tile
		}
	}

	markNeighbors(tiles)
	markBorders(w)

	UpdateAllTilesBuildPermissions(w)

	return tiles
}

func LoadTiles(
	w donburi.World,
	board *save.Board,
	chunks [][]*donburi.Entry,
	showLabels bool,
	startInSight bool,
) [][]*donburi.Entry {
	tiles := make([][]*donburi.Entry, board.Width)

	for i := 0; i < int(board.Width); i++ {
		tiles[i] = make([]*donburi.Entry, board.Height)
	}

	for _, t := range board.Tiles {
		terrain := mapTerrainFromSave(t.Terrain)
		x := int(t.Position.X)
		y := int(t.Position.Y)

		tileResource := component.TileResource{
			Resource:  MapResourceFromSave(t.Resource.Resource),
			Amount:    int(t.Resource.Amount),
			MaxAmount: int(t.Resource.MaxAmount),
		}

		var deposit *domain.Deposit
		if t.Deposit != nil {
			d := mapDepositFromSave(*t.Deposit)
			deposit = &d
		}

		tile := newTile(w, terrain, deposit, x, y, showLabels, startInSight, &tileResource)

		chunkX := x / ChunkSize
		chunkY := y / ChunkSize
		chunk := chunks[chunkX][chunkY]

		transform.AppendChild(chunk, tile, false)
		tiles[x][y] = tile
	}

	markNeighbors(tiles)
	markBorders(w)

	UpdateAllTilesBuildPermissions(w)

	return tiles
}

func markNeighbors(tiles [][]*donburi.Entry) {
	for i := 0; i < len(tiles); i++ {
		for j := 0; j < len(tiles[i]); j++ {
			tile := component.Tile.Get(tiles[i][j])

			if i > 0 {
				tile.NeighborTiles.W = tiles[i-1][j]
			}
			if i < len(tiles)-1 {
				tile.NeighborTiles.E = tiles[i+1][j]
			}
			if j > 0 {
				tile.NeighborTiles.N = tiles[i][j-1]
			}
			if j < len(tiles[i])-1 {
				tile.NeighborTiles.S = tiles[i][j+1]
			}
		}
	}
}

func markBorders(w donburi.World) {
	query.NewQuery(filter.Contains(component.Tile)).Each(w, func(entry *donburi.Entry) {
		tile := component.Tile.Get(entry)

		border, ok := transform.FindChildWithComponent(entry, component.TileBorder)
		if ok {
			// TODO for some reason component.Destroy doesn't work well here
			w.Remove(border.Entity())
		}

		if !tile.InControlRange {
			return
		}

		if !tile.HasBorder() {
			return
		}

		archetype.New(w).
			WithParent(entry).
			WithLayer(component.SpriteLayerBorder).
			WithSprite(component.SpriteData{
				Image: tile.BorderImage(),
				AlphaOverride: &component.AlphaOverride{
					A: 0.8,
				},
			}).
			With(component.TileBorder)
	})
}

func newTile(
	w donburi.World,
	terrain domain.Terrain,
	deposit *domain.Deposit,
	x int,
	y int,
	showLabel bool,
	startInSight bool,
	overrideTileResource *component.TileResource,
) *donburi.Entry {
	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)
	terrainData := domain.Terrains[terrain]

	tile := archetype.New(w).
		WithPosition(math.Vec2{
			X: float64(x % ChunkSize * TileSize),
			Y: float64(y % ChunkSize * TileSize),
		}).
		WithScale(math.Vec2{X: 2, Y: 2}).
		WithLayer(component.SpriteLayerBackground).
		WithSprite(component.SpriteData{
			Image: terrainData.Sprites.BySeason(progress.Season),
		}).
		With(component.SpriteTransition).
		With(component.Collider).
		With(component.Selectable).
		With(component.Tile).
		Entry()

	var resource component.TileResource
	if overrideTileResource != nil {
		resource = *overrideTileResource
	} else if deposit != nil {
		depositData := domain.Deposits[*deposit]
		if !depositData.Discoverable {
			resource = component.NewRandomTileResource(
				depositData.Resource.Resource,
				depositData.Resource.AmountMin,
				depositData.Resource.AmountMax,
				depositData.Resource.AmountStep,
			)
		}
	}

	component.Collider.SetValue(tile, component.ColliderData{
		Width:  16,
		Height: 16,
		Layer:  component.CollisionLayerTiles,
	})

	tileData := component.TileData{
		Position: component.Position{
			X: x,
			Y: y,
		},
		Terrain:      terrain,
		Resource:     resource,
		InSightRange: startInSight,
	}

	if deposit != nil {
		depositData := domain.Deposits[*deposit]
		if depositData.Discoverable {
			tileData.HiddenDeposit = deposit
		} else {
			tileData.Deposit = deposit
		}
	}

	component.Tile.SetValue(tile, tileData)

	if deposit != nil {
		depositData := domain.Deposits[*deposit]
		if !depositData.Discoverable {
			newDeposit(w, tile)
		}
	}

	// TODO Dumb that it's an image, not text?
	if showLabel && resource.Amount > 0 {
		updateTileResourcesDisplay(w, tile)
	}

	if !startInSight {
		fog := archetype.New(w).
			WithParent(tile).
			WithLayer(component.SpriteLayerFog).
			WithSprite(component.SpriteData{
				Image: assets.TileFog,
			}).
			With(component.TileFog).
			With(component.Animation).
			Entry()

		component.Animation.SetValue(fog, component.AnimationData{
			Active: false,
			Update: func(e *donburi.Entry) {
				anim := component.Animation.Get(e)
				if anim.Timer.IsReady() {
					component.Destroy(e)
					return
				}

				sprite := component.Sprite.Get(e)
				sprite.AlphaOverride = &component.AlphaOverride{
					A: 1 - anim.Timer.PercentDone(),
				}
			},
		})
	}

	return tile
}

func newDeposit(w donburi.World, tileEntry *donburi.Entry) {
	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)
	tile := component.Tile.Get(tileEntry)

	archetype.New(w).
		WithParent(tileEntry).
		With(component.TileDeposit).
		WithLayer(component.SpriteLayerForeground).
		WithSprite(component.SpriteData{
			Image: tile.DepositImage(progress.Season),
		}).
		With(component.SpriteTransition)
}

func GenerateChunks(
	w donburi.World,
	parent *donburi.Entry,
	width int,
	height int,
) [][]*donburi.Entry {
	chunksWidth := int(stdmath.Ceil(float64(width) / float64(ChunkSize)))
	chunksHeight := int(stdmath.Ceil(float64(height) / float64(ChunkSize)))

	chunks := make([][]*donburi.Entry, chunksWidth)
	for i := 0; i < chunksWidth; i++ {
		chunks[i] = make([]*donburi.Entry, chunksHeight)
		for j := 0; j < chunksHeight; j++ {
			pos := math.Vec2{
				X: float64(i * ChunkSize * TileSize),
				Y: float64(j * ChunkSize * TileSize),
			}
			chunk := newChunk(w, pos)

			transform.AppendChild(parent, chunk, false)

			chunks[i][j] = chunk
		}
	}

	return chunks
}

func newChunk(w donburi.World, pos math.Vec2) *donburi.Entry {
	chunk := archetype.New(w).
		WithPosition(pos).
		With(component.Chunk).
		Entry()

	component.Chunk.SetValue(chunk, component.ChunkData{
		Size: component.Size{
			Width:  ChunkSize * TileSize,
			Height: ChunkSize * TileSize,
		},
	})

	return chunk
}
