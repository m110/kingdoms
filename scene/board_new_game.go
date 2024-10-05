package scene

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/save"
	"github.com/m110/kingdoms/scene/generator"
	"github.com/m110/kingdoms/system"
)

type newGameBoardCreator struct {
	context    Context
	tutorial   bool
	character  domain.PlayerCharacterType
	saveSlot   *int
	globalSave *save.Global
}

func newNewGameBoardCreator(
	context Context,
	tutorial bool,
	character domain.PlayerCharacterType,
	saveSlot *int,
	globalSave *save.Global,
) *newGameBoardCreator {
	return &newGameBoardCreator{
		context:    context,
		tutorial:   tutorial,
		character:  character,
		saveSlot:   saveSlot,
		globalSave: globalSave,
	}
}

func (c *newGameBoardCreator) ExtraSystems() []System {
	var systems []System

	if c.tutorial {
		systems = append(systems, system.NewTutorial())
	} else {
		systems = append(systems, system.NewActionsRecord())
	}

	return systems
}

func (c *newGameBoardCreator) CreateActions(w donburi.World) {
	if !c.tutorial {
		w.Create(component.ActionsRecorder)
	}
}

func (c *newGameBoardCreator) CreatePlayer(w donburi.World) {
	player := engine.MustFindComponent[component.PlayerData](w, component.Player)
	player.Character = c.character
}

func (c *newGameBoardCreator) CreateProgress(w donburi.World) {
	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)
	progress.Season = domain.SeasonSpring
	progress.Year = 1
	progress.Day = 1
}

func (c *newGameBoardCreator) CreateBoard(w donburi.World) {
	var assetsConfig assets.GeneratorConfig

	if c.tutorial {
		assetsConfig = assets.TutorialWorldGeneratorConfig
	} else {
		assetsConfig = assets.WorldGeneratorConfig
	}
	config := generator.ConfigFromAssets(assetsConfig)

	boardEntry := newBoardEntry(w, config.Width, config.Height)

	var tiles [][]*donburi.Entry

	chunks := system.GenerateChunks(w, boardEntry, config.Width, config.Height)
	gen := generator.NewWorldGenerator(config)
	gw := gen.Generate()
	tiles = system.CreateTiles(w, chunks, gw.Terrains, gw.Deposits, config.Width, config.Height, true, false)

	var startPosition component.Position
	if c.tutorial {
		startPosition.X = 42
		startPosition.Y = 30
	} else {
		var ok bool
		startPosition.X, startPosition.Y, ok = findInitPosition(tiles, config.Width, config.Height)
		if !ok {
			c.context.SceneSwitcher.SwitchToNewGame(*c.saveSlot, c.character)
			return
		}
	}

	board := component.Board.Get(boardEntry)
	board.Tiles = tiles
	board.StartPosition = startPosition
}

func (c *newGameBoardCreator) CreateResearch(w donburi.World) {
	research := engine.MustFindComponent[component.ResearchData](w, component.Research)
	research.KnownTechnologies = map[domain.Technology]bool{}
}

func (c *newGameBoardCreator) CreateResources(w donburi.World) {
	var initResources component.ResourcesData
	if c.tutorial {
		initResources = component.ResourcesData{
			Wood:  5,
			Stone: 10,
			Food:  0,
		}
	} else {
		initResources = initResourcesForNewGame(w, c.character)
		initResources.Diamonds = int(c.globalSave.Diamonds)
	}

	resources := engine.MustFindWithComponent(w, component.Resources)
	component.Resources.SetValue(resources, initResources)
}

func initResourcesForNewGame(w donburi.World, character domain.PlayerCharacterType) component.ResourcesData {
	initResources := component.ResourcesData{
		Wood:  20,
		Stone: 10,
		Food:  0,
	}

	if character == domain.PlayerCharacterPrincess {
		initResources.Food += 5
	} else if character == domain.PlayerCharacterMerchant {
		initResources.Wood += 10
		initResources.Stone += 10
	}

	return initResources
}

func (c *newGameBoardCreator) CreateUI(w donburi.World) {
	newResourcesPanel(w, !c.tutorial)
	archetype.NewSeasonUIPanel(w)
	newMenuPanel(w, c.context, !c.tutorial)
}

func (c *newGameBoardCreator) InitUnits(w donburi.World) {
	board := engine.MustFindComponent[component.BoardData](w, component.Board)
	tiles := board.Tiles

	settlementTile := tiles[board.StartPosition.X][board.StartPosition.Y]
	system.SpawnSettlers(w, settlementTile)
}

func (c *newGameBoardCreator) InitCamera(w donburi.World) {
	board := engine.MustFindComponent[component.BoardData](w, component.Board)
	tiles := board.Tiles

	settlementTile := tiles[board.StartPosition.X][board.StartPosition.Y]
	tilePos := transform.WorldPosition(settlementTile)

	camera := archetype.MustFindCamera(w)
	cameraScale := transform.GetTransform(camera).LocalScale
	transform.GetTransform(camera).LocalPosition = math.Vec2{
		X: tilePos.X - float64(c.context.ScreenWidth)/cameraScale.X/2,
		Y: tilePos.Y - float64(c.context.ScreenHeight)/cameraScale.Y/2,
	}
}

func (c *newGameBoardCreator) InitExtra(w donburi.World) {
	if !c.tutorial {
		archetype.ShowMessage(w, archetype.CharacterSettler, "{player}!\n\nWe have arrived at our new home.\n\nLet's build a settlement.")
	}
}

func findInitPosition(
	tiles [][]*donburi.Entry,
	mapWidth int,
	mapHeight int,
) (int, int, bool) {
	for i := 0; i < 50; i++ {
		settlementX := engine.RandomIntRange(mapWidth/3, mapWidth-mapWidth/3)
		settlementY := engine.RandomIntRange(mapHeight/3, mapHeight-mapHeight/3)

		tile := component.Tile.Get(tiles[settlementX][settlementY])
		if tile.Terrain == domain.TerrainPlains {
			return settlementX, settlementY, true
		}
	}

	return 0, 0, false
}
