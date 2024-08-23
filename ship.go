package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Ship struct {
	pos       Vector2f
	a         float64
	s         float64
	veloVec   Vector2f
	unitVec   Vector2f
	normalVec Vector2f
	state     int32
	curTex    *sdl.Texture
	idleTex   *sdl.Texture
	accelTex  *sdl.Texture
	decelTex  *sdl.Texture
}

func ShipNew(p Vector2f, a, s float64) *Ship {
	//--
	ra := ((a) * math.Pi) / 180.0
	unitVec := Vector2f{math.Cos(ra), math.Sin(ra)}
	v := unitVec
	v.MulScalar(s)
	return &Ship{p, (a), s, v, unitVec, unitVec.NormalVector(), 0, nil, nil, nil, nil}
}

func (sh *Ship) SetAngle(a float64) {
	sh.a = a
	ra := ((a) * math.Pi) / 180.0
	sh.unitVec = Vector2f{math.Cos(ra), math.Sin(ra)}
	sh.veloVec = sh.unitVec
	sh.normalVec = sh.veloVec.NormalVector()
	sh.veloVec.MulScalar(sh.s)

}

func (sh *Ship) OffsetAngle(da float64) {
	sh.SetAngle(sh.a - da)
}

func (sh *Ship) Draw(renderer *sdl.Renderer) {

	src := sdl.Rect{X: 0, Y: 0, W: 32, H: 32}
	x := int32(sh.pos.x) - 15
	y := int32(sh.pos.y) - 15
	dst := sdl.Rect{X: x, Y: y, W: 32, H: 32}
	renderer.CopyEx(sh.curTex, &src, &dst, sh.a+90.0, nil, sdl.FLIP_NONE)

	//--
	renderer.SetDrawColor(255, 0, 0, 255)

	x1 := sh.pos.x
	y1 := sh.pos.y
	renderer.DrawLine(int32(x1-5), int32(y1), int32(x1+5), int32(y1))
	renderer.DrawLine(int32(x1), int32(y1-5), int32(x1), int32(y1+5))

	x2 := x1 + 30.0*sh.unitVec.x
	y2 := y1 + 30.0*sh.unitVec.y
	renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

	x2 = x1 + 30.0*sh.normalVec.x
	y2 = y1 + 30.0*sh.normalVec.y
	renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

}

func (sh *Ship) UpdatePosition() {
	//--
	sh.pos.AddVector(sh.veloVec)
	//fmt.Printf("(%.3f,%.3f)\n", sh.pos.x, sh.pos.y)

}

func (sh *Ship) Accelerate(d float64) {
	sh.s += d
	sh.veloVec = sh.unitVec
	sh.veloVec.MulScalar(sh.s)

}
