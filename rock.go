package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Rock struct {
	pos     Vector2f
	veloVec Vector2f
	mass    float64
	radius  float64
}

func NewRock(p Vector2f, v Vector2f, m float64) *Rock {
	rck := &Rock{pos: p, veloVec: v, mass: m}
	rck.radius = 10.0 * m
	return rck
}

func (r *Rock) UpdatePosition() {
	r.pos.AddVector(r.veloVec)
}

func DrawCircle(renderer *sdl.Renderer, x, y, radius int32) {

	var offsetX int32 = 0
	var offsetY int32 = radius
	d := radius - 1

	for offsetY >= int32(offsetX) {
		renderer.DrawPoint(x+offsetX, y+offsetY)
		renderer.DrawPoint(x+offsetY, y+offsetX)
		renderer.DrawPoint(x-offsetX, y+offsetY)
		renderer.DrawPoint(x-offsetY, y+offsetX)
		renderer.DrawPoint(x+offsetX, y-offsetY)
		renderer.DrawPoint(x+offsetY, y-offsetX)
		renderer.DrawPoint(x-offsetX, y-offsetY)
		renderer.DrawPoint(x-offsetY, y-offsetX)

		if d >= 2*offsetX {
			d -= 2*offsetX + 1
			offsetX += 1
		} else if d < 2*(radius-offsetY) {
			d += 2*offsetY - 1
			offsetY -= 1
		} else {
			d += 2 * (offsetY - offsetX - 1)
			offsetY -= 1
			offsetX += 1
		}

	}

}

func (r *Rock) Draw(renderer *sdl.Renderer) {

	renderer.SetDrawColor(255, 255, 0, 255)
	DrawCircle(renderer, int32(r.pos.x), int32(r.pos.y), int32(r.radius))

}

func (r *Rock) CollideSreenFrame(s sdl.Rect) {
	//---------------------------------------
	left := float64(s.X) + r.radius
	top := float64(s.Y) + r.radius
	right := float64(s.X+s.W) - r.radius
	bottom := float64(s.Y+s.H) - r.radius

	if r.pos.x <= float64(left) || r.pos.x > float64(right) {
		r.veloVec.x = -r.veloVec.x
	}

	if r.pos.y <= float64(top) || r.pos.y > float64(bottom) {
		r.veloVec.y = -r.veloVec.y
	}

}

func (r *Rock) CollideRock(r1 *Rock) {
	//---------------------------------------
	v := r1.pos
	v.SubVector(r.pos)
	d := v.Magnitude()
	if d <= (r.radius + r1.radius) {
		fmt.Print("Collision\n")

		nV12 := v
		tV12 := nV12.NormalVector()

		unV12 := nV12.UnitVector()
		utV12 := tV12.UnitVector()

		nV1 := r.veloVec.Dot(unV12)
		tV1 := r.veloVec.Dot(utV12)
		nV2 := r1.veloVec.Dot(unV12)
		tV2 := r1.veloVec.Dot(utV12)

		sumMass := r.mass + r1.mass
		nV1c := (nV1*(r.mass-r1.mass) + 2*r1.mass*nV2) / sumMass
		nV2c := (nV2*(r1.mass-r.mass) + 2*r.mass*nV1) / sumMass

		//--
		v = unV12
		v.MulScalar(nV1c)
		r.veloVec = utV12
		r.veloVec.MulScalar(tV1)
		r.veloVec.AddVector(v)

		//--
		v = unV12
		v.MulScalar(nV2c)
		r1.veloVec = utV12
		r1.veloVec.MulScalar(tV2)
		r1.veloVec.AddVector(v)

	}

}
