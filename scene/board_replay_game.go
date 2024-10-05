package scene

import (
	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/save"
	"github.com/m110/kingdoms/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"time"
)

type replayGameBoardCreator struct {
	context    Context
	replaySlot int
	replay     *save.GameReplay

	character domain.PlayerCharacterType

	newGameCreator *newGameBoardCreator
}

func newReplayGameBoardCreator(
	context Context,
	replaySlot int,
	replay *save.GameReplay,
) *replayGameBoardCreator {
	character := system.MapCharacterFromSave(replay.Player.Character)

	return &replayGameBoardCreator{
		context:        context,
		replaySlot:     replaySlot,
		replay:         replay,
		character:      character,
		newGameCreator: newNewGameBoardCreator(context, false, character, nil, nil),
	}
}

func (c *replayGameBoardCreator) ExtraSystems() []System {
	systems := []System{
		system.NewActionsPlay(),
	}

	return systems
}

func (c *replayGameBoardCreator) CreateActions(w donburi.World) {
	game := component.MustFindGame(w)
	game.ActionsDisabled = true

	w.Create(component.ActionsPlayer)

	player := engine.MustFindComponent[component.ActionsPlayerData](w, component.ActionsPlayer)

	var actions []component.Action
	for _, a := range c.replay.Actions {
		actions = append(actions, component.Action{
			Tick:       int(a.Tick),
			ActionType: a.Type,
			Payload:    a.Payload,
		})
	}

	player.Load(actions)

	archetype.NewReplayPanel(w, func() {
		c.context.SceneSwitcher.SwitchToReplayedGame(c.replaySlot)
	})
}

func (c *replayGameBoardCreator) CreatePlayer(w donburi.World) {
	c.newGameCreator.CreatePlayer(w)
}

func (c *replayGameBoardCreator) CreateProgress(w donburi.World) {
	c.newGameCreator.CreateProgress(w)
}

func (c *replayGameBoardCreator) CreateBoard(w donburi.World) {
	loadBoard(w, c.context, c.replay.Board)
}

func (c *replayGameBoardCreator) CreateResearch(w donburi.World) {
	c.newGameCreator.CreateResearch(w)
}

func (c *replayGameBoardCreator) CreateResources(w donburi.World) {
	initResources := initResourcesForNewGame(w, c.character)
	resources := engine.MustFindWithComponent(w, component.Resources)
	component.Resources.SetValue(resources, initResources)
}

func (c *replayGameBoardCreator) CreateUI(w donburi.World) {
	newResourcesPanel(w, false)
	archetype.NewSeasonUIPanel(w)
	newMenuPanel(w, c.context, false)
}

func (c *replayGameBoardCreator) InitUnits(w donburi.World) {
	board := engine.MustFindComponent[component.BoardData](w, component.Board)
	board.StartPosition = system.MapPositionFromSave(c.replay.Board.StartPosition)
	c.newGameCreator.InitUnits(w)
}

func (c *replayGameBoardCreator) InitCamera(w donburi.World) {
	c.newGameCreator.InitCamera(w)
}

func (c *replayGameBoardCreator) InitExtra(w donburi.World) {
	textParent := archetype.New(w).
		WithPosition(math.Vec2{X: 620, Y: 100}).
		With(component.UI).
		Entry()

	replayText := archetype.New(w).
		WithParent(textParent).
		WithText(component.TextData{
			Text: "Replay",
			Size: component.TextSizeL,
		}).
		With(component.Animation).
		Entry()

	anim := component.Animation.Get(replayText)
	anim.Timer = engine.NewTimer(time.Second)
	anim.Start()
	anim.Update = func(e *donburi.Entry) {
		anim := component.Animation.Get(e)
		if anim.Timer.IsReady() {
			anim.Timer.Reset()
			text := component.Text.Get(e)
			text.Hidden = !text.Hidden
		}
	}
}
