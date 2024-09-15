package main

import "github.com/veandco/go-sdl2/sdl"

type Bullet struct {
	pos      Vector2f
	veloVect Vector2f
	radius   float64
	fDelete  bool
}

func NewBullet(p Vector2f, vel Vector2f) *Bullet {
	return &Bullet{pos: p, veloVect: vel, radius: 1}
}

func (bul *Bullet) SetDelete(f bool) {
	bul.fDelete = f
}

func (bul *Bullet) IsDelete() bool {
	return bul.fDelete
}

func (bul *Bullet) UpdatePosition() {
	bul.pos.AddVector(bul.veloVect)
}

func (bul *Bullet) HitRock(rock *Rock) bool {
	v := bul.pos
	v.SubVector(rock.pos)
	d := v.Magnitude()
	return d < (rock.radius + bul.radius)
}

func (bul *Bullet) Draw(renderer *sdl.Renderer) {
	uv := bul.veloVect.UnitVector()
	tv := uv.NormalVector()
	uv.MulScalar(5)
	tv.MulScalar(2)
	x1 := bul.pos.x
	y1 := bul.pos.y
	x2 := x1 - uv.x + tv.x
	y2 := y1 - uv.y + tv.y
	x3 := x1 - uv.x - tv.x
	y3 := y1 - uv.y - tv.y
	points := []sdl.FPoint{
		{X: float32(x1), Y: float32(y1)},
		{X: float32(x2), Y: float32(y2)},
		{X: float32(x3), Y: float32(y3)},
		{X: float32(x1), Y: float32(y1)},
	}
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawLinesF(points)

}
