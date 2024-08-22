package main

import "math"

type Vector2f struct {
	x float64
	y float64
}

func (vec *Vector2f) AddVector(v Vector2f) {
	vec.x += v.x
	vec.y += v.y
}

func (vec *Vector2f) SubVector(v Vector2f) {
	vec.x -= v.x
	vec.y -= v.y
}

func (vec *Vector2f) AddScalar(v float64) {
	vec.x += v
	vec.y += v
}

func (vec *Vector2f) MulScalar(v float64) {
	vec.x *= v
	vec.y *= v
}

func (vec *Vector2f) DivScalar(v float64) {
	if v != 0.0 {
		vec.x /= v
		vec.y /= v
	}
}

func (vec *Vector2f) Dot(v Vector2f) float64 {
	return vec.x*v.x + vec.y*v.y
}

func (vec *Vector2f) Magnitude() float64 {
	return math.Sqrt(vec.x*vec.x + vec.y*vec.y)
}

func (vec *Vector2f) UnitVector() Vector2f {
	m := vec.Magnitude()
	if m > 0.0 {
		return Vector2f{vec.x / m, vec.y / m}
	} else {
		return Vector2f{0.0, 0.0}
	}
}

func (vec *Vector2f) NormalVector() Vector2f {
	return Vector2f{-vec.y, vec.x}
}
