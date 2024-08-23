/*--------------------------------------------*\
			Asteroids using sdl2
                 	2024
			Raymond NGUYEN THANH
\*--------------------------------------------*/

package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

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
	//cellSize int32
	//tt_font  *ttf.Font
	surface  *sdl.Surface
	src, dst sdl.Rect

	ship *Ship
)

// func ProcessEventsPlay(renderer *sdl.Renderer) bool {
// 	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// 		switch t := event.(type) {
// 		case *sdl.QuitEvent:
// 			game.fQuitGame = true
// 			return false
// 		case *sdl.KeyboardEvent:

// 			keyCode := t.Keysym.Sym

// 			//keys := ""
// 			if t.State == sdl.PRESSED && t.Repeat == 0 {
// 				switch keyCode {
// 				case sdl.K_p:
// 					game.fPause = !game.fPause
// 				case sdl.K_LEFT:
// 					game.velX = -1
// 					isOutLRBoardLimit = curTetromino.IsOutLeftBoardLimit
// 				case sdl.K_RIGHT:
// 					game.velX = 1
// 					isOutLRBoardLimit = curTetromino.IsOutRightBoardLimit
// 				case sdl.K_UP:
// 					if curTetromino != nil {
// 						curTetromino.RotateLeft()

// 						if curTetromino.HitGround(game.board) {
// 							//-- Undo Rotate
// 							curTetromino.RotateRight()

// 						} else if curTetromino.IsOutRightBoardLimit() {
// 							backupX := curTetromino.x
// 							//-- Move tetromino inside board
// 							for curTetromino.IsOutRightBoardLimit() {
// 								curTetromino.x--
// 							}
// 							if curTetromino.HitGround(game.board) {
// 								curTetromino.x = backupX
// 								//-- Undo Rotate
// 								curTetromino.RotateRight()
// 							}

// 						} else if curTetromino.IsOutLeftBoardLimit() {

// 							backupX := curTetromino.x
// 							//-- Move tetromino inside board
// 							for curTetromino.IsOutLeftBoardLimit() {
// 								curTetromino.x++
// 							}
// 							if curTetromino.HitGround(game.board) {
// 								curTetromino.x = backupX
// 								//-- Undo Rotate
// 								curTetromino.RotateRight()
// 							}

// 						}

// 					}
// 				case sdl.K_DOWN:
// 					game.fFastDown = true
// 				case sdl.K_SPACE:
// 					if curTetromino != nil {
// 						//-- Drop current Tetromino
// 						game.fDrop = true
// 					}
// 				case sdl.K_ESCAPE:
// 					return false
// 				}
// 			} else if t.State == sdl.RELEASED {
// 				switch keyCode {
// 				case sdl.K_LEFT:
// 					game.velX = 0
// 				case sdl.K_RIGHT:
// 					game.velX = 0
// 				case sdl.K_DOWN:
// 					game.fFastDown = false
// 				}

// 			}

// 		}
// 	}
// 	return true
// }

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

	v1 := Vector2f{1.5 * math.Cos(20.0), 1.5 * math.Sin(20.0)}
	fmt.Printf("v1(%3.2f,%3.2f)\n", v1.x, v1.y)
	uv1 := v1.UnitVector()
	fmt.Printf("uv1(%3.2f,%3.2f)\n", uv1.x, uv1.y)
	nv1 := uv1.NormalVector()
	fmt.Printf("nv1(%3.2f,%3.2f)\n", nv1.x, nv1.y)

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

	//var rect sdl.Rect
	//var rects []sdl.Rect

	//--
	//startH := time.Now()
	//startV := startH
	//startR := startH

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
		// renderer.FillRects(rects)

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

		//fmt.Printf("iRotate = %d\n", int32(ship.a))

		//------------------------------------------------------------
		//-- Draw Game

		//renderer.Copy(shipTex, &src, &dst)

		ship.Draw(renderer)

		// if surface, err = window.GetSurface(); err == nil {
		// 	shipSprite.BlitScaled(nil, surface, &sdl.Rect{X: 100, Y: 100, W: 32, H: 32})
		// 	window.UpdateSurface()
		// }

		//--
		renderer.Present()

		sdl.Delay(20)

	}

}
