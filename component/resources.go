package component

import "github.com/yohamta/donburi"

type ResourcesData struct {
	Food     int
	Stone    int
	Wood     int
	Iron     int
	Gold     int
	Diamonds int
}

var Resources = donburi.NewComponentType[ResourcesData]()
