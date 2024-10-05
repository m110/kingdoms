package scene

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	donburievents "github.com/yohamta/donburi/features/events"
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

const saveSlots = 3

type MenuItem struct {
	Text   string
	Action func()
}

type MainMenu struct {
	context        Context
	world          donburi.World
	systems        []System
	drawables      []Drawable
	globalSaveData *save.Global
}

func NewMainMenu(context Context) *MainMenu {
	controls := system.NewControls(false)

	systems := []System{
		system.NewDebug(),
		controls,
		system.NewAnimation(),
		system.NewAudio(),
		system.NewWander(),
		system.NewVelocity(),
		system.NewText(),
		system.NewTimeToLive(),
		system.NewDestroy(),
	}

	drawables := []Drawable{
		system.NewRenderer(),
		controls,
	}

	globalSaveData, err := context.Storage.LoadGlobalData()
	if err != nil {
		// TODO: better error alerting
		panic(err)
	}

	m := &MainMenu{
		context:        context,
		systems:        systems,
		drawables:      drawables,
		globalSaveData: globalSaveData,
	}

	m.world = m.createWorld()
	m.init()

	return m
}

func (m *MainMenu) createWorld() donburi.World {
	w := donburi.NewWorld()
	archetype.NewCamera(w, math.Vec2{}, standardCameraRange())

	game := w.Entry(w.Create(component.Game))
	donburi.SetValue(game, component.Game, component.GameData{
		Settings: component.Settings{
			ScreenWidth:  m.context.ScreenWidth,
			ScreenHeight: m.context.ScreenHeight,
		},
	})

	progress := w.Entry(w.Create(component.Progress))
	component.Progress.SetValue(progress, component.ProgressData{
		Season: domain.RandomSeason(),
		Day:    1,
		Year:   1,
	})

	mapSize := 100

	w.Create(component.Debug)

	camera := archetype.MustFindCamera(w)

	board := newBoardEntry(w, mapSize, mapSize)

	chunks := system.GenerateChunks(w, board, mapSize, mapSize)

	config := generator.ConfigFromAssets(assets.WorldGeneratorConfig)
	gen := generator.NewWorldGenerator(config)
	wg := gen.Generate()
	tiles := system.CreateTiles(w, chunks, wg.Terrains, wg.Deposits, mapSize, mapSize, false, true)

	component.Board.Get(board).Tiles = tiles

	centerTile := tiles[mapSize/2][mapSize/2]
	tilePos := transform.WorldPosition(centerTile)

	cameraPos := math.Vec2{
		X: tilePos.X - float64(m.context.ScreenWidth/4),
		Y: tilePos.Y - float64(m.context.ScreenHeight/4),
	}

	transform.GetTransform(camera).LocalPosition = cameraPos

	camera.AddComponent(component.Wander)
	camera.AddComponent(component.Velocity)
	component.Wander.SetValue(camera, component.WanderData{
		Min: cameraPos,
		Max: cameraPos.MulScalar(1.5),
	})
	component.Velocity.SetValue(camera, component.VelocityData{
		Velocity: math.Vec2{
			X: 0.5,
			Y: 0.5,
		},
	})

	dialog := archetype.NewFrameWithBorder(w, math.Vec2{X: 400, Y: 230}, 300, 100)

	archetype.New(w).
		WithParent(dialog).
		WithPosition(math.Vec2{X: 100, Y: 40}).
		WithText(component.TextData{
			Text: "Kingdoms",
		})

	archetype.New(w).
		WithParent(dialog).
		WithPosition(math.Vec2{X: 50, Y: 80}).
		WithText(component.TextData{
			Text:           "by Milosz Smolka",
			Streaming:      true,
			StreamingTimer: engine.NewTimer(1 * time.Second),
		})

	newGameMenu := newNewGameMenu(w, m.context)
	loadGameMenu := newLoadGameMenu(w, m.context)
	sidePanel := newSidePanel(w, m.context.ScreenWidth, m.context.ScreenHeight)

	diamondsIcon := archetype.NewResourceIconWithValue(w, domain.ResourceDiamonds, int(m.globalSaveData.Diamonds), math.Vec2{X: 10, Y: 60})
	diamondsIcon.AddComponent(component.IconGold)
	transform.AppendChild(sidePanel, diamondsIcon, false)

	tutorialButton := archetype.NewButton(w, "Tutorial", math.Vec2{X: 10, Y: 180}, false, func(w donburi.World, e *donburi.Entry) {
		m.context.SceneSwitcher.SwitchToTutorial()
	})
	newGameButton := archetype.NewButton(w, "New Game", math.Vec2{X: 10, Y: 250}, false, func(w donburi.World, e *donburi.Entry) {
		loadGameMenu.Hide()
		newGameMenu.Show()
	})
	loadGameButton := archetype.NewButton(w, "Load Game", math.Vec2{X: 10, Y: 320}, false, func(w donburi.World, e *donburi.Entry) {
		newGameMenu.Hide()
		loadGameMenu.Show()
	})

	musicOn := false
	musicButton := archetype.NewButton(w, "Toggle music", math.Vec2{X: 10, Y: 500}, false, func(w donburi.World, e *donburi.Entry) {
		musicOn = !musicOn
		m.context.SetMusic(musicOn)
	})

	archetype.New(w).
		WithParent(sidePanel).
		WithPosition(math.Vec2{X: 10, Y: float64(m.context.ScreenHeight - 10)}).
		WithText(component.TextData{
			Text: "Version: " + m.context.Version,
			Size: component.TextSizeS,
		})

	transform.AppendChild(sidePanel, tutorialButton, false)
	transform.AppendChild(sidePanel, newGameButton, false)
	transform.AppendChild(sidePanel, loadGameButton, false)
	transform.AppendChild(sidePanel, musicButton, false)

	return w
}

func (m *MainMenu) init() {
	for _, s := range m.systems {
		s.Init(m.world)
	}
	for _, d := range m.drawables {
		d.Init(m.world)
	}
}

func (m *MainMenu) Update() {
	for _, s := range m.systems {
		s.Update(m.world)
	}

	donburievents.ProcessAllEvents(m.world)
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	for _, s := range m.drawables {
		s.Draw(m.world, screen)
	}
}

type newGameMenu struct {
	world   donburi.World
	context Context

	saveSlot             *int
	saveSlotSelectDialog *donburi.Entry

	character             domain.PlayerCharacterType
	characterSelectDialog *donburi.Entry
	selectCharacter       func()
}

func newNewGameMenu(w donburi.World, context Context) *newGameMenu {
	n := &newGameMenu{
		world:   w,
		context: context,
	}

	return n
}

func (n *newGameMenu) showCharacterSelect() {
	if n.characterSelectDialog != nil {
		return
	}

	characterSelectDialog := newMainMenuDialog(n.world)

	icon := archetype.New(n.world).
		WithParent(characterSelectDialog).
		WithPosition(math.Vec2{X: 40, Y: 50}).
		WithSprite(component.SpriteData{}).
		Entry()

	title := archetype.New(n.world).
		WithParent(icon).
		WithPosition(math.Vec2{X: 0, Y: 160}).
		WithText(component.TextData{}).
		Entry()

	description := archetype.New(n.world).
		WithPosition(math.Vec2{X: 180, Y: 80}).
		WithParent(characterSelectDialog).
		WithText(component.TextData{
			Size: component.TextSizeM,
		}).
		Entry()

	n.selectCharacter = func() {
		character := domain.PlayerCharacters[n.character]

		iconImage := ebiten.NewImage(128, 128)
		iconImage.Fill(assets.CharacterFrameBackgroundColor)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(8, 8)

		iconImage.DrawImage(character.Image, op)
		iconImage.DrawImage(assets.CharacterFrame, op)

		component.Sprite.Get(icon).Image = iconImage
		component.Text.Get(title).Text = character.Name
		component.Text.Get(description).Text = character.Description
	}

	prevButton := archetype.NewButton(n.world, "  <-  ", math.Vec2{X: 190, Y: 220}, false, func(w donburi.World, e *donburi.Entry) {
		n.character--
		if n.character < 0 {
			n.character = domain.PlayerCharacterType(len(domain.PlayerCharacters) - 1)
		}
		n.selectCharacter()
	})
	transform.AppendChild(characterSelectDialog, prevButton, false)

	nextButton := archetype.NewButton(n.world, "  ->  ", math.Vec2{X: 340, Y: 220}, false, func(w donburi.World, e *donburi.Entry) {
		n.character++
		if n.character >= domain.PlayerCharacterType(len(domain.PlayerCharacters)) {
			n.character = 0
		}
		n.selectCharacter()
	})
	transform.AppendChild(characterSelectDialog, nextButton, false)

	startButton := archetype.NewButton(n.world, "Start", math.Vec2{X: 250, Y: 330}, false, func(w donburi.World, e *donburi.Entry) {
		n.context.SceneSwitcher.SwitchToNewGame(*n.saveSlot, n.character)
	})
	transform.AppendChild(characterSelectDialog, startButton, false)

	closeButton := archetype.NewButton(n.world, "X", math.Vec2{X: 550, Y: 10}, false, func(w donburi.World, e *donburi.Entry) {
		n.Hide()
	})
	transform.AppendChild(characterSelectDialog, closeButton, false)

	n.selectCharacter()

	n.characterSelectDialog = characterSelectDialog
}

func (n *newGameMenu) hideCharacterSelect() {
	if n.characterSelectDialog == nil {
		return
	}

	component.Destroy(n.characterSelectDialog)
	n.characterSelectDialog = nil
}

func (n *newGameMenu) showSaveSlotSelect() {
	if n.saveSlotSelectDialog != nil {
		return
	}

	saveSlotSelectDialog := newMainMenuDialog(n.world)

	archetype.New(n.world).
		WithParent(saveSlotSelectDialog).
		WithPosition(math.Vec2{X: 40, Y: 50}).
		WithText(component.TextData{
			Text: "Select save slot",
		})

	setSaveSlot := func(slot int) {
		n.saveSlot = &slot
		n.hideSaveSlotSelect()
		n.showCharacterSelect()
	}

	buttons := newSaveSlotButtons(n.world, saveSlotSelectDialog, n.context, setSaveSlot)
	for _, button := range buttons {
		if !button.occupied {
			continue
		}

		t, ok := transform.FindChildWithComponent(button.entry, component.Text)
		if !ok {
			panic("button has no text")
		}
		text := component.Text.Get(t)
		text.Text += " (overwrite)"
		text.Color = assets.WarningColor
	}

	closeButton := archetype.NewButton(n.world, "X", math.Vec2{X: 550, Y: 10}, false, func(w donburi.World, e *donburi.Entry) {
		n.Hide()
	})
	transform.AppendChild(saveSlotSelectDialog, closeButton, false)

	n.saveSlotSelectDialog = saveSlotSelectDialog
}

func (n *newGameMenu) hideSaveSlotSelect() {
	if n.saveSlotSelectDialog == nil {
		return
	}

	component.Destroy(n.saveSlotSelectDialog)
	n.saveSlotSelectDialog = nil
}

func (n *newGameMenu) Show() {
	n.saveSlot = nil
	n.character = domain.PlayerCharacterKing
	n.showSaveSlotSelect()
}

func (n *newGameMenu) Hide() {
	n.hideCharacterSelect()
	n.hideSaveSlotSelect()
}

type loadGameMenu struct {
	world   donburi.World
	context Context

	dialog *donburi.Entry
}

func newLoadGameMenu(w donburi.World, context Context) *loadGameMenu {
	l := &loadGameMenu{
		world:   w,
		context: context,
	}

	return l
}

func (l *loadGameMenu) Show() {
	if l.dialog != nil {
		return
	}

	loadGameDialog := newMainMenuDialog(l.world)

	archetype.New(l.world).
		WithParent(loadGameDialog).
		WithPosition(math.Vec2{X: 40, Y: 50}).
		WithText(component.TextData{
			Text: "Load Game",
		})

	loadGame := func(slot int) {
		l.context.SceneSwitcher.SwitchToLoadedGame(slot)
	}

	buttons := newSaveSlotButtons(l.world, loadGameDialog, l.context, loadGame)
	for _, button := range buttons {
		if button.occupied {
			continue
		}

		t, ok := transform.FindChildWithComponent(button.entry, component.Text)
		if !ok {
			panic("button has no text")
		}
		text := component.Text.Get(t)
		text.Text += " (empty)"
		text.Color = assets.InactiveColor

		component.Button.Get(button.entry).OnClick = func(w donburi.World, e *donburi.Entry) {}
	}

	for i := range buttons {
		button := buttons[i]
		if !button.occupied {
			continue
		}

		replayButton := archetype.NewButton(l.world, "Replay", math.Vec2{X: 300, Y: 0}, false, func(w donburi.World, e *donburi.Entry) {
			l.context.SceneSwitcher.SwitchToReplayedGame(button.slot)
		})
		transform.AppendChild(button.entry, replayButton, false)

		deleteButton := archetype.NewButton(l.world, "Delete", math.Vec2{X: 435, Y: 0}, false, func(w donburi.World, e *donburi.Entry) {
			err := l.context.Storage.DeleteSlot(button.slot)
			if err != nil {
				// TODO better error alerting
				panic(err)
			}
			l.Hide()
			l.Show()
		})
		transform.AppendChild(button.entry, deleteButton, false)
	}

	closeButton := archetype.NewButton(l.world, "X", math.Vec2{X: 550, Y: 10}, false, func(w donburi.World, e *donburi.Entry) {
		l.Hide()
	})
	transform.AppendChild(loadGameDialog, closeButton, false)

	l.dialog = loadGameDialog
}

func (l *loadGameMenu) Hide() {
	if l.dialog == nil {
		return
	}

	component.Destroy(l.dialog)
	l.dialog = nil
}

type saveSlotButton struct {
	slot     int
	entry    *donburi.Entry
	occupied bool
}

func newSaveSlotButtons(w donburi.World, parent *donburi.Entry, context Context, onClick func(slot int)) []saveSlotButton {
	occupiedSlots, err := context.Storage.OccupiedSlots()
	if err != nil {
		// TODO: better error alerting
		panic(err)
	}

	occupiedSlotsMap := map[int]struct{}{}
	for _, slot := range occupiedSlots {
		occupiedSlotsMap[slot] = struct{}{}
	}

	var buttons []saveSlotButton
	for i := 0; i < saveSlots; i++ {
		slot := i

		occupied := false
		if _, ok := occupiedSlotsMap[slot]; ok {
			occupied = true
		}

		name := fmt.Sprintf("Slot %v", slot+1)

		button := archetype.NewButtonWithWidth(w, name, math.Vec2{X: 40, Y: float64(120 + 100*slot)}, false, func(w donburi.World, e *donburi.Entry) {
			onClick(slot)
		}, 250)

		transform.AppendChild(parent, button, false)

		buttons = append(buttons, saveSlotButton{
			slot:     slot,
			entry:    button,
			occupied: occupied,
		})
	}

	return buttons
}

func newMainMenuDialog(w donburi.World) *donburi.Entry {
	frame := archetype.NewFrameWithBorder(w, math.Vec2{X: 200, Y: 100}, 600, 400)
	frame.AddComponent(component.Dialog)
	component.Layer.Get(frame).Layer = component.SpriteUILayerMessage
	return frame
}

func newSidePanel(w donburi.World, screenWidth int, screenHeight int) *donburi.Entry {
	width := 300

	sidePanelImage := ebiten.NewImage(width, screenHeight)
	sidePanelImage.Fill(assets.PanelColor)

	sidePanel := archetype.New(w).
		WithPosition(math.Vec2{X: float64(screenWidth - width), Y: 0}).
		WithLayer(component.SpriteUILayerBackground).
		WithSprite(component.SpriteData{

			Image: sidePanelImage,
		}).
		With(component.Collider).
		With(component.UI).
		With(component.UIPanel).
		Entry()

	component.Collider.SetValue(sidePanel, component.ColliderData{
		Width:  float64(width),
		Height: float64(screenHeight),
		Layer:  component.CollisionLayerButtons,
	})

	return sidePanel
}
