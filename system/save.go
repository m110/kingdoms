package system

import (
	"fmt"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
	"github.com/m110/kingdoms/save"
)

const (
	saveVersion = 1
)

type Save struct {
	globalSave         *save.Global
	saveGameFunc       func(global *save.Global, game *save.Game)
	saveGameReplayFunc func(replay *save.GameReplay)
	backToMainMenuFunc func()
	saveEnabled        bool
}

func NewSave(
	loadedSaveData *save.Global,
	saveGlobalFunc func(data *save.Global, game *save.Game),
	saveGameReplayFunc func(replay *save.GameReplay),
	backToMainMenuFunc func(),
	saveEnabled bool,
) *Save {
	return &Save{
		globalSave:         loadedSaveData,
		saveGameFunc:       saveGlobalFunc,
		saveGameReplayFunc: saveGameReplayFunc,
		backToMainMenuFunc: backToMainMenuFunc,
		saveEnabled:        saveEnabled,
	}
}

func (s *Save) Init(w donburi.World) {
	events.QuitGameRequestedEvent.Subscribe(w, s.onQuitGameRequested)
	events.RoadSpawnedEvent.Subscribe(w, func(w donburi.World, event events.RoadSpawned) {
		s.onDayProgressed(w)
	})
	events.SettlementSpawnedEvent.Subscribe(w, func(w donburi.World, event events.SettlementSpawned) {
		s.onDayProgressed(w)
	})
	events.SettlementUpgradedEvent.Subscribe(w, func(w donburi.World, event events.SettlementUpgraded) {
		s.onDayProgressed(w)
	})
}

func (s *Save) onQuitGameRequested(w donburi.World, event events.QuitGameRequested) {
	if s.saveEnabled {
		s.save(w)
	}
	s.backToMainMenuFunc()
}

func (s *Save) save(w donburi.World) {
	s.saveGame(w)
	s.saveGameReplay(w)
}

func (s *Save) onDayProgressed(w donburi.World) {
	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)

	if progress.Day == domain.DaysPerYear {
		progress.Day = 1
		progress.Year++
	} else {
		progress.Day++
	}

	seasonChanged := false
	newSeason := domain.SeasonByDay(progress.Day)
	if newSeason != progress.Season {
		seasonChanged = true
		progress.Season = newSeason
	}

	events.ProgressUpdatedEvent.Publish(w, events.ProgressUpdated{
		Progress: progress,
	})

	if seasonChanged {
		events.SeasonChangedEvent.Publish(w, events.SeasonChanged{
			Progress: progress,
		})

		if s.saveEnabled {
			s.save(w)
		}
	}
}

func (s *Save) Update(w donburi.World) {}

func (s *Save) saveGame(w donburi.World) {
	game := &save.Game{}
	game.Header = &save.Header{
		Version: saveVersion,
	}

	resources := engine.MustFindComponent[component.ResourcesData](w, component.Resources)
	game.Resources = &save.Resources{
		Wood:  int64(resources.Wood),
		Stone: int64(resources.Stone),
		Food:  int64(resources.Food),
		Iron:  int64(resources.Iron),
		Gold:  int64(resources.Gold),
	}

	camera := transform.GetTransform(archetype.MustFindCamera(w))
	game.Camera = &save.Camera{
		Position: mapVectorToSave(camera.LocalPosition),
		Zoom:     mapVectorToSave(camera.LocalScale),
	}

	game.Board = newBoardForSave(w, false)

	// TODO support more buildings
	var buildings []*save.Building
	query.NewQuery(filter.Contains(component.Settlement)).Each(w, func(entry *donburi.Entry) {
		parent, ok := transform.GetParent(entry)
		if !ok {
			panic("settlement has no parent")
		}

		tile := component.Tile.Get(parent)
		settlement := component.Settlement.Get(entry)

		buildings = append(buildings, &save.Building{
			Id:       int64(settlement.ID),
			Position: mapPositionToSave(tile.Position),
			Type:     save.BuildingType_BuildingSettlement,
			Level:    int64(settlement.Level),
		})
	})
	game.Buildings = buildings

	var roads []*save.Road
	query.NewQuery(filter.Contains(component.Road)).Each(w, func(entry *donburi.Entry) {
		parent, ok := transform.GetParent(entry)
		if !ok {
			panic("settlement has no parent")
		}

		tile := component.Tile.Get(parent)

		roads = append(roads, &save.Road{
			Position: mapPositionToSave(tile.Position),
		})
	})
	game.Roads = roads

	research := engine.MustFindComponent[component.ResearchData](w, component.Research)

	var technologies []save.TechnologyType
	for tech := range research.KnownTechnologies {
		technologies = append(technologies, mapTechnologyToSave(tech))
	}

	game.Research = &save.Research{
		Technologies: technologies,
	}

	game.Player = newPlayerForSave(w)

	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)
	game.Progress = &save.Progress{
		Season: mapSeasonToSave(progress.Season),
		Year:   int64(progress.Year),
		Day:    int64(progress.Day),
	}

	s.globalSave.Diamonds = int64(resources.Diamonds)

	s.saveGameFunc(s.globalSave, game)
}

func (s *Save) saveGameReplay(w donburi.World) {
	actionsRecorder := engine.MustFindComponent[component.ActionsRecorderData](w, component.ActionsRecorder)

	replay := &save.GameReplay{
		Header: &save.Header{
			Version: saveVersion,
		},
		Board:    newBoardForSave(w, true),
		Player:   newPlayerForSave(w),
		LastTick: int64(actionsRecorder.CurrentTick),
	}

	for _, action := range actionsRecorder.Actions {
		replay.Actions = append(replay.Actions, &save.Action{
			Type:    action.ActionType,
			Tick:    int64(action.Tick),
			Payload: action.Payload,
		})
	}

	s.saveGameReplayFunc(replay)
}

func newBoardForSave(w donburi.World, maxAmount bool) *save.Board {
	var tiles []*save.Tile
	query.NewQuery(filter.Contains(component.Tile)).Each(w, func(entry *donburi.Entry) {
		tile := component.Tile.Get(entry)

		amount := int64(tile.Resource.Amount)
		if maxAmount {
			amount = int64(tile.Resource.MaxAmount)
		}

		saveTile := &save.Tile{
			Position: mapPositionToSave(tile.Position),
			Terrain:  mapTerrainToSave(tile.Terrain),
			Resource: &save.TileResource{
				Resource:  mapResourceToSave(tile.Resource.Resource),
				Amount:    amount,
				MaxAmount: int64(tile.Resource.MaxAmount),
			},
		}

		if tile.Deposit != nil {
			deposit := mapDepositToSave(*tile.Deposit)
			saveTile.Deposit = &deposit
		}

		tiles = append(tiles, saveTile)
	})
	board := engine.MustFindComponent[component.BoardData](w, component.Board)
	return &save.Board{
		Tiles:         tiles,
		Width:         int64(board.Size.Width),
		Height:        int64(board.Size.Height),
		StartPosition: mapPositionToSave(board.StartPosition),
	}
}

func newPlayerForSave(w donburi.World) *save.Player {
	player := engine.MustFindComponent[component.PlayerData](w, component.Player)
	return &save.Player{
		Character: mapCharacterToSave(player.Character),
	}
}

func mapTerrainToSave(terrain domain.Terrain) save.TerrainType {
	switch terrain {
	case domain.TerrainPlains:
		return save.TerrainType_TerrainPlains
	case domain.TerrainForest:
		return save.TerrainType_TerrainForest
	case domain.TerrainDesert:
		return save.TerrainType_TerrainDesert
	case domain.TerrainMountains:
		return save.TerrainType_TerrainMountains
	case domain.TerrainWater:
		return save.TerrainType_TerrainWater
	case domain.TerrainShore:
		return save.TerrainType_TerrainShore
	default:
		panic(fmt.Sprintf("unknown terrain: %v", terrain))
	}
}

func mapTerrainFromSave(terrain save.TerrainType) domain.Terrain {
	switch terrain {
	case save.TerrainType_TerrainPlains:
		return domain.TerrainPlains
	case save.TerrainType_TerrainForest:
		return domain.TerrainForest
	case save.TerrainType_TerrainDesert:
		return domain.TerrainDesert
	case save.TerrainType_TerrainMountains:
		return domain.TerrainMountains
	case save.TerrainType_TerrainWater:
		return domain.TerrainWater
	case save.TerrainType_TerrainShore:
		return domain.TerrainShore
	default:
		panic(fmt.Sprintf("unknown terrain: %v", terrain))
	}
}

func mapDepositToSave(deposit domain.Deposit) save.DepositType {
	switch deposit {
	case domain.DepositTrees:
		return save.DepositType_DepositTrees
	case domain.DepositRocks:
		return save.DepositType_DepositRocks
	case domain.DepositBerries:
		return save.DepositType_DepositBerries
	case domain.DepositMushrooms:
		return save.DepositType_DepositMushrooms
	case domain.DepositDeer:
		return save.DepositType_DepositDeer
	case domain.DepositIronOre:
		return save.DepositType_DepositIronOre
	default:
		panic(fmt.Sprintf("unknown deposit: %v", deposit))
	}
}

func mapDepositFromSave(deposit save.DepositType) domain.Deposit {
	switch deposit {
	case save.DepositType_DepositTrees:
		return domain.DepositTrees
	case save.DepositType_DepositRocks:
		return domain.DepositRocks
	case save.DepositType_DepositBerries:
		return domain.DepositBerries
	case save.DepositType_DepositMushrooms:
		return domain.DepositMushrooms
	case save.DepositType_DepositDeer:
		return domain.DepositDeer
	case save.DepositType_DepositIronOre:
		return domain.DepositIronOre
	default:
		panic(fmt.Sprintf("unknown deposit: %v", deposit))
	}
}

func mapResourceToSave(resource domain.Resource) save.ResourceType {
	switch resource {
	case domain.ResourceWood:
		return save.ResourceType_ResourceWood
	case domain.ResourceStone:
		return save.ResourceType_ResourceStone
	case domain.ResourceFood:
		return save.ResourceType_ResourceFood
	case domain.ResourceIron:
		return save.ResourceType_ResourceIron
	case domain.ResourceGold:
		return save.ResourceType_ResourceGold
	default:
		panic(fmt.Sprintf("unknown resource: %v", resource))
	}
}

func MapResourceFromSave(resource save.ResourceType) domain.Resource {
	switch resource {
	case save.ResourceType_ResourceWood:
		return domain.ResourceWood
	case save.ResourceType_ResourceStone:
		return domain.ResourceStone
	case save.ResourceType_ResourceFood:
		return domain.ResourceFood
	case save.ResourceType_ResourceIron:
		return domain.ResourceIron
	case save.ResourceType_ResourceGold:
		return domain.ResourceGold
	default:
		panic(fmt.Sprintf("unknown resource: %v", resource))
	}
}

func mapTechnologyToSave(tech domain.Technology) save.TechnologyType {
	switch tech {
	case domain.TechnologyMining:
		return save.TechnologyType_TechnologyMining
	case domain.TechnologyStonecraft:
		return save.TechnologyType_TechnologyStonecraft
	case domain.TechnologyIronworking:
		return save.TechnologyType_TechnologyIronworking
	case domain.TechnologySentryTower:
		return save.TechnologyType_TechnologySentryTower
	case domain.TechnologySwords:
		return save.TechnologyType_TechnologySwords
	case domain.TechnologyKnighthood:
		return save.TechnologyType_TechnologyKnighthood
	case domain.TechnologyGranary:
		return save.TechnologyType_TechnologyGranary
	case domain.TechnologyHunting:
		return save.TechnologyType_TechnologyHunting
	case domain.TechnologyFarming:
		return save.TechnologyType_TechnologyFarming
	case domain.TechnologyCurrency:
		return save.TechnologyType_TechnologyCurrency
	case domain.TechnologyDocks:
		return save.TechnologyType_TechnologyDocks
	case domain.TechnologyTrade:
		return save.TechnologyType_TechnologyTrade
	default:
		panic(fmt.Sprintf("unknown technology: %v", tech))
	}
}

func MapTechnologyFromSave(tech save.TechnologyType) domain.Technology {
	switch tech {
	case save.TechnologyType_TechnologyMining:
		return domain.TechnologyMining
	case save.TechnologyType_TechnologyStonecraft:
		return domain.TechnologyStonecraft
	case save.TechnologyType_TechnologyIronworking:
		return domain.TechnologyIronworking
	case save.TechnologyType_TechnologySentryTower:
		return domain.TechnologySentryTower
	case save.TechnologyType_TechnologySwords:
		return domain.TechnologySwords
	case save.TechnologyType_TechnologyKnighthood:
		return domain.TechnologyKnighthood
	case save.TechnologyType_TechnologyGranary:
		return domain.TechnologyGranary
	case save.TechnologyType_TechnologyHunting:
		return domain.TechnologyHunting
	case save.TechnologyType_TechnologyFarming:
		return domain.TechnologyFarming
	case save.TechnologyType_TechnologyCurrency:
		return domain.TechnologyCurrency
	case save.TechnologyType_TechnologyDocks:
		return domain.TechnologyDocks
	case save.TechnologyType_TechnologyTrade:
		return domain.TechnologyTrade
	default:
		panic(fmt.Sprintf("unknown technology: %v", tech))
	}
}

func mapSeasonToSave(season domain.Season) save.SeasonType {
	switch season {
	case domain.SeasonSpring:
		return save.SeasonType_SeasonSpring
	case domain.SeasonSummer:
		return save.SeasonType_SeasonSummer
	case domain.SeasonFall:
		return save.SeasonType_SeasonFall
	case domain.SeasonWinter:
		return save.SeasonType_SeasonWinter
	default:
		panic(fmt.Sprintf("unknown season: %v", season))
	}
}

func MapSeasonFromSave(season save.SeasonType) domain.Season {
	switch season {
	case save.SeasonType_SeasonSpring:
		return domain.SeasonSpring
	case save.SeasonType_SeasonSummer:
		return domain.SeasonSummer
	case save.SeasonType_SeasonFall:
		return domain.SeasonFall
	case save.SeasonType_SeasonWinter:
		return domain.SeasonWinter
	default:
		panic(fmt.Sprintf("unknown season: %v", season))
	}
}

func mapCharacterToSave(character domain.PlayerCharacterType) save.PlayerCharacterType {
	switch character {
	case domain.PlayerCharacterKing:
		return save.PlayerCharacterType_PlayerCharacterKing
	case domain.PlayerCharacterQueen:
		return save.PlayerCharacterType_PlayerCharacterQueen
	case domain.PlayerCharacterPrincess:
		return save.PlayerCharacterType_PlayerCharacterPrincess
	case domain.PlayerCharacterMerchant:
		return save.PlayerCharacterType_PlayerCharacterMerchant
	default:
		panic(fmt.Sprintf("unknown character: %v", character))
	}
}

func MapCharacterFromSave(character save.PlayerCharacterType) domain.PlayerCharacterType {
	switch character {
	case save.PlayerCharacterType_PlayerCharacterKing:
		return domain.PlayerCharacterKing
	case save.PlayerCharacterType_PlayerCharacterQueen:
		return domain.PlayerCharacterQueen
	case save.PlayerCharacterType_PlayerCharacterPrincess:
		return domain.PlayerCharacterPrincess
	case save.PlayerCharacterType_PlayerCharacterMerchant:
		return domain.PlayerCharacterMerchant
	default:
		panic(fmt.Sprintf("unknown character: %v", character))
	}
}

func mapPositionToSave(p component.Position) *save.Position {
	return &save.Position{
		X: int64(p.X),
		Y: int64(p.Y),
	}
}

func MapPositionFromSave(p *save.Position) component.Position {
	return component.Position{
		X: int(p.X),
		Y: int(p.Y),
	}
}

func mapVectorToSave(v math.Vec2) *save.Vector {
	return &save.Vector{
		X: float32(v.X),
		Y: float32(v.Y),
	}
}

func MapVectorFromSave(v *save.Vector) math.Vec2 {
	return math.Vec2{
		X: float64(v.X),
		Y: float64(v.Y),
	}
}
