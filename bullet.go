package main

import (
	"sdl2_asteroids/vector"

	"github.com/veandco/go-sdl2/sdl"
)

type Bullet struct {
	pos      vector.Vector2f
	veloVect vector.Vector2f
	radius   float64
	fDelete  bool
}

func NewBullet(p vector.Vector2f, vel vector.Vector2f) *Bullet {
	return &Bullet{pos: p, veloVect: vel, radius: 1}
}

func (bul *Bullet) SetDelete(f bool) {
	bul.fDelete = f
}

func (bul *Bullet) IsDelete() bool {
	return bul.fDelete
}

func (bul *Bullet) UpdatePosition() {
	bul.pos.Add(bul.veloVect)
}

func (bul *Bullet) HitRock(rock *Rock) bool {
	v := bul.pos
	v.Sub(rock.pos)
	d := v.Magnitude()
	return d < (rock.radius + bul.radius)
}

func (bul *Bullet) Draw(renderer *sdl.Renderer) {
	uv := vector.Normalize(bul.veloVect)
	tv := uv.Normal()
	uv.Mul(5)
	tv.Mul(2)
	x1 := bul.pos.X
	y1 := bul.pos.Y
	x2 := x1 - uv.X + tv.X
	y2 := y1 - uv.Y + tv.Y
	x3 := x1 - uv.X - tv.X
	y3 := y1 - uv.Y - tv.Y
	points := []sdl.FPoint{
		{X: float32(x1), Y: float32(y1)},
		{X: float32(x2), Y: float32(y2)},
		{X: float32(x3), Y: float32(y3)},
		{X: float32(x1), Y: float32(y1)},
	}
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawLinesF(points)

}
