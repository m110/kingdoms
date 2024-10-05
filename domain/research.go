package domain

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/m110/kingdoms/assets"
)

type Technology int

type TechnologyData struct {
	ID          Technology
	Name        string
	Icon        *ebiten.Image
	Description string
	Cost        ResourceCost
	Enables     []Technology
}

var Technologies map[Technology]TechnologyData

func loadTechnologies() {
	t := map[Technology]TechnologyData{
		TechnologyMining: {
			Icon:        assets.IconMining,
			Name:        "Mining",
			Description: "Collect 1 extra stone from mountains.",
			Cost: ResourceCost{
				Wood: 5,
			},
			Enables: []Technology{
				TechnologyStonecraft,
				TechnologyIronworking,
			},
		},
		TechnologyStonecraft: {
			Icon:        assets.IconStonecraft,
			Name:        "Stonecraft",
			Description: "Building roads on your territory costs no stone.",
			Cost: ResourceCost{
				Stone: 10,
			},
			Enables: []Technology{
				TechnologySentryTower,
			},
		},
		TechnologyIronworking: {
			Icon:        assets.IconIronWorking,
			Name:        "Ironworking",
			Description: "You're able to find and collect iron.",
			Cost: ResourceCost{
				Stone: 5,
				Wood:  5,
			},
			Enables: []Technology{
				TechnologySwords,
			},
		},
		TechnologySentryTower: {
			Icon:        assets.IconSentry,
			Name:        "Sentry Tower",
			Description: "Increases the sight range of all settlements by 3.",
			Cost: ResourceCost{
				Stone: 5,
				Wood:  5,
			},
		},
		TechnologySwords: {
			Icon:        assets.IconSwords,
			Name:        "Swords",
			Description: "Allows invading enemy territory.",
			Cost: ResourceCost{
				Iron: 5,
			},
			Enables: []Technology{
				TechnologyKnighthood,
			},
		},
		TechnologyKnighthood: {
			Icon:        assets.IconKnighthood,
			Name:        "Knighthood",
			Description: "You're able to find and enter ancient ruins.",
			Cost: ResourceCost{
				Food: 5,
				Iron: 5,
			},
		},
		TechnologyGranary: {
			Icon:        assets.IconGranary,
			Name:        "Granary",
			Description: "Increases all collected food by 1.",
			Cost: ResourceCost{
				Stone: 5,
				Wood:  5,
			},
			Enables: []Technology{
				TechnologyHunting,
				TechnologyFarming,
			},
		},
		TechnologyHunting: {
			Icon:        assets.IconHunting,
			Name:        "Hunting",
			Description: "You're able to find and hunt animals.",
			Cost: ResourceCost{
				Wood: 5,
			},
		},
		TechnologyFarming: {
			Icon:        assets.IconFarming,
			Name:        "Farming",
			Description: "Allows building farms",
			Cost: ResourceCost{
				Wood: 10,
			},
		},
		TechnologyCurrency: {
			Icon:        assets.IconCurrency,
			Name:        "Currency",
			Description: "Allows building castles (settlements level 5).\nAllows collecting gold.",
			Cost: ResourceCost{
				Food:  5,
				Stone: 5,
				Wood:  5,
			},
			Enables: []Technology{
				TechnologyTrade,
				TechnologyDocks,
			},
		},
		TechnologyDocks: {
			Icon:        assets.IconDocks,
			Name:        "Docks",
			Description: "Allows building docks and building settlements overseas.",
			Cost: ResourceCost{
				Wood: 10,
			},
		},
		TechnologyTrade: {
			Icon:        assets.IconTrade,
			Name:        "Trade",
			Description: "Allows trading with free cities.",
			Cost: ResourceCost{
				Food: 10,
			},
		},
	}

	for k, tech := range t {
		tech.ID = k
		t[k] = tech
	}

	Technologies = t
}

var RootTechnologies = []Technology{
	TechnologyMining,
	TechnologyGranary,
	TechnologyCurrency,
}

const (
	TechnologyMining Technology = iota
	TechnologyStonecraft
	TechnologyIronworking
	TechnologySentryTower
	TechnologySwords
	TechnologyKnighthood
	TechnologyGranary
	TechnologyHunting
	TechnologyFarming
	TechnologyCurrency
	TechnologyDocks
	TechnologyTrade
)
