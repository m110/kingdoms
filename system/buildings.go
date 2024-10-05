package system

import (
	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
)

const (
	fogRevealDuration = 150 * time.Millisecond
)

type Buildings struct {
	board *component.BoardData

	roads       [][]*entryWithPos
	settlements [][]*entryWithPos

	resources *component.ResourcesData

	tilesQuery       *donburi.Query
	roadsQuery       *donburi.Query
	settlementsQuery *donburi.Query
}

func NewBuildings() *Buildings {
	return &Buildings{
		tilesQuery: donburi.NewQuery(
			filter.Contains(
				component.Tile,
			),
		),
		roadsQuery: donburi.NewQuery(
			filter.Contains(
				component.Road,
			),
		),
		settlementsQuery: donburi.NewQuery(
			filter.Contains(
				component.Settlement,
			),
		),
	}
}

func (b *Buildings) Init(w donburi.World) {
	b.board = archetype.MustFindBoard(w)

	b.roads = make([][]*entryWithPos, b.board.Size.Width)
	b.settlements = make([][]*entryWithPos, b.board.Size.Width)
	for i := 0; i < b.board.Size.Width; i++ {
		b.roads[i] = make([]*entryWithPos, b.board.Size.Height)
		b.settlements[i] = make([]*entryWithPos, b.board.Size.Height)
	}

	b.resources = component.Resources.Get(engine.MustFindWithComponent(w, component.Resources))

	events.NewRoadRequestedEvent.Subscribe(w, b.onRoadRequested)
	events.NewSettlementRequestedEvent.Subscribe(w, b.onSettlementRequested)
	events.SettlementUpgradeRequestedEvent.Subscribe(w, b.onSettlementUpgradeRequested)
	events.TechnologyResearchedEvent.Subscribe(w, b.onTechnologyResearched)
}

func (b *Buildings) Update(w donburi.World) {}

func (b *Buildings) onRoadRequested(w donburi.World, event events.NewRoadRequested) {
	if b.SpawnRoad(w, event.Tile, true) {
		showRoadPlaceholders(w, event.Tile)
	}
}

// SpawnRoad returns true if the road was spawned or already exists
func (b *Buildings) SpawnRoad(w donburi.World, tileEntry *donburi.Entry, requestedByPlayer bool) bool {
	tile := component.Tile.Get(tileEntry)

	if requestedByPlayer {
		// Idempotency
		if tile.HasRoad {
			return true
		}

		if !canBuildRoad(tile) {
			return false
		}
	}

	existingRoad := b.roads[tile.Position.X][tile.Position.Y]
	if existingRoad != nil {
		return true
	}

	neighbors := b.neighborsForRoad(tile.Position.X, tile.Position.Y)
	var roads []*entryWithPos
	if neighbors.N != nil {
		roads = append(roads, neighbors.N)
	}
	if neighbors.S != nil {
		roads = append(roads, neighbors.S)
	}
	if neighbors.W != nil {
		roads = append(roads, neighbors.W)
	}
	if neighbors.E != nil {
		roads = append(roads, neighbors.E)
	}

	// Can't build without another road nearby
	if requestedByPlayer && len(roads) == 0 {
		return false
	}

	if requestedByPlayer {
		cost := roadCost(w, tileEntry)
		if b.resources.Stone < cost {
			return false
		}

		UpdateResources(w, events.ResourcesDelta{
			Stone: -cost,
		})
	}

	newRoad := archetype.New(w).
		WithParent(tileEntry).
		With(component.Road).
		WithLayer(component.SpriteLayerRoad).
		WithSprite(component.SpriteData{
			Image: tile.RoadImage(),
		}).
		Entry()

	tile.HasRoad = true

	b.roads[tile.Position.X][tile.Position.Y] = &entryWithPos{
		pos: component.Position{
			X: tile.Position.X,
			Y: tile.Position.Y,
		},
		entry:  newRoad,
		parent: tileEntry,
	}

	updateTileBuildPermissions(w, tileEntry)
	neighborTiles := tile.NeighborTiles.All()
	for i := range neighborTiles {
		n := neighborTiles[i]
		updateTileBuildPermissions(w, n.Entry)
	}

	for i := range roads {
		rr := roads[i]
		component.Sprite.Get(rr.entry).Image = component.Tile.Get(rr.parent).RoadImage()
	}

	b.updateRoadsRanges(w)
	b.updateHarvestConditions(w, tileEntry)

	if requestedByPlayer {
		events.RoadSpawnedEvent.Publish(w, events.RoadSpawned{
			Road: newRoad,
			Tile: tileEntry,
		})
	}

	harvestTiles(w)

	return true
}

func UpdateAllTilesBuildPermissions(w donburi.World) {
	donburi.NewQuery(filter.Contains(component.Tile)).Each(w, func(entry *donburi.Entry) {
		updateTileBuildPermissions(w, entry)
	})
}

func updateTileBuildPermissions(w donburi.World, tileEntry *donburi.Entry) {
	tile := component.Tile.Get(tileEntry)

	for _, neighbor := range tile.NeighborTiles.All() {
		neighborTile := component.Tile.Get(neighbor.Entry)
		if neighborTile.HasRoad {
			tile.HasNeighborRoad = true
			break
		}
	}

	originSettlementSpawned := donburi.NewQuery(filter.Contains(component.Settlement)).Count(w) > 0

	tile.CanBuildBuildings = tile.Terrain == domain.TerrainPlains && tile.Deposit == nil && !tile.HasBuilding && (tile.HasRoad || !originSettlementSpawned)
	tile.CanBuildRoad = canBuildRoad(tile)
}

func canBuildRoad(tile *component.TileData) bool {
	// TODO water/mountains canBuildRoad should be defined in domain
	return !tile.HasRoad && tile.HasNeighborRoad && tile.Terrain != domain.TerrainWater && tile.Terrain != domain.TerrainMountains
}

func (b *Buildings) updateRoadsRanges(w donburi.World) {
	b.roadsQuery.Each(w, func(entry *donburi.Entry) {
		parent, ok := transform.GetParent(entry)
		if !ok {
			panic("road has no parent")
		}
		tile := component.Tile.Get(parent)
		forTilesInRange(w, tile.Position.X, tile.Position.Y, domain.BaseRoadSightRange, func(tile *donburi.Entry) {
			markTileInSight(parent, tile)
		})
	})
}

// TODO Perhaps a better fit for map system
func (b *Buildings) updateHarvestConditions(w donburi.World, tileEntry *donburi.Entry) {
	tile := component.Tile.Get(tileEntry)

	if tile.Resource.Amount > 0 && tile.InControlRange && tile.HasRoad && !tile.CanBeHarvested {
		tile.CanBeHarvested = true

		archetype.New(w).
			WithParent(tileEntry).
			WithLayer(component.SpriteLayerIndicator).
			WithSprite(component.SpriteData{
				Image: assets.FlagWhite,
			})
	}
}

type entryWithPos struct {
	pos    component.Position
	entry  *donburi.Entry
	parent *donburi.Entry
}

type neighborRoads struct {
	N *entryWithPos
	S *entryWithPos
	W *entryWithPos
	E *entryWithPos
}

func (b *Buildings) neighborsForRoad(x, y int) neighborRoads {
	n := neighborRoads{}

	if y > 0 {
		n.N = b.roads[x][y-1]
	}
	if y < len(b.roads[x])-1 {
		n.S = b.roads[x][y+1]
	}
	if x > 0 {
		n.W = b.roads[x-1][y]
	}
	if x < len(b.roads)-1 {
		n.E = b.roads[x+1][y]
	}

	return n
}

func SpawnSettlers(w donburi.World, tileEntry *donburi.Entry) {
	if donburi.NewQuery(filter.Contains(component.Settlers)).Count(w) > 0 {
		panic("already have settlers")
	}

	archetype.New(w).
		WithParent(tileEntry).
		WithLayer(component.SpriteLayerBuildings).
		WithSprite(component.SpriteData{
			Image: assets.Settlers,
		}).
		With(component.Settlers).
		Entry()

	RevealInitialSettlersTiles(w, tileEntry)
}

func RevealInitialSettlersTiles(w donburi.World, tileEntry *donburi.Entry) {
	player := engine.MustFindComponent[component.PlayerData](w, component.Player)
	sightRange := domain.SettlersSightRange
	if player.Character == domain.PlayerCharacterQueen {
		sightRange += domain.QueenBonusSightRange
	}

	tile := component.Tile.Get(tileEntry)
	forTilesInRange(w, tile.Position.X, tile.Position.Y, sightRange, func(tile *donburi.Entry) {
		markTileInSight(tileEntry, tile)
	})
}

func (b *Buildings) onSettlementRequested(w donburi.World, event events.NewSettlementRequested) {
	count := b.settlementsQuery.Count(w)
	if count == 0 {
		b.spawnOriginSettlement(w, event.Tile)
	} else {
		b.spawnSubsequentSettlement(w, event.Tile)
	}
}

func (b *Buildings) spawnOriginSettlement(w donburi.World, tileEntry *donburi.Entry) {
	sett := b.spawnSettlement(w, tileEntry)
	if sett == nil {
		return
	}

	tile := component.Tile.Get(tileEntry)
	x := tile.Position.X
	y := tile.Position.Y

	settlementNeighborTiles := []*donburi.Entry{}
	if x > 0 {
		settlementNeighborTiles = append(settlementNeighborTiles, b.board.Tiles[x-1][y])
	}
	if x < len(b.board.Tiles)-1 {
		settlementNeighborTiles = append(settlementNeighborTiles, b.board.Tiles[x+1][y])
	}
	if y > 0 {
		settlementNeighborTiles = append(settlementNeighborTiles, b.board.Tiles[x][y-1])
	}
	if y < len(b.board.Tiles[x])-1 {
		settlementNeighborTiles = append(settlementNeighborTiles, b.board.Tiles[x][y+1])
	}

	b.SpawnRoad(w, tileEntry, false)

	for i := range settlementNeighborTiles {
		neighborTile := settlementNeighborTiles[i]

		// TODO It would be better to have a list of allowed tiles
		if component.Tile.Get(neighborTile).Terrain == domain.TerrainWater {
			continue
		}

		b.SpawnRoad(w, neighborTile, false)
	}

	donburi.NewQuery(filter.Contains(component.Settlers)).Each(w, func(entry *donburi.Entry) {
		component.Destroy(entry)
	})

	player := engine.MustFindComponent[component.PlayerData](w, component.Player)
	if player.Character == domain.PlayerCharacterKing {
		b.upgradeSettlement(w, sett, false)
	}

	UpdateAllTilesBuildPermissions(w)

	events.OriginSettlementSpawnedEvent.Publish(w, events.OriginSettlementSpawned{
		Settlement: sett,
		Tile:       tileEntry,
	})
}

func (b *Buildings) spawnSubsequentSettlement(w donburi.World, tileEntry *donburi.Entry) {
	tile := component.Tile.Get(tileEntry)

	existingRoad := b.roads[tile.Position.X][tile.Position.Y]
	if existingRoad == nil {
		return
	}

	if b.resources.Wood < domain.SettlementCost {
		return
	}

	sett := b.spawnSettlement(w, tileEntry)
	if sett == nil {
		return
	}

	UpdateResources(w, events.ResourcesDelta{
		Wood: -domain.SettlementCost,
	})

	events.SettlementSpawnedEvent.Publish(w, events.SettlementSpawned{
		Settlement: sett,
		Tile:       tileEntry,
	})
}

func (b *Buildings) spawnSettlement(w donburi.World, tileEntry *donburi.Entry) *donburi.Entry {
	tile := component.Tile.Get(tileEntry)

	if !tile.InSightRange {
		return nil
	}

	if tile.Terrain != domain.TerrainPlains {
		return nil
	}

	existingSettlement := b.settlements[tile.Position.X][tile.Position.Y]
	if existingSettlement != nil {
		return nil
	}

	settlementsCount := b.settlementsQuery.Count(w)
	return b.CreateSettlementEntry(w, tileEntry, settlementsCount, 1)
}

func (b *Buildings) CreateSettlementEntry(
	w donburi.World,
	tileEntry *donburi.Entry,
	id int,
	level int,
) *donburi.Entry {
	sett := archetype.New(w).
		WithParent(tileEntry).
		WithLayer(component.SpriteLayerBuildings).
		WithSprite(component.SpriteData{
			Image: archetype.SettlementImageFromLevel(level),
		}).
		With(component.Settlement).
		Entry()

	archetype.New(w).
		WithParent(tileEntry).
		WithLayer(component.SpriteLayerIndicator).
		WithSprite(component.SpriteData{
			Image: assets.FlagWhite,
		})

	component.Settlement.SetValue(sett, component.SettlementData{
		ID:    id,
		Level: level,
	})

	tile := component.Tile.Get(tileEntry)

	b.settlements[tile.Position.X][tile.Position.Y] = &entryWithPos{
		pos: component.Position{
			X: tile.Position.X,
			Y: tile.Position.Y,
		},
		parent: tileEntry,
		entry:  sett,
	}

	tile.HasBuilding = true

	b.updateSettlementsRanges(w)
	updateTileBuildPermissions(w, tileEntry)

	return sett
}

func (b *Buildings) settlementExistsInDirection(start component.Position, dir component.Position) bool {
	if dir.X == 0 && dir.Y == 0 {
		panic("invalid direction")
	}

	x := start.X + dir.X
	y := start.Y + dir.Y
	for x >= 0 && x < b.board.Size.Width && y >= 0 && y < b.board.Size.Height {
		if b.roads[x][y] == nil {
			return false
		}

		if b.settlements[x][y] != nil {
			return true
		}

		x += dir.X
		y += dir.Y
	}

	return false
}

func (b *Buildings) updateSettlementsRanges(w donburi.World) {
	research := component.Research.Get(engine.MustFindWithComponent(w, component.Research))

	bonusSightRange := 0
	if research.KnownTechnologies[domain.TechnologySentryTower] {
		bonusSightRange = domain.SentrySightRangeBonus
	}

	bonusControlRange := 0

	b.settlementsQuery.Each(w, func(entry *donburi.Entry) {
		sett := component.Settlement.Get(entry)
		controlRange := domain.BaseSettlementControlRange + sett.Level + bonusControlRange
		sightRange := domain.BaseSettlementSightRange + sett.Level + bonusSightRange

		parent, ok := transform.GetParent(entry)
		if !ok {
			panic("settlement has no parent")
		}

		tile := component.Tile.Get(parent)

		forTilesInRange(w, tile.Position.X, tile.Position.Y, controlRange, func(tileEntry *donburi.Entry) {
			markTileInControl(tileEntry)
			b.updateHarvestConditions(w, tileEntry)
		})

		forTilesInRange(w, tile.Position.X, tile.Position.Y, sightRange, func(tile *donburi.Entry) {
			markTileInSight(parent, tile)
		})
	})

	harvestTiles(w)
	markBorders(w)
}

func forTilesInRange(w donburi.World, startX int, startY int, totalRange int, callback func(tile *donburi.Entry)) {
	board := engine.MustFindComponent[component.BoardData](w, component.Board)

	for i := -totalRange; i <= totalRange; i++ {
		effectiveRange := totalRange - engine.Abs(i)

		for j := -effectiveRange; j <= effectiveRange; j++ {
			x := startX + i
			y := startY + j

			if x >= 0 && x < board.Size.Width && y >= 0 && y < board.Size.Height {
				callback(board.Tiles[x][y])
			}
		}
	}
}

func markTileInControl(tileEntry *donburi.Entry) {
	tile := component.Tile.Get(tileEntry)
	if tile.InControlRange {
		return
	}

	tile.InControlRange = true
}

func markTileInSight(originTileEntry *donburi.Entry, tileEntry *donburi.Entry) {
	tile := component.Tile.Get(tileEntry)
	if tile.InSightRange {
		return
	}

	tile.InSightRange = true

	originTile := component.Tile.Get(originTileEntry)
	distance := engine.Abs(originTile.Position.X-tile.Position.X) + engine.Abs(originTile.Position.Y-tile.Position.Y)

	fog, ok := transform.FindChildWithComponent(tileEntry, component.TileFog)
	if !ok {
		panic("tile has no fog")
	}

	anim := component.Animation.Get(fog)
	anim.Timer = engine.NewTimer(fogRevealDuration * time.Duration(distance))
	anim.Start()
}

func (b *Buildings) onSettlementUpgradeRequested(w donburi.World, event events.SettlementUpgradeRequested) {
	b.upgradeSettlement(w, event.Settlement, true)
}

func (b *Buildings) upgradeSettlement(w donburi.World, settlement *donburi.Entry, requestedByPlayer bool) {
	sett := component.Settlement.Get(settlement)

	if sett.Level >= maxSettlementLevel(w) {
		return
	}

	if requestedByPlayer {
		if b.resources.Food < domain.SettlementUpgradeCost {
			return
		}

		UpdateResources(w, events.ResourcesDelta{
			Food: -domain.SettlementUpgradeCost,
		})
	}

	sett.Level++

	sprite := component.Sprite.Get(settlement)
	sprite.Image = archetype.SettlementImageFromLevel(sett.Level)

	b.updateSettlementsRanges(w)

	if sett.Level == domain.MaxSettlementLevelWithCurrency {
		UpdateResources(w, events.ResourcesDelta{
			Gold: domain.BonusGoldForMaxLevelSettlement,
		})
	}

	if requestedByPlayer {
		parent, ok := transform.GetParent(settlement)
		if !ok {
			panic("settlement has no parent")
		}
		events.SettlementUpgradedEvent.Publish(w, events.SettlementUpgraded{
			Settlement: settlement,
			Tile:       parent,
		})
	}
}

func (b *Buildings) onTechnologyResearched(w donburi.World, event events.TechnologyResearched) {
	if event.Technology == domain.TechnologySentryTower {
		b.updateSettlementsRanges(w)
	}
}

func UpdateResources(w donburi.World, delta events.ResourcesDelta,
) {
	resources := component.Resources.Get(engine.MustFindWithComponent(w, component.Resources))

	resources.Food += delta.Food
	resources.Stone += delta.Stone
	resources.Wood += delta.Wood
	resources.Iron += delta.Iron
	resources.Gold += delta.Gold
	resources.Diamonds += delta.Diamonds

	events.ResourcesUpdatedEvent.Publish(w, events.ResourcesUpdated{
		Resources: resources,
		Delta:     delta,
	})
}

func roadCost(w donburi.World, tileEntry *donburi.Entry) int {
	research := component.Research.Get(engine.MustFindWithComponent(w, component.Research))
	if research.KnownTechnologies[domain.TechnologyStonecraft] {
		tile := component.Tile.Get(tileEntry)
		if tile.InControlRange {
			return 0
		}
	}

	return domain.BaseRoadCost
}

func maxSettlementLevel(w donburi.World) int {
	research := component.Research.Get(engine.MustFindWithComponent(w, component.Research))
	if !research.KnownTechnologies[domain.TechnologyCurrency] {
		return domain.MaxSettlementLevelBase
	}

	return domain.MaxSettlementLevelWithCurrency
}
