package main

import "github.com/veandco/go-sdl2/sdl"

type Bullet struct {
	pos      Vector2f
	veloVect Vector2f
	fDelete  bool
}

func NewBullet(p Vector2f, vel Vector2f) *Bullet {
	return &Bullet{pos: p, veloVect: vel}
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
	return d < rock.radius
}

func (bul *Bullet) Draw(renderer *sdl.Renderer) {
	uv := bul.veloVect.UnitVector()
	uv.MulScalar(5.0)
	x1 := bul.pos.x
	y1 := bul.pos.y
	x2 := x1 - uv.x
	y2 := y1 - uv.y
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

}
