package scene

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"

	"github.com/m110/kingdoms/domain"

	"github.com/m110/kingdoms/archetype"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/save"
	"github.com/m110/kingdoms/system"
)

type loadGameBoardCreator struct {
	context    Context
	globalSave *save.Global
	game       *save.Game
	replay     *save.GameReplay

	// TODO This should not be the case
	// Rework the system so it has no "roads" and "settlements" cache
	// Then, there can be simple functions that create roads and settlements
	buildings *system.Buildings
}

func newLoadGameBoardCreator(
	context Context,
	globalSave *save.Global,
	game *save.Game,
	replay *save.GameReplay,
	buildings *system.Buildings,
) *loadGameBoardCreator {
	return &loadGameBoardCreator{
		context:    context,
		globalSave: globalSave,
		game:       game,
		replay:     replay,
		buildings:  buildings,
	}
}

func (c *loadGameBoardCreator) ExtraSystems() []System {
	systems := []System{
		system.NewActionsRecord(),
	}

	return systems
}

func (c *loadGameBoardCreator) CreateActions(w donburi.World) {
	w.Create(component.ActionsRecorder)

	recorder := engine.MustFindComponent[component.ActionsRecorderData](w, component.ActionsRecorder)

	var actions []component.Action
	for _, a := range c.replay.Actions {
		actions = append(actions, component.Action{
			Tick:       int(a.Tick),
			ActionType: a.Type,
			Payload:    a.Payload,
		})
	}

	recorder.Load(int(c.replay.LastTick), actions)
}

func (c *loadGameBoardCreator) CreatePlayer(w donburi.World) {
	character := system.MapCharacterFromSave(c.game.Player.Character)
	player := engine.MustFindComponent[component.PlayerData](w, component.Player)
	player.Character = character
}

func (c *loadGameBoardCreator) CreateProgress(w donburi.World) {
	season := system.MapSeasonFromSave(c.game.Progress.Season)
	progress := engine.MustFindComponent[component.ProgressData](w, component.Progress)
	progress.Season = season
	progress.Year = int(c.game.Progress.Year)
	progress.Day = int(c.game.Progress.Day)
}

func (c *loadGameBoardCreator) CreateBoard(w donburi.World) {
	loadBoard(w, c.context, c.game.Board)
}

func loadBoard(w donburi.World, context Context, board *save.Board) {
	boardWidth := int(board.Width)
	boardHeight := int(board.Height)

	boardEntry := newBoardEntry(w, boardWidth, boardHeight)

	var tiles [][]*donburi.Entry

	chunks := system.GenerateChunks(w, boardEntry, boardWidth, boardHeight)
	tiles = system.LoadTiles(w, board, chunks, true, false)

	component.Board.Get(boardEntry).Tiles = tiles
}

func (c *loadGameBoardCreator) CreateResearch(w donburi.World) {
	research := engine.MustFindComponent[component.ResearchData](w, component.Research)
	research.KnownTechnologies = map[domain.Technology]bool{}

	for _, t := range c.game.Research.Technologies {
		technology := system.MapTechnologyFromSave(t)
		research.KnownTechnologies[technology] = true
	}
}

func (c *loadGameBoardCreator) CreateResources(w donburi.World) {
	resources := engine.MustFindWithComponent(w, component.Resources)

	component.Resources.SetValue(resources, component.ResourcesData{
		Food:     int(c.game.Resources.Food),
		Stone:    int(c.game.Resources.Stone),
		Wood:     int(c.game.Resources.Wood),
		Iron:     int(c.game.Resources.Iron),
		Gold:     int(c.game.Resources.Gold),
		Diamonds: int(c.globalSave.Diamonds),
	})
}

func (c *loadGameBoardCreator) CreateUI(w donburi.World) {
	newResourcesPanel(w, true)
	archetype.NewSeasonUIPanel(w)
	newMenuPanel(w, c.context, true)
}

func (c *loadGameBoardCreator) InitUnits(w donburi.World) {
	board := engine.MustFindComponent[component.BoardData](w, component.Board)
	board.StartPosition = system.MapPositionFromSave(c.game.Board.StartPosition)

	tiles := board.Tiles
	settlementTile := tiles[c.game.Board.StartPosition.X][c.game.Board.StartPosition.Y]

	if len(c.game.Buildings) == 0 {
		system.SpawnSettlers(w, settlementTile)
	} else {
		// Don't show settlers, but still reveal the initial tiles
		system.RevealInitialSettlersTiles(w, settlementTile)
	}

	for _, b := range c.game.Buildings {
		if b.Type == save.BuildingType_BuildingSettlement {
			tile := tiles[b.Position.X][b.Position.Y]
			c.buildings.CreateSettlementEntry(w, tile, int(b.Id), int(b.Level))
		}
	}

	for _, u := range c.game.Roads {
		tile := tiles[u.Position.X][u.Position.Y]
		c.buildings.SpawnRoad(w, tile, false)
	}

	system.UpdateAllTilesBuildPermissions(w)
}

func (c *loadGameBoardCreator) InitCamera(w donburi.World) {
	camera := archetype.MustFindCamera(w)
	cameraPosition := math.Vec2{
		X: float64(c.game.Camera.Position.X),
		Y: float64(c.game.Camera.Position.Y),
	}
	cameraScale := math.Vec2{
		X: float64(c.game.Camera.Zoom.X),
		Y: float64(c.game.Camera.Zoom.Y),
	}
	transform.GetTransform(camera).LocalPosition = cameraPosition
	transform.GetTransform(camera).LocalScale = cameraScale
}

func (c *loadGameBoardCreator) InitExtra(w donburi.World) {}
