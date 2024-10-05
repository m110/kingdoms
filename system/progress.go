package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"

	"github.com/m110/kingdoms/domain"

	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/events"
)

const targetMaxLevelSettlements = 5

type Progress struct {
	settlementsQuery *donburi.Query
	gameWon          bool
}

func NewProgress() *Progress {
	return &Progress{
		settlementsQuery: query.NewQuery(filter.Contains(component.Settlement)),
	}
}

func (p *Progress) Init(w donburi.World) {
	events.ResourcesUpdatedEvent.Subscribe(w, p.onResourcesUpdated)
	events.SettlementSpawnedEvent.Subscribe(w, p.onSettlementSpawned)
	events.SettlementUpgradedEvent.Subscribe(w, p.onSettlementUpgraded)
}

func (p *Progress) Update(w donburi.World) {}

func (p *Progress) onResourcesUpdated(w donburi.World, event events.ResourcesUpdated) {
	if p.gameWon {
		return
	}

	if event.Resources.Food < domain.SettlementUpgradeCost && event.Resources.Stone <= 0 && event.Resources.Wood < domain.SettlementCost {
		events.GameOverEvent.Publish(w, events.GameOver{})
	}
}

func (p *Progress) onSettlementSpawned(w donburi.World, event events.SettlementSpawned) {
	p.checkGameWon(w)
}

func (p *Progress) onSettlementUpgraded(w donburi.World, event events.SettlementUpgraded) {
	p.checkGameWon(w)
}

func (p *Progress) checkGameWon(w donburi.World) {
	if p.gameWon {
		return
	}

	progress := p.calculateProgress(w)

	if progress >= targetMaxLevelSettlements {
		events.GameWonEvent.Publish(w, events.GameWon{})
		p.gameWon = true
	}
}

func (p *Progress) calculateProgress(w donburi.World) int {
	maxLevelSettlements := 0

	p.settlementsQuery.Each(w, func(entry *donburi.Entry) {
		level := component.Settlement.Get(entry).Level
		if level >= domain.MaxSettlementLevelWithCurrency {
			maxLevelSettlements++
		}
	})

	return maxLevelSettlements
}
