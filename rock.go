package main

import (
	"math"
	"sdl2_asteroids/vector"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	NB_COSSINS = 6
)

var (
	cos []float64
	sin []float64
)

func PreCalculateCosSin() {
	aOffset := 360 / NB_COSSINS
	for i := range NB_COSSINS {
		ra := float64(i*aOffset) * math.Pi / 180
		cos = append(cos, math.Cos(ra))
		sin = append(sin, math.Sin(ra))
	}
}

type Rock struct {
	pos      vector.Vector2f
	veloVec  vector.Vector2f
	mass     float64
	radius   float64
	fDelete  bool
	iExplode int
	explVecs []vector.Vector2f
	points   []vector.Vector2f
}

func NewRock(p vector.Vector2f, v vector.Vector2f, m float64) *Rock {
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
		pos:     vector.Vector2f{X: float64(px), Y: float64(py)},
		veloVec: vector.Vector2f{X: 1.35 * math.Cos(ra), Y: 1.35 * math.Sin(ra)},
		mass:    m,
		radius:  10.0 * m,
		fDelete: false,
	}
	rck.iExplode = 0

	return rck
}

func (r *Rock) UpdatePosition() {
	r.pos.Add(r.veloVec)
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

func (r *Rock) InitExplosion() {

	for i := range NB_COSSINS {
		d := float64(2)
		v := vector.Mul(vector.Vector2f{X: cos[i], Y: sin[i]}, d)
		r.explVecs = append(r.explVecs, v)
		p1 := vector.Add(r.pos, v)
		r.points = append(r.points, p1)
	}
}

func (r *Rock) UpdateExplosion() {
	for i := range NB_COSSINS {
		v := vector.Mul(r.veloVec, 3)
		p := vector.Add(r.points[i], v)
		p.Add(r.explVecs[i])
		r.points[i] = p
	}
}

func (r *Rock) Draw(renderer *sdl.Renderer) {

	if r.iExplode == 0 {

		renderer.SetDrawColor(255, 255, 0, 255)
		DrawCircle(renderer, int32(r.pos.X), int32(r.pos.Y), int32(r.radius))
		x1 := r.pos.X
		y1 := r.pos.Y
		v := vector.Mul(r.veloVec, 10)
		x2 := x1 + v.X
		y2 := y1 + v.Y
		renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

	} else {

		renderer.SetDrawColor(255, 255, 0, 255)

		for i, p1 := range r.points {
			p2 := p1
			uv := vector.Normalize(r.explVecs[i])
			uv.Mul(float64(r.iExplode))
			p2.Add(uv)
			renderer.DrawLine(int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y))
		}

	}

}

func (r *Rock) CollideRock(r1 *Rock) {
	//---------------------------------------
	v := r1.pos
	v.Sub(r.pos)
	d := v.Magnitude()
	if d <= (r.radius + r1.radius) {
		//mt.Print("Collision\n")

		nV12 := v
		tV12 := nV12.Normal()

		unV12 := vector.Normalize(nV12)
		utV12 := vector.Normalize(tV12)

		nV1 := r.veloVec.Dot(unV12)
		tV1 := r.veloVec.Dot(utV12)
		nV2 := r1.veloVec.Dot(unV12)
		tV2 := r1.veloVec.Dot(utV12)

		sumMass := r.mass + r1.mass
		nV1c := (nV1*(r.mass-r1.mass) + 2*r1.mass*nV2) / sumMass
		nV2c := (nV2*(r1.mass-r.mass) + 2*r.mass*nV1) / sumMass

		//--
		v = unV12
		v.Mul(nV1c)
		r.veloVec = utV12
		r.veloVec.Mul(tV1)
		r.veloVec.Add(v)

		//--
		v = unV12
		v.Mul(nV2c)
		r1.veloVec = utV12
		r1.veloVec.Mul(tV2)
		r1.veloVec.Add(v)

	}

}

func (r *Rock) GetPosition() vector.Vector2f {
	return r.pos
}

func (r *Rock) SetPosition(p vector.Vector2f) {
	r.pos = p
}

func (r *Rock) GetVelocity() vector.Vector2f {
	return r.veloVec
}

func (r *Rock) SetVelocity(v vector.Vector2f) {
	r.veloVec = v
}

func (r *Rock) GetMass() float64 {
	return r.mass
}

func (r *Rock) GetRadius() float64 {
	return r.radius
}
