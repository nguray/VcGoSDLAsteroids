package main

import "math"

type Ship struct {
	x         float64
	y         float64
	a         float64
	s         float64
	veloVec   Vector2f
	unitVec   Vector2f
	normalVec Vector2f
}

func ShipNew(x, y, a, s float64) *Ship {
	//--
	v := Vector2f{s * math.Cos(a*180.0/math.Pi), s * math.Sin(a*180.0/math.Pi)}
	return &Ship{x, y, a, s, v, v.UnitVector(), v.NormalVector()}
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
