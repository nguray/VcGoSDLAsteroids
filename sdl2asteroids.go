/*--------------------------------------------*\
			Asteroids using sdl2
                 	2024
			Raymond NGUYEN THANH
\*--------------------------------------------*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sdl2_asteroids/vector"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GameMode int

type GameObject interface {
	SetPosition(p vector.Vector2f)
	GetPosition() vector.Vector2f
	GetVelocity() vector.Vector2f
	SetVelocity(v vector.Vector2f)
	GetMass() float64
	GetRadius() float64
}

const (
	STANDBY GameMode = iota
	PLAY
	GAMEPAUSE
	GAMEOVER
	HIGHSCORES
)

const (
	LEFT       = 10
	TOP        = 10
	NB_ROWS    = 20
	NB_COLUMNS = 12
	WIN_WIDTH  = 800
	WIN_HEIGHT = 600
	TITLE      = "Go SDL2 Asteroids"
)

var (
	//tt_font  *ttf.Font
	//surface *sdl.Surface
	//src, dst sdl.Rect

	ship *Ship

	bullets       []*Bullet
	rocks         []*Rock
	myRand        *rand.Rand
	fPause        bool
	laser_snd     *mix.Chunk
	explosion_snd *mix.Chunk
	joysticks     [16]*sdl.Joystick
)

func NewGame() {

	bullets = bullets[:0]
	//--
	for range 5 {
		r := NewRandomRock()
		rocks = append(rocks, r)
	}
	ship.SetPosition(vector.Vector2f{X: WIN_WIDTH / 2, Y: WIN_HEIGHT / 2})

}

func FireBullet() {

	if fPause {
		fPause = false
	}
	v := vector.Mul(ship.DirectionVec(), 5.0)
	b := NewBullet(ship.pos, v)
	bullets = append(bullets, b)
	laser_snd.Play(-1, 0)

}

func DoCollision(object0, object1 GameObject) {
	//---------------------------------------
	p0 := object0.GetPosition()
	p1 := object1.GetPosition()
	m0 := object0.GetMass()
	m1 := object1.GetMass()
	r0 := object0.GetRadius()
	r1 := object1.GetRadius()
	veloVec0 := object0.GetVelocity()
	veloVec1 := object1.GetVelocity()

	v := p1
	v.Sub(p0)
	d := v.Magnitude()
	if d <= (r0 + r1) {
		//mt.Print("Collision\n")

		nV12 := v
		tV12 := nV12.Normal()

		unV12 := vector.Normalize(nV12)
		utV12 := vector.Normalize(tV12)

		nV1 := veloVec0.Dot(unV12)
		tV1 := veloVec0.Dot(utV12)
		nV2 := veloVec1.Dot(unV12)
		tV2 := veloVec1.Dot(utV12)

		sumMass := m0 + m1
		nV1c := (nV1*(m0-m1) + 2*m1*nV2) / sumMass
		nV2c := (nV2*(m1-m0) + 2*m0*nV1) / sumMass

		//--
		v0 := unV12
		v0.Mul(nV1c)
		newVeloVec0 := utV12
		newVeloVec0.Mul(tV1)
		newVeloVec0.Add(v0)
		object0.SetVelocity(newVeloVec0)

		//--
		v1 := unV12
		v1.Mul(nV2c)
		newVeloVec1 := utV12
		newVeloVec1.Mul(tV2)
		newVeloVec1.Add(v1)
		object1.SetVelocity(newVeloVec1)

	}

}

type GenericBounce interface {
	*Rock | *Ship
	GetRadius() float64
	GetPosition() vector.Vector2f
	GetVelocity() vector.Vector2f
	SetVelocity(v vector.Vector2f)
}

func DoSreenFrameCollison[T GenericBounce](obj T, s sdl.Rect) {
	//---------------------------------------
	radius := obj.GetRadius()
	pos := obj.GetPosition()
	veloVec := obj.GetVelocity()

	left := float64(s.X) + radius
	top := float64(s.Y) + radius
	right := float64(s.X+s.W) - radius
	bottom := float64(s.Y+s.H) - radius

	if pos.X <= float64(left) || pos.X > float64(right) {
		veloVec.X = -veloVec.X
	}

	if pos.Y <= float64(top) || pos.Y > float64(bottom) {
		veloVec.Y = -veloVec.Y
	}
	obj.SetVelocity(veloVec)

}

func SubDivideRock(r Rock, m float64) {

	uv := vector.Normalize(r.veloVec)
	un := uv.Normal()
	normeV := r.veloVec.Magnitude() * 1.5

	v10 := vector.Add(uv, un)
	v10.Mul(10)
	p10 := vector.Add(r.pos, v10)
	uv10 := vector.Normalize(v10)
	uv10.Mul(normeV)
	rocks = append(rocks, NewRock(p10, uv10, m))

	v20 := vector.Sub(uv, un)
	v20.Mul(10)
	p20 := r.pos
	p20.Add(v20)
	uv20 := vector.Normalize(v20)
	uv20.Mul(normeV)
	rocks = append(rocks, NewRock(p20, uv20, m))

	v30 := vector.Sub(un, uv)
	v30.Mul(10)
	p30 := r.pos
	p30.Add(v30)
	uv30 := vector.Normalize(v30)
	uv30.Mul(normeV)
	rocks = append(rocks, NewRock(p30, uv30, m))

	v40 := vector.Add(uv, un)
	v40.Mul(-1)
	v40.Mul(10)
	p40 := r.pos
	p40.Add(v40)
	uv40 := vector.Normalize(v40)
	uv40.Mul(normeV)
	rocks = append(rocks, NewRock(p40, uv40, m))

}

func main() {

	var renderer *sdl.Renderer

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	nbJoysticks := sdl.NumJoysticks()
	//fmt.Printf("nb joysticks = %d\n", nbJoysticks)ock.mass / 2)

	if nbJoysticks != 0 {
		sdl.JoystickEventState(sdl.ENABLE)
	}

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIN_WIDTH, WIN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	ttf.Init()
	defer ttf.Quit()

	curDir, _ := os.Getwd()
	fullPathName := filepath.Join(curDir, "resources", "Plane00.png")
	shipImg0, err := img.Load(fullPathName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load image: %s\n", err)
		return
	}
	defer shipImg0.Free()
	fullPathName = filepath.Join(curDir, "resources", "Plane01.png")
	shipImg1, err := img.Load(fullPathName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load image: %s\n", err)
		return
	}
	defer shipImg1.Free()
	fullPathName = filepath.Join(curDir, "resources", "Plane02.png")
	shipImg2, err := img.Load(fullPathName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load image: %s\n", err)
		return
	}
	defer shipImg2.Free()

	// v1 := Vector2f{1.5 * math.Cos(20.0), 1.5 * math.Sin(20.0)}
	// fmt.Printf("v1(%3.2f,%3.2f)\n", v1.x, v1.y)
	// uv1 := v1.UnitVector()ock.mass / 2)
	// fmt.Printf("uv1(%3.2f,%3.2f)\n", uv1.x, uv1.y)
	// nv1 := uv1.NormalVector()
	// fmt.Printf("nv1(%3.2f,%3.2f)\n", nv1.x, nv1.y)

	mix.OpenAudio(44100, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, 1024)
	fullPathName = filepath.Join(curDir, "resources", "344276__nsstudios__laser3.wav")
	laser_snd, err = mix.LoadWAV(fullPathName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load Sound: %s\n", err)
		panic(err)
	}
	defer laser_snd.Free()

	fullPathName = filepath.Join(curDir, "resources", "asteroid-94614.mp3")
	explosion_snd, err = mix.LoadWAV(fullPathName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Load Sound: %s\n", err)
		panic(err)
	}
	defer explosion_snd.Free()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	//renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}
	defer renderer.Destroy()

	a := -90.0
	ship = ShipNew(vector.Vector2f{X: WIN_WIDTH / 2, Y: WIN_HEIGHT / 2}, a)

	shipTex0, _ := renderer.CreateTextureFromSurface(shipImg0)
	defer shipTex0.Destroy()
	shipTex1, _ := renderer.CreateTextureFromSurface(shipImg1)
	defer shipTex1.Destroy()
	shipTex2, _ := renderer.CreateTextureFromSurface(shipImg2)
	defer shipTex2.Destroy()

	ship.idleTex = shipTex0
	ship.accelTex = shipTex1
	ship.decelTex = shipTex2
	ship.curTex = shipTex0

	myRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	NewGame()

	//--drawObjects
	startExplodeUpdate := time.Now()
	//startV := startH
	//startR := startH

	screenFrame := sdl.Rect{X: 0, Y: 0, W: WIN_WIDTH, H: WIN_HEIGHT}

	// Precalculate Cosinus and Sinus Values for rock Explosion animation
	PreCalculateCosSin()

	iRotate := 0
	iAccel := 0

	fPause = true
	fStep := false
	running := true

	for running {

		//-- Draw Background
		renderer.SetDrawColor(16, 16, 64, 64)
		renderer.Clear()

		// rect = sdl.Rect{X: int32(LEFT), Y: int32(TOP), W: int32(cellSize * NB_COLUMNS), H: int32(cellSize * NB_ROWS)}
		// renderer.SetDrawColor(10, 10, 100, 255)
		// renderer.FillRect(&rect)20

		elapsedExplodeUpdate := time.Since(startExplodeUpdate)

		//-- Process current mode Events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.JoyAxisEvent:

				fmt.Printf("[%d ms] JoyAxis\ttype:%d\twhich:%c\taxis:%d\tvalue:%d\n",
					t.Timestamp, t.Type, t.Which, t.Axis, t.Value)

				switch t.Axis {
				case 1:
					if t.Value < 500 && t.Value > -500 {
						iAccel = 0
					} else if t.Value < 500 {
						iAccel = 1
					} else if t.Value > 500 {
						iAccel = -1
					}
				case 3:

					if t.Value < 500 && t.Value > -500 {
						iRotate = 0
					} else if t.Value < 500 {
						iRotate = -1
					} else if t.Value > 500 {
						iRotate = 1
					}

				}

			case *sdl.JoyBallEvent:
				fmt.Println("Joystick", t.Which, "trackball moved by", t.XRel, t.YRel)
			case *sdl.JoyButtonEvent:
				if t.State == sdl.PRESSED {
					fmt.Println("Joystick", t.Which, "button", t.Button, "pressed")
					if t.Button == 4 || t.Button == 5 {
						FireBullet()
					}
				} else {
					fmt.Println("Joystick", t.Which, "button", t.Button, "released")
				}
			case *sdl.JoyHatEvent:
				position := ""
				switch t.Value {
				case sdl.HAT_LEFTUP:
					position = "top-left"
				case sdl.HAT_UP:
					position = "top"
				case sdl.HAT_RIGHTUP:
					position = "top-right"
				case sdl.HAT_RIGHT:
					position = "right"
				case sdl.HAT_RIGHTDOWN:
					position = "bottom-right"
				case sdl.HAT_DOWN:
					position = "bottom"
				case sdl.HAT_LEFTDOWN:
					position = "bottom-left"
				case sdl.HAT_LEFT:
					position = "left"
				case sdl.HAT_CENTERED:
					position = "center"
				}

				fmt.Println("Joystick", t.Which, "hat", t.Hat, "moved to", position, "position")
			case *sdl.JoyDeviceAddedEvent:
				// Open joystick for use
				joysticks[int(t.Which)] = sdl.JoystickOpen(int(t.Which))
				if joysticks[int(t.Which)] != nil {
					fmt.Println("Joystick", t.Which, "connected")
				}
			case *sdl.JoyDeviceRemovedEvent:
				if joystick := joysticks[int(t.Which)]; joystick != nil {
					joystick.Close()
				}
				fmt.Println("Joystick", t.Which, "disconnected")

			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym

				if t.State == sdl.PRESSED && t.Repeat == 0 {
					switch keyCode {
					case sdl.K_LEFT:
						iRotate = -1
					case sdl.K_RIGHT:
						iRotate = 1
					case sdl.K_UP:
						iAccel = 1
					case sdl.K_DOWN:
						iAccel = -1
					case sdl.K_p:
						fPause = !fPause
					case sdl.K_SPACE:
						FireBullet()
					case sdl.K_ESCAPE:
						return
					}
				} else if t.State == sdl.RELEASED {
					switch keyCode {
					case sdl.K_LEFT:
						iRotate = 0
					case sdl.K_RIGHT:
						iRotate = 0
					case sdl.K_UP:
						iAccel = 0
					case sdl.K_DOWN:
						iAccel = 0
					case sdl.K_s:
						if fPause {
							fStep = true
							fPause = false
						}
					}

				}

			}

		}

		//-- Game Mode Update States

		// rects = []sdl.Rect{{500, 300, 100, 100}, {200, 300, 200, 200}}
		// renderer.SetDrawColor(255, 0, 255, 255)
		// renderer.FillRects(rects)

		if iRotate < 0 {
			ship.OffsetAngle(2)
		} else if iRotate > 0 {
			ship.OffsetAngle(-2)
		}

		if !fPause {

			if iAccel > 0 {
				ship.Accelerate(0.1)
				ship.SetForwardThrush()
			} else if iAccel < 0 {
				ship.Accelerate(-0.1)
				ship.SetBackwardTrush()
			} else {
				ship.SetIdle()
			}

			ship.UpdatePosition()

			// Keep Ship inside) screen
			DoSreenFrameCollison(ship, screenFrame)

			//-- Bullets
			tmpBullets := bullets[:0]
			for _, b := range bullets {

				b.UpdatePosition()
				fHit := false

				//--
				for _, rock := range rocks {

					if rock.iExplode == 0 && b.HitRock(rock) {
						fHit = true
						explosion_snd.Play(-1, 0)

						if rock.mass == 2 {
							rock.fDelete = true
							//-- SubDivide
							SubDivideRock(*rock, rock.mass/3)
							//fPause = true
						} else if rock.mass == 1 {
							rock.fDelete = true
							//-- SubDivide
							SubDivideRock(*rock, rock.mass/2)

							//fPause = true
						} else {
							rock.iExplode = 1
							rock.InitExplosion()
							//fPause = true

						}
						break
					}
				}

				if !fHit {
					//-- Check for out of window frame
					if (b.pos.X < 0) || (b.pos.X > WIN_WIDTH) || (b.pos.Y < 0) || (b.pos.Y > WIN_HEIGHT) {
						//--
						fHit = true
					}
				}

				if !fHit {
					tmpBullets = append(tmpBullets, b)
				}
			}
			bullets = tmpBullets

			//-- Rocks
			tmpRock1 := rocks[:0]
			for _, r := range rocks {
				if !r.IsDelete() {
					r.UpdatePosition()
					DoSreenFrameCollison(r, screenFrame)
					tmpRock1 = append(tmpRock1, r)
				}
			}
			rocks = tmpRock1

			// Do collison Ship<->Rock
			for _, r := range rocks {
				if !r.fDelete && r.iExplode == 0 {
					DoCollision(ship, r)
				}
			}

			// Do collison between rocks
			var r, r1 *Rock
			for i := 0; i < len(rocks); i++ {
				r = rocks[i]
				if !r.fDelete && r.iExplode == 0 {
					for j := i + 1; j < len(rocks); j++ {
						r1 = rocks[j]
						if !r1.fDelete && r1.iExplode == 0 {
							DoCollision(r, r1)
						}
					}
				}
			}

			if elapsedExplodeUpdate.Milliseconds() > 120 || fStep {
				startExplodeUpdate = time.Now()
				for _, r := range rocks {
					if r.iExplode > 0 {
						r.iExplode += 1
						r.UpdateExplosion()
						if r.iExplode > 4 {
							r.fDelete = true
						}
					}
				}
			}

		} else {
			if iAccel != 0 {
				fPause = false
			}
		}

		//------------------------------------------------------------
		//-- Draw Game

		//renderer.Copy(shipTex, &src, &dst)

		ship.Draw(renderer)

		for _, b := range bullets {
			b.Draw(renderer)
		}

		rocksTemp := rocks[:0]
		for _, r := range rocks {
			if !r.IsDelete() {
				r.Draw(renderer)
				rocksTemp = append(rocksTemp, r)
			}
		}
		rocks = rocksTemp

		// if surface, err = window.GetSurface(); err == nil {
		// 	shipSprite.BlitScaled(nil, surface, &sdl.Rect{X: 100, Y: 100, W: 32, H: 32})
		// 	window.UpdateSurface()
		// }

		if fStep {
			fPause = true
			fStep = false
		}

		//--
		renderer.Present()

		if len(rocks) == 0 {
			NewGame()
			fPause = true
			for sdl.PollEvent() != nil {
			}
			sdl.Delay(500)
		}

		//fmt.Printf("nb bullets = %d\n", len(bullets))

		sdl.Delay(20)

	}

}
