package main

import "math"

type Vec2f struct {
	X float64 `json:"dir_x"`
	Y float64 `json:"dir_y"`
}

func makeVec2f(x, y float64) Vec2f {
	return Vec2f{x, y}
}

func (vec Vec2f) magnitude() float64 {
	return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y)
}

func (vec Vec2f) normalize() Vec2f {
	return Vec2f{vec.X / vec.magnitude(), vec.Y / vec.magnitude()}
}

func (vec Vec2f) scale(amount float64) Vec2f {
	return Vec2f{vec.X * amount, vec.Y * amount}
}
