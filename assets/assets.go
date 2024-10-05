package assets

import (
	"bytes"
	"embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"io/fs"
	"math"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio/mp3"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type TerrainSprites struct {
	Spring *ebiten.Image
	Summer *ebiten.Image
	Fall   *ebiten.Image
	Winter *ebiten.Image
}

type DepositSprites struct {
	Spring SpritesByAmount
	Summer SpritesByAmount
	Fall   SpritesByAmount
	Winter SpritesByAmount
}

type SpritesByAmount struct {
	Full  *ebiten.Image
	Half  *ebiten.Image
	Empty *ebiten.Image
}

var (
	//go:embed fonts/kenney-mini-square.ttf
	squareFontData []byte

	//go:embed *
	assetsFS embed.FS

	WorldGeneratorConfig         GeneratorConfig
	TutorialWorldGeneratorConfig GeneratorConfig

	TileFog *ebiten.Image

	TerrainPlains    TerrainSprites
	TerrainForest    *ebiten.Image
	TerrainDesert    *ebiten.Image
	TerrainShore     TerrainSprites
	TerrainWater     TerrainSprites
	TerrainMountains *ebiten.Image

	DepositTrees     DepositSprites
	DepositRocks     DepositSprites
	DepositIronOre   DepositSprites
	DepositBerries   DepositSprites
	DepositMushrooms DepositSprites
	DepositDeer      DepositSprites

	TileBorderN *ebiten.Image
	TileBorderE *ebiten.Image
	TileBorderS *ebiten.Image
	TileBorderW *ebiten.Image

	IconWood    *ebiten.Image
	IconStone   *ebiten.Image
	IconFood    *ebiten.Image
	IconIron    *ebiten.Image
	IconDiamond *ebiten.Image

	IconCurrency    *ebiten.Image
	IconDocks       *ebiten.Image
	IconFarming     *ebiten.Image
	IconGranary     *ebiten.Image
	IconHunting     *ebiten.Image
	IconIronWorking *ebiten.Image
	IconKnighthood  *ebiten.Image
	IconMining      *ebiten.Image
	IconSentry      *ebiten.Image
	IconStonecraft  *ebiten.Image
	IconSwords      *ebiten.Image
	IconTrade       *ebiten.Image

	IconFrame    *ebiten.Image
	IconUpgrade  *ebiten.Image
	IconResearch *ebiten.Image
	IconMenu     *ebiten.Image

	IconNext      *ebiten.Image
	IconResources *ebiten.Image

	SettlementLevel1 *ebiten.Image
	SettlementLevel2 *ebiten.Image
	SettlementLevel3 *ebiten.Image
	SettlementLevel4 *ebiten.Image
	SettlementLevel5 *ebiten.Image

	Settlers *ebiten.Image

	FlagWhite *ebiten.Image

	SeasonSpring *ebiten.Image
	SeasonSummer *ebiten.Image
	SeasonFall   *ebiten.Image
	SeasonWinter *ebiten.Image

	CharacterSettler *ebiten.Image
	CharacterAdvisor *ebiten.Image
	CharacterHunter  *ebiten.Image
	CharacterFrame   *ebiten.Image

	PlayerKing     *ebiten.Image
	PlayerQueen    *ebiten.Image
	PlayerPrincess *ebiten.Image
	PlayerMerchant *ebiten.Image

	Road  *ebiten.Image
	RoadN *ebiten.Image
	RoadE *ebiten.Image
	RoadS *ebiten.Image
	RoadW *ebiten.Image

	RoadNS *ebiten.Image
	RoadWE *ebiten.Image

	RoadNE *ebiten.Image
	RoadSE *ebiten.Image
	RoadSW *ebiten.Image
	RoadNW *ebiten.Image

	RoadNSE *ebiten.Image
	RoadSWE *ebiten.Image
	RoadNSW *ebiten.Image
	RoadNWE *ebiten.Image

	RoadNSWE *ebiten.Image

	Crosshair *ebiten.Image
	Selected  *ebiten.Image

	SeasonsWheel *ebiten.Image

	LargeSquareFont  font.Face
	MediumSquareFont font.Face
	SmallSquareFont  font.Face

	Click1Sound []byte
	Click2Sound []byte
	Click3Sound []byte

	WinSound  []byte
	LoseSound []byte

	Chime1Sound []byte
	Chime2Sound []byte
	Chime3Sound []byte
	Chime4Sound []byte

	Theme1 []byte

	ShaderDistortion *ebiten.Shader

	Version = "dev"
)

type GeneratorConfig struct {
	Seed int64 `json:"seed"`

	Width  int `json:"width"`
	Height int `json:"height"`

	Altitude    LayerConfig `json:"altitude"`
	Temperature LayerConfig `json:"temperature"`
	Humidity    LayerConfig `json:"humidity"`
}

type LayerConfig struct {
	Scale  float64       `json:"scale"`
	Layers map[int]Range `json:"layers"`
}

type Range struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func MustLoadAssets() {
	WorldGeneratorConfig = mustLoadWorldGeneratorConfig("map/world.yaml")
	TutorialWorldGeneratorConfig = mustLoadWorldGeneratorConfig("map/tutorial.yaml")

	TileFog = mustLoadImage("tiles/fog.png")

	TerrainPlains = mustLoadTerrainSprites("plains")
	TerrainForest = mustLoadImage("tiles/terrains/forest.png")
	TerrainDesert = mustLoadImage("tiles/terrains/desert.png")
	TerrainShore = mustLoadTerrainSprites("shore")
	TerrainWater = mustLoadTerrainSprites("water")
	TerrainMountains = mustLoadImage("tiles/terrains/mountains.png")

	DepositTrees = mustLoadDepositSprites("trees")
	DepositRocks = mustLoadDepositSprites("rocks")
	DepositIronOre = mustLoadDepositSprites("iron-ore")
	DepositBerries = mustLoadDepositSprites("berries")
	DepositMushrooms = mustLoadDepositSprites("mushrooms")
	DepositDeer = mustLoadDepositSprites("deer")

	TileBorderN = mustLoadImage("tiles/border.png")
	TileBorderE = rotatedImage(TileBorderN, math.Pi*0.5)
	TileBorderS = rotatedImage(TileBorderN, math.Pi)
	TileBorderW = rotatedImage(TileBorderN, math.Pi*1.5)

	IconWood = mustLoadImage("icons/wood.png")
	IconStone = mustLoadImage("icons/stone.png")
	IconFood = mustLoadImage("icons/food.png")
	IconIron = mustLoadImage("icons/iron.png")
	IconDiamond = mustLoadImage("icons/diamond.png")

	IconCurrency = mustLoadImage("icons/currency.png")
	IconDocks = mustLoadImage("icons/docks.png")
	IconFarming = mustLoadImage("icons/farming.png")
	IconGranary = mustLoadImage("icons/granary.png")
	IconHunting = mustLoadImage("icons/hunting.png")
	IconIronWorking = mustLoadImage("icons/iron-working.png")
	IconKnighthood = mustLoadImage("icons/knighthood.png")
	IconMining = mustLoadImage("icons/mining.png")
	IconSentry = mustLoadImage("icons/sentry.png")
	IconStonecraft = mustLoadImage("icons/stonecraft.png")
	IconSwords = mustLoadImage("icons/swords.png")
	IconTrade = mustLoadImage("icons/trade.png")

	IconFrame = mustLoadImage("ui/icon-frame.png")
	IconUpgrade = mustLoadImage("icons/upgrade.png")
	IconResearch = mustLoadImage("icons/research.png")
	IconMenu = mustLoadImage("icons/menu.png")

	IconNext = mustLoadImage("icons/next.png")
	IconResources = mustLoadImage("icons/resources.png")

	SettlementLevel1 = mustLoadImage("buildings/settlement-1.png")
	SettlementLevel2 = mustLoadImage("buildings/settlement-2.png")
	SettlementLevel3 = mustLoadImage("buildings/settlement-3.png")
	SettlementLevel4 = mustLoadImage("buildings/settlement-4.png")
	SettlementLevel5 = mustLoadImage("buildings/settlement-5.png")

	Settlers = mustLoadImage("buildings/settlers.png")

	FlagWhite = mustLoadImage("tiles/flag-white.png")

	SeasonSpring = mustLoadImage("seasons/spring.png")
	SeasonSummer = mustLoadImage("seasons/summer.png")
	SeasonFall = mustLoadImage("seasons/fall.png")
	SeasonWinter = mustLoadImage("seasons/winter.png")

	CharacterSettler = mustLoadImage("characters/settler.png")
	CharacterAdvisor = mustLoadImage("characters/advisor.png")
	CharacterHunter = mustLoadImage("characters/hunter.png")
	CharacterFrame = mustLoadImage("characters/frame.png")

	PlayerKing = mustLoadImage("characters/player/king.png")
	PlayerQueen = mustLoadImage("characters/player/queen.png")
	PlayerPrincess = mustLoadImage("characters/player/princess.png")
	PlayerMerchant = mustLoadImage("characters/player/merchant.png")

	Road = mustLoadImage("roads/road.png")
	RoadN = mustLoadImage("roads/road-n.png")
	RoadE = mustLoadImage("roads/road-e.png")
	RoadS = mustLoadImage("roads/road-s.png")
	RoadW = mustLoadImage("roads/road-w.png")

	RoadNS = mustLoadImage("roads/road-ns.png")
	RoadWE = mustLoadImage("roads/road-we.png")

	RoadNE = mustLoadImage("roads/road-ne.png")
	RoadSE = mustLoadImage("roads/road-se.png")
	RoadSW = mustLoadImage("roads/road-sw.png")
	RoadNW = mustLoadImage("roads/road-nw.png")

	RoadNSE = mustLoadImage("roads/road-nse.png")
	RoadSWE = mustLoadImage("roads/road-swe.png")
	RoadNSW = mustLoadImage("roads/road-nsw.png")
	RoadNWE = mustLoadImage("roads/road-nwe.png")

	RoadNSWE = mustLoadImage("roads/road-nswe.png")

	Crosshair = mustLoadImage("crosshair.png")
	Selected = mustLoadImage("selected.png")

	LargeSquareFont = mustLoadFont(squareFontData, 24)
	MediumSquareFont = mustLoadFont(squareFontData, 16)
	SmallSquareFont = mustLoadFont(squareFontData, 8)

	Click1Sound = mustLoadMP3Stream("sounds/click1.mp3")
	Click2Sound = mustLoadMP3Stream("sounds/click2.mp3")
	Click3Sound = mustLoadMP3Stream("sounds/click3.mp3")

	WinSound = mustLoadMP3Stream("sounds/win.mp3")
	LoseSound = mustLoadMP3Stream("sounds/lose.mp3")

	Chime1Sound = mustLoadMP3Stream("sounds/chime1.mp3")
	Chime2Sound = mustLoadMP3Stream("sounds/chime2.mp3")
	Chime3Sound = mustLoadMP3Stream("sounds/chime3.mp3")
	Chime4Sound = mustLoadMP3Stream("sounds/chime4.mp3")

	Theme1 = mustLoadMP3Stream("sounds/theme1.mp3")

	ShaderDistortion = mustLoadShader("shaders/distortion.kage")

	SeasonsWheel = newSeasonsWheel()

	version, err := fs.ReadFile(assetsFS, "version")
	if err == nil {
		Version = strings.TrimSpace(string(version))
	}
}

func mustLoadFont(data []byte, size int) font.Face {
	f, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}

	return face
}

func mustLoadImage(path string) *ebiten.Image {
	data, err := fs.ReadFile(assetsFS, path)
	if err != nil {
		panic(err)
	}

	return mustNewEbitenImage(data)
}

func mustNewEbitenImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadMP3Stream(path string) []byte {
	data, err := fs.ReadFile(assetsFS, path)
	if err != nil {
		panic(err)
	}

	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	streamBytes, err := io.ReadAll(stream)
	if err != nil {
		panic(err)
	}

	return streamBytes
}

func mustLoadShader(path string) *ebiten.Shader {
	data, err := fs.ReadFile(assetsFS, path)
	if err != nil {
		panic(err)
	}

	shader, err := ebiten.NewShader(data)
	if err != nil {
		panic(err)
	}

	return shader
}

func mustLoadTerrainSprites(terrain string) TerrainSprites {
	basePath := fmt.Sprintf("tiles/terrains/%s", terrain)
	return TerrainSprites{
		Spring: mustLoadImage(path.Join(basePath, "spring.png")),
		Summer: mustLoadImage(path.Join(basePath, "summer.png")),
		Fall:   mustLoadImage(path.Join(basePath, "fall.png")),
		Winter: mustLoadImage(path.Join(basePath, "winter.png")),
	}
}

func mustLoadDepositSprites(deposit string) DepositSprites {
	return DepositSprites{
		Spring: mustLoadDepositSpriteSeason(deposit, "spring"),
		Summer: mustLoadDepositSpriteSeason(deposit, "summer"),
		Fall:   mustLoadDepositSpriteSeason(deposit, "fall"),
		Winter: mustLoadDepositSpriteSeason(deposit, "winter"),
	}
}

func mustLoadDepositSpriteSeason(deposit, season string) SpritesByAmount {
	basePath := fmt.Sprintf("tiles/deposits/%s/%s", deposit, season)
	return SpritesByAmount{
		Full:  mustLoadImage(path.Join(basePath, "full.png")),
		Half:  mustLoadImage(path.Join(basePath, "half.png")),
		Empty: mustLoadImage(path.Join(basePath, "empty.png")),
	}
}

func mustLoadWorldGeneratorConfig(path string) GeneratorConfig {
	data, err := fs.ReadFile(assetsFS, path)
	if err != nil {
		panic(err)
	}

	var config GeneratorConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func newSeasonsWheel() *ebiten.Image {
	size := 256
	center := float64(size) / 2
	img := ebiten.NewImage(size, size)
	drawTriangleInRect(img, 0, 0, float64(size), 0, center, center, SpringColor)
	drawTriangleInRect(img, float64(size), 0, float64(size), float64(size), center, center, SummerColor)
	drawTriangleInRect(img, float64(size), float64(size), 0, float64(size), center, center, FallColor)
	drawTriangleInRect(img, 0, float64(size), 0, 0, center, center, WinterColor)

	halfIcon := 8.0
	scale := 4.0
	margin := float64(size) / 8

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(margin*3, 0)
	img.DrawImage(SeasonSpring, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfIcon, -halfIcon)
	op.GeoM.Rotate(math.Pi / 2)
	op.GeoM.Translate(halfIcon, halfIcon)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(margin*6, margin*3)
	img.DrawImage(SeasonSummer, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfIcon, -halfIcon)
	op.GeoM.Rotate(math.Pi)
	op.GeoM.Translate(halfIcon, halfIcon)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(margin*3, margin*6)
	img.DrawImage(SeasonFall, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfIcon, -halfIcon)
	op.GeoM.Rotate(math.Pi * 1.5)
	op.GeoM.Translate(halfIcon, halfIcon)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(0, margin*3)
	img.DrawImage(SeasonWinter, op)

	return img
}

func drawTriangleInRect(img *ebiten.Image, x1, y1, x2, y2, x3, y3 float64, col color.RGBA) {
	src := ebiten.NewImage(1, 1)
	src.Fill(color.White)

	r := float32(col.R) / 255.0
	g := float32(col.G) / 255.0
	b := float32(col.B) / 255.0
	a := float32(col.A) / 255.0

	vertices := []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x2), DstY: float32(y2), ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x3), DstY: float32(y3), ColorR: r, ColorG: g, ColorB: b, ColorA: a},
	}
	img.DrawTriangles(vertices, []uint16{0, 1, 2}, src, nil)
}

func rotatedImage(img *ebiten.Image, rad float64) *ebiten.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newImg := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(width)/2, -float64(height)/2)
	op.GeoM.Rotate(rad)
	op.GeoM.Translate(float64(width)/2, float64(height)/2)
	newImg.DrawImage(img, op)

	return newImg
}
