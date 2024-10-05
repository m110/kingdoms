package generator

import (
	stdmath "math"
	"math/rand"

	"github.com/ojrac/opensimplex-go"

	"github.com/m110/kingdoms/assets"
	"github.com/m110/kingdoms/domain"
)

// The lower the R, the more points
var DepositOnTerrainR = map[domain.Terrain]map[domain.Deposit]int{
	domain.TerrainForest: {
		domain.DepositTrees:     1,
		domain.DepositDeer:      4,
		domain.DepositBerries:   4,
		domain.DepositMushrooms: 2,
	},
	domain.TerrainDesert: {
		domain.DepositRocks: 3,
	},
	domain.TerrainPlains: {
		domain.DepositRocks:     2,
		domain.DepositTrees:     2,
		domain.DepositBerries:   2,
		domain.DepositMushrooms: 2,
		domain.DepositDeer:      2,
	},
}

var TerrainsMap = map[Altitude]map[Temperature]map[Humidity]domain.Terrain{
	AltitudeLow: {
		TemperatureCold: {
			HumidityDry:    domain.TerrainWater,
			HumidityNormal: domain.TerrainWater,
			HumidityWet:    domain.TerrainWater,
		},
		TemperatureMild: {
			HumidityDry:    domain.TerrainWater,
			HumidityNormal: domain.TerrainWater,
			HumidityWet:    domain.TerrainWater,
		},
		TemperatureHot: {
			HumidityDry:    domain.TerrainWater,
			HumidityNormal: domain.TerrainWater,
			HumidityWet:    domain.TerrainWater,
		},
	},
	AltitudeLowMedium: {
		TemperatureCold: {
			HumidityDry:    domain.TerrainShore,
			HumidityNormal: domain.TerrainShore,
			HumidityWet:    domain.TerrainShore,
		},
		TemperatureMild: {
			HumidityDry:    domain.TerrainShore,
			HumidityNormal: domain.TerrainShore,
			HumidityWet:    domain.TerrainShore,
		},
		TemperatureHot: {
			HumidityDry:    domain.TerrainShore,
			HumidityNormal: domain.TerrainShore,
			HumidityWet:    domain.TerrainShore,
		},
	},
	AltitudeHigh: {
		TemperatureCold: {
			HumidityDry:    domain.TerrainMountains,
			HumidityNormal: domain.TerrainMountains,
			HumidityWet:    domain.TerrainMountains,
		},
		TemperatureMild: {
			HumidityDry:    domain.TerrainMountains,
			HumidityNormal: domain.TerrainMountains,
			HumidityWet:    domain.TerrainMountains,
		},
		TemperatureHot: {
			HumidityDry:    domain.TerrainMountains,
			HumidityNormal: domain.TerrainMountains,
			HumidityWet:    domain.TerrainMountains,
		},
	},
	AltitudeMedium: {
		TemperatureCold: {
			HumidityDry:    domain.TerrainPlains,
			HumidityNormal: domain.TerrainPlains,
			HumidityWet:    domain.TerrainPlains,
		},
		TemperatureMild: {
			HumidityDry:    domain.TerrainPlains,
			HumidityNormal: domain.TerrainPlains,
			HumidityWet:    domain.TerrainForest,
		},
		TemperatureHot: {
			HumidityDry:    domain.TerrainDesert,
			HumidityNormal: domain.TerrainPlains,
			HumidityWet:    domain.TerrainPlains,
		},
	},
}

type LayerID = int

type Altitude = LayerID

type Layer struct {
	LayerID LayerID
	Value   float64
}

const (
	AltitudeLow Altitude = iota
	AltitudeLowMedium
	AltitudeMedium
	AltitudeHigh
)

type Temperature = LayerID

const (
	TemperatureCold Temperature = iota
	TemperatureMild
	TemperatureHot
)

type Humidity = LayerID

const (
	HumidityDry Humidity = iota
	HumidityNormal
	HumidityWet
)

type World struct {
	Altitude    [][]Layer
	Temperature [][]Layer
	Humidity    [][]Layer

	Terrains [][]domain.Terrain
	Deposits [][]*domain.Deposit
}

func ConfigFromAssets(cfg assets.GeneratorConfig) WorldGeneratorConfig {
	return WorldGeneratorConfig{
		Seed: cfg.Seed,

		Width:  cfg.Width,
		Height: cfg.Height,

		Altitude:    layerConfigFromAssets(cfg.Altitude),
		Temperature: layerConfigFromAssets(cfg.Temperature),
		Humidity:    layerConfigFromAssets(cfg.Humidity),
	}
}

func layerConfigFromAssets(cfg assets.LayerConfig) LayerConfig {
	layer := make(map[LayerID]Range, len(cfg.Layers))
	for layerID, r := range cfg.Layers {
		layer[layerID] = Range{
			Min: r.Min,
			Max: r.Max,
		}
	}
	return LayerConfig{
		Scale:  cfg.Scale,
		Layers: layer,
	}
}

type WorldGeneratorConfig struct {
	Seed int64

	Width  int
	Height int

	Altitude    LayerConfig
	Temperature LayerConfig
	Humidity    LayerConfig
}

type LayerConfig struct {
	Scale  float64
	Layers map[LayerID]Range
}

type Range struct {
	Min float64
	Max float64
}

type WorldGenerator struct {
	config WorldGeneratorConfig
}

func NewWorldGenerator(config WorldGeneratorConfig) *WorldGenerator {
	if config.Seed == 0 {
		config.Seed = rand.Int63n(stdmath.MaxInt64)
	}
	return &WorldGenerator{
		config: config,
	}
}

func (w *WorldGenerator) Generate() World {
	gw := World{
		Altitude:    w.generateLayers(w.config.Seed*19+19, w.config.Altitude.Scale, w.config.Altitude.Layers),
		Temperature: w.generateLayers(w.config.Seed*191+191, w.config.Temperature.Scale, w.config.Temperature.Layers),
		Humidity:    w.generateLayers(w.config.Seed*1919+1919, w.config.Humidity.Scale, w.config.Humidity.Layers),
	}

	gw.Terrains = w.generateTerrains(gw)
	gw.Deposits = w.generateDeposits(gw)

	return gw
}

func (w *WorldGenerator) generateLayers(
	seed int64,
	scale float64,
	ranges map[int]Range,
) [][]Layer {
	data := make([][]Layer, w.config.Width)

	simplex := opensimplex.NewNormalized(seed)

	for x := 0; x < w.config.Width; x++ {
		data[x] = make([]Layer, w.config.Height)

		for y := 0; y < w.config.Height; y++ {
			normalizedX := float64(x)/float64(w.config.Width) - 0.5
			normalizedY := float64(y)/float64(w.config.Height) - 0.5

			noise := simplex.Eval2(normalizedX*scale, normalizedY*scale)

			for layerID, r := range ranges {
				if noise >= r.Min && noise < r.Max {
					data[x][y] = Layer{
						LayerID: layerID,
						Value:   noise,
					}
					break
				}
			}
		}
	}

	return data
}

func (w *WorldGenerator) blueNoise(seed int64, r int) [][]bool {
	result := make([][]bool, w.config.Width)
	blueNoise := make([][]float64, w.config.Width)

	for x := 0; x < w.config.Width; x++ {
		result[x] = make([]bool, w.config.Height)
		blueNoise[x] = make([]float64, w.config.Height)
	}

	simplex := opensimplex.NewNormalized(seed)

	for x := 0; x < w.config.Width; x++ {
		for y := 0; y < w.config.Height; y++ {
			normalizedX := float64(x)/float64(w.config.Width) - 0.5
			normalizedY := float64(y)/float64(w.config.Height) - 0.5

			blueNoise[x][y] = simplex.Eval2(normalizedX*50, normalizedY*50)
		}
	}

	for xc := 0; xc < w.config.Width; xc++ {
		for yc := 0; yc < w.config.Height; yc++ {
			max := 0.0

			for dx := -r; dx <= r; dx++ {
				for dy := -r; dy <= r; dy++ {
					xn := dx + xc
					yn := dy + yc

					if 0 <= xn && xn < w.config.Width && 0 <= yn && yn < w.config.Height {
						e := blueNoise[xn][yn]
						if e > max {
							max = e
						}
					}
				}
			}

			if blueNoise[xc][yc] == max {
				result[xc][yc] = true
			}
		}
	}

	return result
}

func (w *WorldGenerator) generateTerrains(world World) [][]domain.Terrain {
	data := make([][]domain.Terrain, w.config.Width)

	for x := 0; x < w.config.Width; x++ {
		data[x] = make([]domain.Terrain, w.config.Height)

		for y := 0; y < w.config.Height; y++ {
			altitude := world.Altitude[x][y].LayerID
			temperature := world.Temperature[x][y].LayerID
			humidity := world.Humidity[x][y].LayerID

			data[x][y] = TerrainsMap[altitude][temperature][humidity]
		}
	}

	return data
}

func (w *WorldGenerator) generateDeposits(world World) [][]*domain.Deposit {
	data := make([][]*domain.Deposit, w.config.Width)
	for x := 0; x < w.config.Width; x++ {
		data[x] = make([]*domain.Deposit, w.config.Height)
	}

	var factor int64 = 33

	for deposit := range domain.Deposits {
		deposit := deposit

		noiseOnTerrain := map[domain.Terrain][][]bool{}

		for terrain := range domain.Terrains {
			r, ok := DepositOnTerrainR[terrain][deposit]
			if !ok {
				continue
			}

			noise := w.blueNoise(w.config.Seed*factor+factor, r)
			noiseOnTerrain[terrain] = noise
		}

		for x := 0; x < w.config.Width; x++ {
			for y := 0; y < w.config.Height; y++ {
				if data[x][y] != nil {
					continue
				}

				terrain := world.Terrains[x][y]
				// TODO Could be defined in the domain
				if terrain == domain.TerrainWater || terrain == domain.TerrainShore || terrain == domain.TerrainMountains {
					continue
				}

				noise, ok := noiseOnTerrain[terrain]
				if !ok {
					continue
				}

				if noise[x][y] {
					data[x][y] = &deposit
				}
			}
		}

		factor += 33
	}

	return data
}
