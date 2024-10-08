package main

import (
	"math"
	"sdl2_asteroids/vector"

	"github.com/veandco/go-sdl2/sdl"
)

type Ship struct {
	pos           vector.Vector2f
	a             float64
	veloVec       vector.Vector2f
	mass          float64
	radius        float64
	thrushUnitVec vector.Vector2f
	curTex        *sdl.Texture
	idleTex       *sdl.Texture
	accelTex      *sdl.Texture
	decelTex      *sdl.Texture
	shieldLevel   float64
}

func ShipNew(p vector.Vector2f, a float64) *Ship {
	//--
	ra := ((a) * math.Pi) / 180.0
	unitVec := vector.Vector2f{X: math.Cos(ra), Y: math.Sin(ra)}
	v := unitVec
	v.Mul(0.1)

	sh := &Ship{pos: p, a: a, veloVec: v, thrushUnitVec: unitVec}
	sh.mass = 2
	sh.shieldLevel = 3
	sh.radius = (8.0 + sh.shieldLevel*1.5) * sh.mass
	return sh
}

func (sh *Ship) DecShieldLevel() {
	if sh.shieldLevel > 0 {
		sh.shieldLevel--
		sh.radius = (8.0 + sh.shieldLevel*1.5) * sh.mass
	}
}

func (sh *Ship) IncShieldLevel() {
	if sh.shieldLevel < 3 {
		sh.shieldLevel++
		sh.radius = (8.0 + sh.shieldLevel*1.5) * sh.mass
	}
}

func (sh *Ship) SetAngle(a float64) {
	sh.a = a
	ra := ((a) * math.Pi) / 180.0
	sh.thrushUnitVec = vector.Vector2f{X: math.Cos(ra), Y: math.Sin(ra)}
}

func (sh *Ship) SetIdle() {
	sh.curTex = sh.idleTex
}

func (sh *Ship) SetForwardThrush() {
	sh.curTex = sh.accelTex
}

func (sh *Ship) SetBackwardTrush() {
	sh.curTex = sh.decelTex
}

func (sh *Ship) OffsetAngle(da float64) {
	sh.SetAngle(sh.a - da)
}

func (sh *Ship) Draw(renderer *sdl.Renderer) {

	//--

	red := (255 - uint8(sh.shieldLevel)*64)
	renderer.SetDrawColor(red, 60, 0, 255)
	DrawCircle(renderer, int32(sh.pos.X), int32(sh.pos.Y), int32(sh.radius))

	//--
	src := sdl.Rect{X: 0, Y: 0, W: 32, H: 32}
	x := int32(sh.pos.X) - 15
	y := int32(sh.pos.Y) - 15
	dst := sdl.Rect{X: x, Y: y, W: 32, H: 32}
	renderer.CopyEx(sh.curTex, &src, &dst, sh.a+90.0, nil, sdl.FLIP_NONE)

	// x1 := sh.pos.x
	// y1 := sh.pos.y
	// renderer.DrawLine(int32(x1-5), int32(y1), int32(x1+5), int32(y1))
	// renderer.DrawLine(int32(x1), int32(y1-5), int32(x1), int32(y1+5))

	// x2 := x1 + 30.0*sh.thrushUnitVec.x
	// y2 := y1 + 30.0*sh.thrushUnitVec.y
	// renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

	// thrushNormalVec := sh.thrushUnitVec.NormalVector()
	// x2 = x1 + 30.0*thrushNormalVec.x
	// y2 = y1 + 30.0*thrushNormalVec.y
	// renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

}

func (sh *Ship) UpdatePosition() {
	//--
	sh.pos.Add(sh.veloVec)
	//fmt.Printf("(%.3f,%.3f)\n", sh.pos.x, sh.pos.y)

}

func (sh *Ship) Accelerate(d float64) {
	v := sh.thrushUnitVec
	v.Mul(d)
	sh.veloVec.Add(v)

}

func (sh *Ship) DirectionVec() vector.Vector2f {
	ra := ((sh.a) * math.Pi) / 180.0
	return vector.Vector2f{X: math.Cos(ra), Y: math.Sin(ra)}
}

func (sh *Ship) GetPosition() vector.Vector2f {
	return sh.pos
}

func (sh *Ship) SetPosition(p vector.Vector2f) {
	sh.pos = p
}

func (sh *Ship) GetVelocity() vector.Vector2f {
	return sh.veloVec
}

func (sh *Ship) SetVelocity(v vector.Vector2f) {
	sh.veloVec = v
}

func (sh *Ship) GetMass() float64 {
	return sh.mass
}

func (sh *Ship) GetRadius() float64 {
	return sh.radius
}
