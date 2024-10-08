//
//

package vector

import "math"

type Vector2f struct {
	X float64
	Y float64
}

func Add(lhs, rhs Vector2f) Vector2f {
	return Vector2f{lhs.X + rhs.X, lhs.Y + rhs.Y}
}

func (vec *Vector2f) Add(v Vector2f) {
	vec.X += v.X
	vec.Y += v.Y
}

func Sub(lhs, rhs Vector2f) Vector2f {
	return Vector2f{lhs.X - rhs.X, lhs.Y - rhs.Y}
}

func (vec *Vector2f) Sub(v Vector2f) {
	vec.X -= v.X
	vec.Y -= v.Y
}

func (vec *Vector2f) Mul(v float64) {
	vec.X *= v
	vec.Y *= v
}

func Mul(v Vector2f, fval float64) Vector2f {
	return Vector2f{v.X * fval, v.Y * fval}
}

func (vec *Vector2f) Div(v float64) {
	if v != 0.0 {
		vec.X /= v
		vec.Y /= v
	}
}

func (vec *Vector2f) Dot(v Vector2f) float64 {
	// Produit Scalaire
	return vec.X*v.X + vec.Y*v.Y
}

func (vec *Vector2f) Magnitude() float64 {
	return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y)
}

func (vec *Vector2f) Normalize() Vector2f {
	m := vec.Magnitude()
	if m > 0.0 {
		return Vector2f{vec.X / m, vec.Y / m}
	} else {
		return Vector2f{0.0, 0.0}
	}
}

func Normalize(v Vector2f) Vector2f {
	m := v.Magnitude()
	if m > 0.0 {
		return Vector2f{v.X / m, v.Y / m}
	} else {
		return Vector2f{0.0, 0.0}
	}
}

func (vec *Vector2f) Normal() Vector2f {
	return Vector2f{-vec.Y, vec.X}
}
