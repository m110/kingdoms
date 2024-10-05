package system

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"

	"github.com/m110/kingdoms/domain"

	"github.com/m110/kingdoms/archetype"
	"github.com/m110/kingdoms/component"
	"github.com/m110/kingdoms/engine"
	"github.com/m110/kingdoms/events"
	events2 "github.com/m110/kingdoms/events"
)

type Research struct{}

func NewResearch() *Research {
	return &Research{}
}

func (r *Research) Init(w donburi.World) {
	events.ResearchRequestedEvent.Subscribe(w, r.OnResearchRequested)
}

func (r *Research) Update(w donburi.World) {}

func (r *Research) OnResearchRequested(w donburi.World, event events.ResearchRequested) {
	tech := domain.Technologies[event.Technology]
	research := component.Research.Get(engine.MustFindWithComponent(w, component.Research))

	if research.KnownTechnologies[event.Technology] {
		return
	}

	for _, t := range domain.Technologies {
		for _, e := range t.Enables {
			if e == event.Technology && !research.KnownTechnologies[t.ID] {
				return
			}
		}
	}

	resources := component.Resources.Get(engine.MustFindWithComponent(w, component.Resources))

	delta := events2.ResourcesDelta{
		Food:  -tech.Cost.Food,
		Stone: -tech.Cost.Stone,
		Wood:  -tech.Cost.Wood,
		Iron:  -tech.Cost.Iron,
		Gold:  -tech.Cost.Gold,
	}

	diff := events2.ResourcesDelta{
		Food:  resources.Food + delta.Food,
		Stone: resources.Stone + delta.Stone,
		Wood:  resources.Wood + delta.Wood,
		Iron:  resources.Iron + delta.Iron,
		Gold:  resources.Gold + delta.Gold,
	}

	if diff.Food < 0 || diff.Stone < 0 || diff.Wood < 0 || diff.Iron < 0 || diff.Gold < 0 {
		return
	}

	UpdateResources(w, delta)
	research.KnownTechnologies[event.Technology] = true

	events.TechnologyResearchedEvent.Publish(w, events.TechnologyResearched{
		Technology: event.Technology,
	})

	// TODO not sure if should be here?
	switch event.Technology {
	case domain.TechnologyHunting:
		archetype.ShowMessage(w, archetype.CharacterHunter, "{player}, we learned how to track deer.\n\nWe can now hunt deer for food.")
	}

	// TODO Probably a hack - close and open the panel instead of updating
	panel, ok := donburi.NewQuery(filter.Contains(component.ResearchPanel)).First(w)
	if ok {
		component.Destroy(panel)
		archetype.NewResearchPanel(w, &event.Technology)
	}
}
