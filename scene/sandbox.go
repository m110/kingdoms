package scene

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	donburievents "github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
	"golang.org/x/image/colornames"
	"gopkg.in/yaml.v3"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/domain"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/scene/generator"
	"github.com/m110/kingdoms/system"
)

var (
	TagAltitude    = donburi.NewTag()
	TagTemperature = donburi.NewTag()
	TagHumidity    = donburi.NewTag()
)

type Sandbox struct {
	context   Context
	world     donburi.World
	systems   []System
	drawables []Drawable

	showTemperature bool
	showHumidity    bool
	showAltitude    bool

	debugAltitudeImages    map[generator.Altitude]*ebiten.Image
	debugTemperatureImages map[generator.Temperature]*ebiten.Image
	debugHumidityImages    map[generator.Humidity]*ebiten.Image
}

func NewSandbox(context Context) *Sandbox {
	controls := system.NewControls(true)

	systems := []System{
		system.NewDebug(),
		controls,
		system.NewAnimation(),
		system.NewText(),
		system.NewTimeToLive(),
		system.NewDestroy(),
	}

	drawables := []Drawable{
		system.NewRenderer(),
		controls,
	}

	s := &Sandbox{
		context:   context,
		systems:   systems,
		drawables: drawables,

		debugAltitudeImages: map[generator.Altitude]*ebiten.Image{
			generator.AltitudeLow:       altitudeDebugImage(generator.AltitudeLow),
			generator.AltitudeLowMedium: altitudeDebugImage(generator.AltitudeLowMedium),
			generator.AltitudeMedium:    altitudeDebugImage(generator.AltitudeMedium),
			generator.AltitudeHigh:      altitudeDebugImage(generator.AltitudeHigh),
		},
		debugTemperatureImages: map[generator.Temperature]*ebiten.Image{
			generator.TemperatureCold: temperatureDebugImage(generator.TemperatureCold),
			generator.TemperatureMild: temperatureDebugImage(generator.TemperatureMild),
			generator.TemperatureHot:  temperatureDebugImage(generator.TemperatureHot),
		},
		debugHumidityImages: map[generator.Humidity]*ebiten.Image{
			generator.HumidityDry:    humidityDebugImage(generator.HumidityDry),
			generator.HumidityNormal: humidityDebugImage(generator.HumidityNormal),
			generator.HumidityWet:    humidityDebugImage(generator.HumidityWet),
		},
	}

	s.world = s.createWorld()
	s.init()

	return s
}

func (s *Sandbox) createUI(w donburi.World) {
	panel := archetype.NewFrameWithBorder(
		w,
		math.Vec2{X: float64(s.context.ScreenWidth - 250), Y: 100},
		250,
		500,
	)
	panel.AddComponent(component.UIPanel)

	reloadButton := archetype.NewButton(w, "Reload", math.Vec2{X: 20, Y: 50}, false, func(w donburi.World, e *donburi.Entry) {
		s.world = s.createWorld()
		s.init()
	})
	transform.AppendChild(panel, reloadButton, false)

	toggleBoardButton := archetype.NewButton(w, "Toggle Board", math.Vec2{X: 20, Y: 150}, false, func(w donburi.World, e *donburi.Entry) {
		query.NewQuery(filter.Contains(component.Tile)).Each(w, func(e *donburi.Entry) {
			if !e.HasComponent(component.Active) {
				e.AddComponent(component.Active)
				component.Active.SetValue(e, component.ActiveData{
					Active: false,
				})
				return
			}

			active := component.Active.Get(e)
			active.Active = !active.Active
		})
	})
	transform.AppendChild(panel, toggleBoardButton, false)

	toggleAltitudeButton := archetype.NewButton(w, "Toggle Altitude", math.Vec2{X: 20, Y: 250}, false, func(w donburi.World, e *donburi.Entry) {
		s.showAltitude = !s.showAltitude
		toggleDebugTiles(w, TagAltitude, s.showAltitude)
	})
	transform.AppendChild(panel, toggleAltitudeButton, false)

	toggleTemperatureButton := archetype.NewButton(w, "Toggle Temp", math.Vec2{X: 20, Y: 310}, false, func(w donburi.World, e *donburi.Entry) {
		s.showTemperature = !s.showTemperature
		toggleDebugTiles(w, TagTemperature, s.showTemperature)
	})
	transform.AppendChild(panel, toggleTemperatureButton, false)

	toggleHumidityButton := archetype.NewButton(w, "Toggle Humidity", math.Vec2{X: 20, Y: 370}, false, func(w donburi.World, e *donburi.Entry) {
		s.showHumidity = !s.showHumidity
		toggleDebugTiles(w, TagHumidity, s.showHumidity)
	})
	transform.AppendChild(panel, toggleHumidityButton, false)

	toggleDebugTiles(w, TagAltitude, s.showAltitude)
	toggleDebugTiles(w, TagTemperature, s.showTemperature)
	toggleDebugTiles(w, TagHumidity, s.showHumidity)
}

func toggleDebugTiles(w donburi.World, tag donburi.IComponentType, isActive bool) {
	query.NewQuery(filter.Contains(tag)).Each(w, func(e *donburi.Entry) {
		active := component.Active.Get(e)
		active.Active = isActive
	})
}

func (s *Sandbox) createWorld() donburi.World {
	w := donburi.NewWorld()
	archetype.NewCamera(w, math.Vec2{}, engine.FloatRange{Min: 0.2, Max: 1.0})

	game := w.Entry(w.Create(component.Game))
	donburi.SetValue(game, component.Game, component.GameData{
		Settings: component.Settings{
			ScreenWidth:  s.context.ScreenWidth,
			ScreenHeight: s.context.ScreenHeight,
		},
	})

	progress := w.Entry(w.Create(component.Progress))
	component.Progress.SetValue(progress, component.ProgressData{
		Season: domain.SeasonSummer,
		Day:    1,
		Year:   1,
	})

	mapSize := 200

	w.Create(component.Debug)

	board := newBoardEntry(w, mapSize, mapSize)
	chunks := system.GenerateChunks(w, board, mapSize, mapSize)

	var config generator.WorldGeneratorConfig
	configContent, err := os.ReadFile("sandbox.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		panic(err)
	}

	gen := generator.NewWorldGenerator(config)
	gw := gen.Generate()

	s.createDebugEntries(w, chunks, gw)

	tiles := system.CreateTiles(w, chunks, gw.Terrains, gw.Deposits, mapSize, mapSize, false, true)
	component.Board.Get(board).Tiles = tiles

	s.createUI(w)

	return w
}

func (s *Sandbox) init() {
	for _, ss := range s.systems {
		ss.Init(s.world)
	}
	for _, d := range s.drawables {
		d.Init(s.world)
	}

	camera := archetype.MustFindCamera(s.world)
	cam := component.Camera.Get(camera)
	transform.GetTransform(camera).LocalScale = math.Vec2{
		X: cam.Zoom.Min,
		Y: cam.Zoom.Min,
	}
}

func (s *Sandbox) Update() {
	for _, d := range s.systems {
		d.Update(s.world)
	}

	donburievents.ProcessAllEvents(s.world)
}

func (s *Sandbox) Draw(screen *ebiten.Image) {
	for _, d := range s.drawables {
		d.Draw(s.world, screen)
	}
}

func (s *Sandbox) createDebugEntries(world donburi.World, chunks [][]*donburi.Entry, gw generator.World) {
	for x := 0; x < len(gw.Altitude); x++ {
		for y := 0; y < len(gw.Altitude[0]); y++ {
			altitude := gw.Altitude[x][y]
			temperature := gw.Temperature[x][y]
			humidity := gw.Humidity[x][y]

			chunkX := x / system.ChunkSize
			chunkY := y / system.ChunkSize
			chunk := chunks[chunkX][chunkY]

			newDebugTile(world, chunk, x, y, component.SpriteLayerFog, s.debugAltitudeImages[altitude.LayerID], TagAltitude)
			newDebugTile(world, chunk, x, y, component.SpriteLayerFog+1, s.debugTemperatureImages[temperature.LayerID], TagTemperature)
			newDebugTile(world, chunk, x, y, component.SpriteLayerFog+2, s.debugHumidityImages[humidity.LayerID], TagHumidity)
		}
	}
}

func newDebugTile(
	w donburi.World,
	parent *donburi.Entry,
	x int,
	y int,
	layer component.LayerID,
	image *ebiten.Image,
	tag donburi.IComponentType,
) {
	archetype.New(w).
		WithParent(parent).
		WithPosition(math.Vec2{
			X: float64(x % system.ChunkSize * system.TileSize),
			Y: float64(y % system.ChunkSize * system.TileSize),
		}).
		WithScale(math.Vec2{X: 2, Y: 2}).
		WithLayer(layer).
		With(component.Active).
		WithSprite(component.SpriteData{
			Image: image,
		}).
		With(tag)
}

func altitudeDebugImage(id generator.Altitude) *ebiten.Image {
	img := ebiten.NewImage(16, 16)

	switch id {
	case generator.AltitudeLow:
		img.Fill(semiTransparent(colornames.White))
	case generator.AltitudeLowMedium:
		img.Fill(semiTransparent(colornames.Lightslategray))
	case generator.AltitudeMedium:
		img.Fill(semiTransparent(colornames.Lightgray))
	case generator.AltitudeHigh:
		img.Fill(semiTransparent(colornames.Black))
	}

	return img
}

func temperatureDebugImage(temperature generator.LayerID) *ebiten.Image {
	img := ebiten.NewImage(16, 16)

	switch temperature {
	case generator.TemperatureCold:
		img.Fill(semiTransparent(colornames.Lightblue))
	case generator.TemperatureMild:
		img.Fill(semiTransparent(colornames.Lightyellow))
	case generator.TemperatureHot:
		img.Fill(semiTransparent(colornames.Darkorange))
	}

	return img
}

func humidityDebugImage(humidity generator.LayerID) *ebiten.Image {
	img := ebiten.NewImage(16, 16)

	switch humidity {
	case generator.HumidityDry:
		img.Fill(semiTransparent(colornames.Lightgoldenrodyellow))
	case generator.HumidityNormal:
		img.Fill(semiTransparent(colornames.Darkorange))
	case generator.HumidityWet:
		img.Fill(semiTransparent(colornames.Indianred))
	}

	return img
}

func semiTransparent(c color.RGBA) color.Color {
	rgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	rgba.A = 128

	return rgba
}
