package engine

import "github.com/yohamta/donburi/features/math"

func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampVector(v math.Vec2, min math.Vec2, max math.Vec2) math.Vec2 {
	return math.Vec2{
		X: Clamp(v.X, min.X, max.X),
		Y: Clamp(v.Y, min.Y, max.Y),
	}
}
