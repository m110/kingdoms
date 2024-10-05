package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/domain"
	"github.com/yohamta/donburi"

	"github.com/m110/kingdoms/save"
)

type Context struct {
	Version string

	ScreenWidth  int
	ScreenHeight int

	Storage       Storage
	SceneSwitcher Switcher

	SetMusic func(bool)
}

type Switcher interface {
	SwitchToMainMenu()
	SwitchToNewGame(saveSlot int, character domain.PlayerCharacterType)
	SwitchToLoadedGame(saveSlot int)
	SwitchToReplayedGame(saveSlot int)
	SwitchToTutorial()
}

type Storage interface {
	LoadGlobalData() (*save.Global, error)
	OccupiedSlots() ([]int, error)
	LoadSlot(slot int) (*save.Game, error)
	LoadReplay(slot int) (*save.GameReplay, error)
	DeleteSlot(slot int) error

	SaveGlobalData(*save.Global) error
	SaveSlot(int, *save.Global, *save.Game) error
	SaveReplay(int, *save.GameReplay) error
}

type Initializable interface {
	// Init is called once the world is set up
	Init(w donburi.World)
}

type System interface {
	Initializable
	Update(w donburi.World)
}

type Drawable interface {
	Initializable
	Draw(w donburi.World, screen *ebiten.Image)
}
