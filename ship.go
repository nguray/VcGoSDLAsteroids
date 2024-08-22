package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Ship struct {
	x         float64
	y         float64
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

func ShipNew(x, y, a, s float64) *Ship {
	//--
	v := Vector2f{s * math.Cos(a*180.0/math.Pi), s * math.Sin(a*180.0/math.Pi)}
	return &Ship{x, y, a, s, v, v.UnitVector(), v.NormalVector(), 0, nil, nil, nil, nil}
}

func (sh *Ship) SetAngle(a float64) {
	sh.a = a
	coeffRadDeg := 180.0 / math.Pi
	sh.veloVec = Vector2f{sh.s * math.Cos(sh.a*coeffRadDeg), sh.s * math.Sin(sh.a*coeffRadDeg)}
	sh.unitVec = sh.veloVec.UnitVector()
	sh.normalVec = sh.veloVec.NormalVector()
}

func (sh *Ship) OffsetAngle(da float64) {
	sh.SetAngle(sh.a + da)
}

func (sh *Ship) Draw(renderer *sdl.Renderer) {

	src := sdl.Rect{X: 0, Y: 0, W: 32, H: 32}
	dst := sdl.Rect{X: int32(ship.x), Y: int32(ship.y), W: 32, H: 32}
	renderer.CopyEx(sh.curTex, &src, &dst, ship.a, nil, sdl.FLIP_NONE)

	left := dst.X
	top := dst.Y
	right := left + dst.W
	bottom := top + dst.H
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawLine(left, top, right, top)
	renderer.DrawLine(right, top, right, bottom)

}
