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
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GameMode int

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

	bullets []*Bullet
	rocks   []*Rock
	myRand  *rand.Rand
)

func main() {

	var renderer *sdl.Renderer

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

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
	// uv1 := v1.UnitVector()
	// fmt.Printf("uv1(%3.2f,%3.2f)\n", uv1.x, uv1.y)
	// nv1 := uv1.NormalVector()
	// fmt.Printf("nv1(%3.2f,%3.2f)\n", nv1.x, nv1.y)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	//renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		//return 2
		panic(err)
	}
	defer renderer.Destroy()

	a := -90.0
	ship = ShipNew(Vector2f{400.0, 500.0}, a)

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

	for i := 0; i < 10; i++ {
		rocks = append(rocks, NewRandomRock())

	}

	//var rect sdl.Rect
	//var rects []sdl.Rect

	//--
	//startH := time.Now()
	//startV := startH
	//startR := startH

	screenFrame := sdl.Rect{X: 0, Y: 0, W: WIN_WIDTH, H: WIN_HEIGHT}

	iRotate := 0
	iAccel := 0

	running := true
	for running {

		//-- Draw Background
		renderer.SetDrawColor(16, 16, 64, 64)
		renderer.Clear()

		// rect = sdl.Rect{X: int32(LEFT), Y: int32(TOP), W: int32(cellSize * NB_COLUMNS), H: int32(cellSize * NB_ROWS)}
		// renderer.SetDrawColor(10, 10, 100, 255)
		// renderer.FillRect(&rect)

		//-- Process current mode Events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				return
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
					case sdl.K_SPACE:
						v := ship.DirectionVec()
						v.MulScalar(5.0)
						bullets = append(bullets, NewBullet(ship.pos, v))
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
					}

				}

			}

		}

		//running = processEvents(renderer)

		//-- Game Mode Update States

		// rects = []sdl.Rect{{500, 300, 100, 100}, {200, 300, 200, 200}}
		// renderer.SetDrawColor(255, 0, 255, 255)
		// renderer.FillReocks[i]cts(rects)

		if iRotate < 0 {
			ship.OffsetAngle(2.0)
		} else if iRotate > 0 {
			ship.OffsetAngle(-2.0)
		}

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

		// Keep Ship inside screen
		p := ship.pos
		if p.x < 0.0 {
			p.x = WIN_WIDTH
		} else if p.x > WIN_WIDTH {
			p.x = 0.0
		}
		if p.y < 0.0 {
			p.y = WIN_HEIGHT
		} else if p.y > WIN_HEIGHT {
			p.y = 0.0
		}
		ship.pos = p

		//-- Bullets
		for _, b := range bullets {
			b.UpdatePosition()
			//-- Check for out range
			if (b.pos.x < 0) || (b.pos.x > WIN_WIDTH) || (b.pos.y < 0) || (b.pos.y > WIN_HEIGHT) {
				b.SetDelete(true)
			}
		}

		//-- Rocks
		for _, rock := range rocks {
			rock.UpdatePosition()
			rock.CollideSreenFrame(screenFrame)

		}

		var r *Rock
		for i := 0; i < len(rocks); i++ {
			r = rocks[i]
			for j := i + 1; j < len(rocks); j++ {
				r.CollideRock(rocks[j])
			}
		}

		//fmt.Printf("iRotate = %d\n", int32(ship.a))

		//------------------------------------------------------------
		//-- Draw Game

		//renderer.Copy(shipTex, &src, &dst)

		ship.Draw(renderer)

		for _, b := range bullets {
			if !b.fDelete {
				b.Draw(renderer)
			}
		}

		for _, rock := range rocks {
			rock.Draw(renderer)
		}

		// if surface, err = window.GetSurface(); err == nil {
		// 	shipSprite.BlitScaled(nil, surface, &sdl.Rect{X: 100, Y: 100, W: 32, H: 32})
		// 	window.UpdateSurface()
		// }

		//--
		renderer.Present()

		//-- Update Bullets Slices
		tmp := bullets[:0]
		for _, b := range bullets {
			if !b.fDelete {
				tmp = append(tmp, b)
			}
		}
		bullets = tmp

		//fmt.Printf("nb bullets = %d\n", len(bullets))

		sdl.Delay(20)

	}

}
