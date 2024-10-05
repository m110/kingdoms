package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/scene"
	"github.com/m110/kingdoms/storage"
	"image/color"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	sceneContext scene.Context
	scene        Scene

	nextSceneFunc func() Scene
	loadingScreen *ebiten.Image

	storage *storage.Storage

	musicOn     bool
	themePlayer *audio.Player
}

type Config struct {
	StartScene string

	ScreenWidth  int
	ScreenHeight int

	SampleRate int
}

func NewGame(config Config) *Game {
	assets.MustLoadAssets()
	domain.LoadDomain()

	audioCtx := audio.NewContext(config.SampleRate)
	themePlayer := audioCtx.NewPlayerFromBytes(assets.Theme1)

	s, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	loadingScreen := ebiten.NewImage(config.ScreenWidth, config.ScreenHeight)
	loadingScreen.Fill(color.Black)
	text.Draw(loadingScreen, "Loading...", assets.LargeSquareFont, config.ScreenWidth/2-100, config.ScreenHeight/2, color.White)

	g := &Game{
		themePlayer:   themePlayer,
		musicOn:       false,
		loadingScreen: loadingScreen,
	}

	g.sceneContext = scene.Context{
		Version:       assets.Version,
		Storage:       s,
		SceneSwitcher: g,
		SetMusic:      g.setMusic,
		ScreenWidth:   config.ScreenWidth,
		ScreenHeight:  config.ScreenHeight,
	}

	switch config.StartScene {
	case "":
		g.SwitchToMainMenu()
	default:
		panic("unknown start scene")
	}

	return g
}

func (g *Game) setMusic(b bool) {
	g.musicOn = b
	if b {
		g.themePlayer.SetVolume(0.5)
		g.themePlayer.Play()
	} else {
		g.themePlayer.Pause()
	}
}

func (g *Game) switchScene(newSceneFunc func() Scene) {
	g.scene = nil
	g.nextSceneFunc = newSceneFunc
}

func (g *Game) SwitchToMainMenu() {
	g.switchScene(func() Scene {
		return scene.NewMainMenu(g.sceneContext)
	})
}

func (g *Game) SwitchToNewGame(saveSlot int, character domain.PlayerCharacterType) {
	g.switchScene(func() Scene {
		return scene.NewBoard(g.sceneContext, &saveSlot, character, false)
	})
}

func (g *Game) SwitchToLoadedGame(saveSlot int) {
	g.switchScene(func() Scene {
		game, err := g.sceneContext.Storage.LoadSlot(saveSlot)
		if err != nil {
			// TODO better error handling
			panic(err)
		}
		replay, err := g.sceneContext.Storage.LoadReplay(saveSlot)
		if err != nil {
			// TODO better error handling
			panic(err)
		}

		return scene.NewBoardFromSave(g.sceneContext, saveSlot, game, replay)
	})
}

func (g *Game) SwitchToReplayedGame(saveSlot int) {
	g.switchScene(func() Scene {
		replay, err := g.sceneContext.Storage.LoadReplay(saveSlot)
		if err != nil {
			// TODO better error handling
			panic(err)
		}

		return scene.NewBoardFromReplay(g.sceneContext, saveSlot, replay)
	})
}

func (g *Game) SwitchToTutorial() {
	g.switchScene(func() Scene {
		return scene.NewBoard(g.sceneContext, nil, domain.PlayerCharacterKing, true)
	})
}

func (g *Game) Update() error {
	if g.nextSceneFunc != nil {
		g.scene = g.nextSceneFunc()
		g.nextSceneFunc = nil
	}

	g.scene.Update()

	if g.musicOn && !g.themePlayer.IsPlaying() {
		_ = g.themePlayer.Rewind()
		g.themePlayer.Play()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.scene == nil {
		screen.DrawImage(g.loadingScreen, nil)
		return
	}

	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	if g.sceneContext.ScreenWidth == 0 || g.sceneContext.ScreenHeight == 0 {
		return width, height
	}

	return g.sceneContext.ScreenWidth, g.sceneContext.ScreenHeight
}
