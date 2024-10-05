package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	gameevents "github.com/m110/kingdoms/events"
	"github.com/m110/kingdoms/save"
	"github.com/m110/kingdoms/system"
)

type boardCreator interface {
	ExtraSystems() []System

	CreateActions(w donburi.World)
	CreatePlayer(w donburi.World)
	CreateProgress(w donburi.World)
	CreateBoard(w donburi.World)
	CreateResearch(w donburi.World)
	CreateResources(w donburi.World)
	CreateUI(w donburi.World)

	InitUnits(w donburi.World)
	InitCamera(w donburi.World)
	InitExtra(w donburi.World)
}

type Board struct {
	context   Context
	world     donburi.World
	systems   []System
	drawables []Drawable

	buildings *system.Buildings
	progress  *system.Progress

	saveSlot *int
	tutorial bool
}

func NewBoard(
	context Context,
	saveSlot *int,
	character domain.PlayerCharacterType,
	tutorial bool,
) *Board {
	globalSave, err := context.Storage.LoadGlobalData()
	if err != nil {
		// TODO Better error alerting
		panic(err)
	}

	buildings := system.NewBuildings()
	creator := newNewGameBoardCreator(context, tutorial, character, saveSlot, globalSave)
	return newBoard(context, saveSlot, globalSave, true, tutorial, creator, buildings)
}

func NewBoardFromSave(
	context Context,
	saveSlot int,
	game *save.Game,
	replay *save.GameReplay,
) *Board {
	globalSave, err := context.Storage.LoadGlobalData()
	if err != nil {
		// TODO Better error alerting
		panic(err)
	}

	// TODO hack fix asap :(
	buildings := system.NewBuildings()
	creator := newLoadGameBoardCreator(context, globalSave, game, replay, buildings)
	return newBoard(context, &saveSlot, globalSave, true, false, creator, buildings)
}

func NewBoardFromReplay(
	context Context,
	replaySlot int,
	replay *save.GameReplay,
) *Board {
	// TODO deduplicate across constructors
	globalSave, err := context.Storage.LoadGlobalData()
	if err != nil {
		// TODO Better error alerting
		panic(err)
	}

	// TODO hack fix asap :(
	buildings := system.NewBuildings()
	creator := newReplayGameBoardCreator(context, replaySlot, replay)
	return newBoard(context, nil, globalSave, false, false, creator, buildings)
}

func newBoard(
	context Context,
	saveSlot *int,
	saveData *save.Global,
	mapInteractive bool,
	tutorial bool,
	creator boardCreator,
	buildings *system.Buildings,
) *Board {
	progress := system.NewProgress()
	controls := system.NewControls(mapInteractive)

	// TODO probably could use the strategy pattern as well
	systems := creator.ExtraSystems()

	// A simplified save handler, so the save system doesn't need to concern
	// about the slot number or error handling
	saveFunc := func(global *save.Global, game *save.Game) {
		err := context.Storage.SaveSlot(*saveSlot, global, game)
		if err != nil {
			// TODO Better error alerting
			panic(err)
		}
	}

	saveReplayFunc := func(replay *save.GameReplay) {
		err := context.Storage.SaveReplay(*saveSlot, replay)
		if err != nil {
			// TODO Better error alerting
			panic(err)
		}
	}

	canSave := saveSlot != nil && !tutorial

	systems = append(systems, []System{
		system.NewDebug(),
		system.NewVelocity(),
		system.NewAnimation(),
		buildings,
		controls,
		system.NewUI(),
		system.NewAudio(),
		progress,
		system.NewResearch(),
		system.NewMap(),
		system.NewSave(
			saveData,
			saveFunc,
			saveReplayFunc,
			context.SceneSwitcher.SwitchToMainMenu,
			canSave,
		),
		system.NewSpriteTransition(),
		system.NewText(),
		system.NewTimeToLive(),
		system.NewDestroy(),
	}...)

	drawables := []Drawable{
		system.NewRenderer(),
		controls,
	}

	b := &Board{
		context:   context,
		systems:   systems,
		drawables: drawables,
		buildings: buildings,
		progress:  progress,
		saveSlot:  saveSlot,
		tutorial:  tutorial,
	}

	b.world = b.createWorld()

	creator.CreateActions(b.world)
	creator.CreatePlayer(b.world)
	creator.CreateProgress(b.world)
	creator.CreateBoard(b.world)
	creator.CreateResearch(b.world)
	creator.CreateResources(b.world)
	creator.CreateUI(b.world)

	// Init initializes all systems, so they subscribe to events
	// Any setup that publishes events should be done after this point
	b.init()

	creator.InitUnits(b.world)
	creator.InitCamera(b.world)
	creator.InitExtra(b.world)

	resources := component.Resources.Get(engine.MustFindWithComponent(b.world, component.Resources))
	gameevents.ResourcesUpdatedEvent.Publish(b.world, gameevents.ResourcesUpdated{
		Resources: resources,
		Delta: gameevents.ResourcesDelta{
			Wood:  resources.Wood,
			Stone: resources.Stone,
			Food:  resources.Food,
		},
	})

	return b
}

func (b *Board) createWorld() donburi.World {
	w := donburi.NewWorld()

	archetype.NewCamera(w, math.Vec2{}, standardCameraRange())

	game := w.Entry(w.Create(component.Game))
	donburi.SetValue(game, component.Game, component.GameData{
		Settings: component.Settings{
			ScreenWidth:  b.context.ScreenWidth,
			ScreenHeight: b.context.ScreenHeight,
		},
	})

	w.Create(component.Debug)
	w.Create(component.Progress)
	w.Create(component.Resources)
	w.Create(component.Research)
	w.Create(component.Player)

	return w
}

func (b *Board) init() {
	for _, s := range b.systems {
		s.Init(b.world)
	}
	for _, d := range b.drawables {
		d.Init(b.world)
	}
}

func (b *Board) Update() {
	for _, s := range b.systems {
		s.Update(b.world)
	}

	events.ProcessAllEvents(b.world)
}

func (b *Board) Draw(screen *ebiten.Image) {
	for _, s := range b.drawables {
		s.Draw(b.world, screen)
	}
}

func newBoardEntry(w donburi.World, width int, height int) *donburi.Entry {
	board := archetype.New(w).
		With(component.Board).
		Entry()

	component.Board.SetValue(board, component.BoardData{
		Size: component.Size{
			Width:  width,
			Height: height,
		},
		TileSize: component.Size{
			Width:  32,
			Height: 32,
		},
	})

	return board
}

func newResourcesPanel(w donburi.World, showDiamonds bool) *donburi.Entry {
	panel := archetype.NewFrameWithBorder(w, math.Vec2{X: 96, Y: -2}, 700, 40)
	component.Layer.Get(panel).Layer = component.SpriteUILayerUI
	panel.AddComponent(component.UIPanel)

	posY := 8.0

	foodIcon := archetype.NewResourceIconWithValue(w, domain.ResourceFood, 0, math.Vec2{X: 10, Y: posY})
	foodIcon.AddComponent(component.IconFood)
	stoneIcon := archetype.NewResourceIconWithValue(w, domain.ResourceStone, 0, math.Vec2{X: 110, Y: posY})
	stoneIcon.AddComponent(component.IconStone)
	woodIcon := archetype.NewResourceIconWithValue(w, domain.ResourceWood, 0, math.Vec2{X: 210, Y: posY})
	woodIcon.AddComponent(component.IconWood)

	ironIcon := archetype.NewResourceIconWithValue(w, domain.ResourceIron, 0, math.Vec2{X: 310, Y: posY})
	ironIcon.AddComponent(component.IconIron)
	ironIcon.AddComponent(component.Active)
	component.Active.Get(ironIcon).Active = false

	goldIcon := archetype.NewResourceIconWithValue(w, domain.ResourceGold, 0, math.Vec2{X: 410, Y: posY})
	goldIcon.AddComponent(component.IconGold)
	goldIcon.AddComponent(component.Active)
	component.Active.Get(goldIcon).Active = false

	transform.AppendChild(panel, foodIcon, false)
	transform.AppendChild(panel, woodIcon, false)
	transform.AppendChild(panel, stoneIcon, false)
	transform.AppendChild(panel, ironIcon, false)
	transform.AppendChild(panel, goldIcon, false)

	if showDiamonds {
		diamondsIcon := archetype.NewResourceIconWithValue(w, domain.ResourceDiamonds, 0, math.Vec2{X: 510, Y: posY})
		diamondsIcon.AddComponent(component.IconDiamonds)
		transform.AppendChild(panel, diamondsIcon, false)
	}

	return panel
}

func newMenuPanel(w donburi.World, ctx Context, canSave bool) *donburi.Entry {
	menuPanel := archetype.NewFrameWithBorder(w, math.Vec2{X: -2, Y: 150}, 80, 288)
	menuPanel.AddComponent(component.UIPanel)
	component.Layer.Get(menuPanel).Layer = component.SpriteUILayerUI

	researchButton := archetype.NewIconButton(w, assets.IconResearch, math.Vec2{X: 10, Y: 12}, false, func(w donburi.World, e *donburi.Entry) {
		existing, ok := donburi.NewQuery(filter.Contains(component.ResearchPanel)).First(w)
		if ok {
			component.Destroy(existing)
			return
		}

		archetype.NewResearchPanel(w, nil)
	})
	transform.AppendChild(menuPanel, researchButton, false)

	pauseMenu := newPauseMenu(w, ctx, canSave)

	menuButton := archetype.NewIconButton(w, assets.IconMenu, math.Vec2{X: 10, Y: 216}, false, func(w donburi.World, e *donburi.Entry) {
		active := component.Active.Get(pauseMenu)
		active.Active = true

		game := component.MustFindGame(w)
		game.Paused = true
	})
	transform.AppendChild(menuPanel, menuButton, false)

	return menuPanel
}

func newPauseMenu(w donburi.World, ctx Context, canSave bool) *donburi.Entry {
	pauseImage := ebiten.NewImage(ctx.ScreenWidth, ctx.ScreenHeight)
	pauseImage.Fill(assets.PauseScreenColor)
	pauseScreen := archetype.NewFrameWithBorder(w, math.Vec2{X: 0, Y: 0}, ctx.ScreenWidth, ctx.ScreenHeight)
	pauseScreen.AddComponent(component.UIPanel)
	pauseScreen.AddComponent(component.Active)
	component.Active.Get(pauseScreen).Active = false
	component.Sprite.Get(pauseScreen).Image = pauseImage
	component.Layer.Get(pauseScreen).Layer = component.SpriteUILayerPauseOverlay

	cheatsPanel := archetype.NewFrameWithBorder(w, math.Vec2{X: 100, Y: 100}, 300, 400)
	cheatsPanel.AddComponent(component.UIPanel)
	cheatsPanel.AddComponent(component.Active)
	component.Layer.Get(cheatsPanel).Layer = component.SpriteUILayerPauseDialog
	component.Active.Get(cheatsPanel).Active = false

	foodIcon := archetype.NewIconButton(w, assets.IconFood, math.Vec2{X: 10, Y: 12}, true, func(w donburi.World, e *donburi.Entry) {
		system.UpdateResources(w, gameevents.ResourcesDelta{
			Food: 5,
		})
	})
	component.Layer.Get(foodIcon).Layer = component.SpriteUILayerPauseDialog
	transform.AppendChild(cheatsPanel, foodIcon, false)

	stoneIcon := archetype.NewIconButton(w, assets.IconStone, math.Vec2{X: 10, Y: 112}, true, func(w donburi.World, e *donburi.Entry) {
		system.UpdateResources(w, gameevents.ResourcesDelta{
			Stone: 5,
		})
	})
	component.Layer.Get(stoneIcon).Layer = component.SpriteUILayerPauseDialog
	transform.AppendChild(cheatsPanel, stoneIcon, false)

	woodIcon := archetype.NewIconButton(w, assets.IconWood, math.Vec2{X: 10, Y: 212}, true, func(w donburi.World, e *donburi.Entry) {
		system.UpdateResources(w, gameevents.ResourcesDelta{
			Wood: 5,
		})
	})
	component.Layer.Get(woodIcon).Layer = component.SpriteUILayerPauseDialog
	transform.AppendChild(cheatsPanel, woodIcon, false)

	diamondIcon := archetype.NewIconButton(w, assets.IconDiamond, math.Vec2{X: 110, Y: 12}, true, func(w donburi.World, e *donburi.Entry) {
		system.UpdateResources(w, gameevents.ResourcesDelta{
			Diamonds: 5,
		})
	})
	component.Layer.Get(diamondIcon).Layer = component.SpriteUILayerPauseDialog
	transform.AppendChild(cheatsPanel, diamondIcon, false)

	toggleFogButton := archetype.NewButton(w, "Toggle Fog", math.Vec2{X: 10, Y: 312}, true, func(w donburi.World, e *donburi.Entry) {
		donburi.NewQuery(filter.Contains(component.TileFog)).Each(w, func(e *donburi.Entry) {
			sprite := component.Sprite.Get(e)
			sprite.Hidden = !sprite.Hidden
		})
	})
	component.Layer.Get(toggleFogButton).Layer = component.SpriteUILayerPauseDialog
	transform.AppendChild(cheatsPanel, toggleFogButton, false)

	pauseMenuWidth := 384
	buttonMargin := 16.0
	buttonWidth := pauseMenuWidth - 64

	pauseMenu := archetype.NewFrameWithBorder(w, math.Vec2{X: 420, Y: 130}, pauseMenuWidth, 256)
	component.Layer.Get(pauseMenu).Layer = component.SpriteUILayerPauseDialog

	archetype.New(w).
		WithParent(pauseMenu).
		WithPosition(math.Vec2{X: 150, Y: 32}).
		WithText(component.TextData{
			Text: "Pause",
		})

	continueButton := archetype.NewButtonWithWidth(w, "Continue", math.Vec2{X: buttonMargin, Y: 48}, false, func(w donburi.World, e *donburi.Entry) {
		component.Active.Get(pauseScreen).Active = false
		game := component.MustFindGame(w)
		game.Paused = false
	}, buttonWidth)
	transform.AppendChild(pauseMenu, continueButton, false)

	debugButton := archetype.NewButtonWithWidth(w, "Toggle Debug", math.Vec2{X: buttonMargin, Y: 128}, false, func(w donburi.World, e *donburi.Entry) {
		debug := engine.MustFindComponent[component.DebugData](w, component.Debug)
		debug.Enabled = !debug.Enabled

		component.Active.Get(cheatsPanel).Active = debug.Enabled
	}, buttonWidth)
	transform.AppendChild(pauseMenu, debugButton, false)

	quitText := "Quit"
	if canSave {
		quitText = "Save & Quit"
	}

	quitButton := archetype.NewButtonWithWidth(w, quitText, math.Vec2{X: buttonMargin, Y: 192}, false, func(w donburi.World, e *donburi.Entry) {
		gameevents.QuitGameRequestedEvent.Publish(w, gameevents.QuitGameRequested{})
	}, buttonWidth)
	transform.AppendChild(pauseMenu, quitButton, false)

	transform.AppendChild(pauseScreen, pauseMenu, false)

	return pauseScreen
}

func standardCameraRange() engine.FloatRange {
	if system.UseTouchControls() {
		return engine.FloatRange{Min: 0.5, Max: 1.5}
	} else {
		return engine.FloatRange{Min: 0.5, Max: 1}
	}
}
