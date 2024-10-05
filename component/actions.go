package component

import (
	"github.com/m110/kingdoms/save"
	"github.com/yohamta/donburi"
	"sort"
)

type Action struct {
	Tick       int
	ActionType save.ActionType
	Payload    *save.ActionPayload
}

type ActionsRecorderData struct {
	CurrentTick int
	Actions     []Action
}

func (d *ActionsRecorderData) Load(lastTick int, actions []Action) {
	d.CurrentTick = lastTick
	d.Actions = actions
}

func (d *ActionsRecorderData) Record(actionType save.ActionType, payload *save.ActionPayload) {
	d.Actions = append(d.Actions, Action{
		Tick:       d.CurrentTick,
		ActionType: actionType,
		Payload:    payload,
	})
}

var ActionsRecorder = donburi.NewComponentType[ActionsRecorderData]()

type ActionsPlayerData struct {
	CurrentTick float64
	Actions     []Action

	Finished bool
	Paused   bool
	Speed    float64
}

func (a *ActionsPlayerData) Load(actions []Action) {
	// Make sure actions are sorted by tick
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Tick < actions[j].Tick
	})

	a.CurrentTick = 0
	a.Actions = actions
	a.Speed = 1.0
}

func (a *ActionsPlayerData) Play() {
	if !a.Finished {
		a.Paused = false
	}
}

func (a *ActionsPlayerData) Pause() {
	a.Paused = true
}

func (a *ActionsPlayerData) Finish() {
	a.Finished = true
	a.Pause()
}

func (a *ActionsPlayerData) SetSpeed(speed float64) {
	a.Speed = speed
}

func (a *ActionsPlayerData) GetActionsForCurrentTick() []Action {
	var actions []Action
	lastIndex := 0

	for i, action := range a.Actions {
		if action.Tick > int(a.CurrentTick) {
			break
		}

		lastIndex = i
		actions = append(actions, action)
	}

	if len(actions) > 0 {
		a.Actions = a.Actions[lastIndex+1:]
	}

	if len(a.Actions) == 0 {
		a.Finish()
	}

	return actions
}

var ActionsPlayer = donburi.NewComponentType[ActionsPlayerData]()
