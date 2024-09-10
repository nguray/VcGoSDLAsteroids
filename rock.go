package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Rock struct {
	pos      Vector2f
	veloVec  Vector2f
	mass     float64
	radius   float64
	iExplode int
	fDelete  bool
}

func NewRock(p Vector2f, v Vector2f, m float64) *Rock {
	rck := &Rock{pos: p, veloVec: v, mass: m}
	rck.radius = 10.0 * m
	rck.fDelete = false
	rck.iExplode = 0
	return rck
}

func NewRandomRock() *Rock {

	m := float64(1 + myRand.Intn(2))
	px := myRand.Intn(WIN_WIDTH)
	ri := int(10 * m)
	if px < ri {
		px = ri + 1
	} else if px > (WIN_WIDTH - ri) {
		px = WIN_WIDTH - ri - 1
	}
	py := myRand.Intn(WIN_HEIGHT)
	if py < ri {
		py = ri + 1
	} else if py > (WIN_HEIGHT - ri) {
		py = WIN_HEIGHT - ri - 1
	}

	ra := float64(myRand.Intn(360)) * math.Pi / 180.0
	rck := &Rock{
		pos:     Vector2f{float64(px), float64(py)},
		veloVec: Vector2f{1.35 * math.Cos(ra), 1.35 * math.Sin(ra)},
		mass:    m,
		radius:  10.0 * m,
		fDelete: false,
	}
	rck.iExplode = 0

	return rck
}

func (r *Rock) UpdatePosition() {
	r.pos.AddVector(r.veloVec)
}

func (r Rock) IsDelete() bool {
	return r.fDelete
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

	if r.iExplode == 0 {

		renderer.SetDrawColor(255, 255, 0, 255)
		DrawCircle(renderer, int32(r.pos.x), int32(r.pos.y), int32(r.radius))
		x1 := r.pos.x
		y1 := r.pos.y
		v := r.veloVec
		v.MulScalar(10)
		x2 := x1 + v.x
		y2 := y1 + v.y
		renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

	} else {

		renderer.SetDrawColor(255, 0, 0, 255)

		for a := 0.0; a < 360; a += 60 {
			ra := a * math.Pi / 180
			uv10 := Vector2f{math.Cos(ra), math.Sin(ra)}
			uv11 := uv10
			d := float64(4 * r.iExplode)
			uv10.MulScalar(d)
			uv11.MulScalar(d + 2)
			p1 := r.pos
			p1.AddVector(uv10)
			p2 := r.pos
			p2.AddVector(uv11)
			renderer.DrawLine(int32(p1.x), int32(p1.y), int32(p2.x), int32(p2.y))
		}

	}

}

func (r *Rock) CollideRock(r1 *Rock) {
	//---------------------------------------
	v := r1.pos
	v.SubVector(r.pos)
	d := v.Magnitude()
	if d <= (r.radius + r1.radius) {
		//mt.Print("Collision\n")

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

func (r *Rock) GetPosition() Vector2f {
	return r.pos
}

func (r *Rock) SetPosition(p Vector2f) {
	r.pos = p
}

func (r *Rock) GetVelocity() Vector2f {
	return r.veloVec
}

func (r *Rock) SetVelocity(v Vector2f) {
	r.veloVec = v
}

func (r *Rock) GetMass() float64 {
	return r.mass
}

func (r *Rock) GetRadius() float64 {
	return r.radius
}
