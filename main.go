package main

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	winWidth  = 420
	winHeight = 660
	scl       = 15
)

const (
	up = iota
	right
	down
	left
)

type state struct {
	score            int
	scoreToGive      int
	scoreTick        float32
	scoreColour      rl.Color
	scoreColourTimer float32
	showScoreColour  bool
	timer            float32
	timerCount       float32
	isDead           bool
	framelimit       bool
}

type snake struct {
	body    []rl.Vector2
	prev    rl.Vector2
	dir     int
	nextDir int
}

type fruit struct {
	tick float32
	pos  rl.Vector2
}

type hud struct {
	rec rl.Rectangle
}

func main() {
	rl.InitWindow(winWidth, winHeight, "Snake")
	rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))

	state := state{
		score:            0,
		scoreToGive:      15,
		scoreTick:        0,
		scoreColour:      rl.White,
		scoreColourTimer: 0,
		showScoreColour:  false,
		timer:            0,
		timerCount:       0.3,
		isDead:           false,
		framelimit:       false,
	}

	snake := snake{
		body:    make([]rl.Vector2, 1),
		prev:    rl.NewVector2(0, 0),
		dir:     up,
		nextDir: up,
	}

	fruit := fruit{
		tick: 0,
		pos:  rl.NewVector2(float32(rl.GetRandomValue(0, (winWidth/scl)-1)*scl), float32(rl.GetRandomValue(2, (winHeight/scl)-1)*scl)),
	}

	hud := hud{
		rec: rl.NewRectangle(0, 0, winWidth, 30),
	}

	restart(&snake, &state, &fruit)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		if rl.IsKeyPressed(rl.KeyF) && state.framelimit {
			rl.SetTargetFPS(int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor())))
			state.framelimit = false
		} else if rl.IsKeyPressed(rl.KeyF) && !state.framelimit {
			rl.SetTargetFPS(0)
			state.framelimit = true
		}

		snakeUpdate(&state, &snake, dt, &fruit, &hud)
		fruitUpdate(&fruit, &snake, &state, dt)

		rl.BeginDrawing()

		rl.ClearBackground(rl.Gray)

		draw(&snake, &fruit, &hud, &state)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func snakeUpdate(state *state, snake *snake, dt float32, fruit *fruit, hud *hud) {
	if rl.IsKeyPressed(rl.KeyR) {
		restart(snake, state, fruit)
	}

	if state.isDead {
		return
	}

	state.timer += dt

	if (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp)) && snake.dir != down {
		snake.nextDir = up
	} else if (rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown)) && snake.dir != up {
		snake.nextDir = down
	} else if (rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && snake.dir != left {
		snake.nextDir = right
	} else if (rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && snake.dir != right {
		snake.nextDir = left
	}

	if snake.body[0].X > winWidth-scl {
		snake.body[0].X = 0
	} else if snake.body[0].X < 0 {
		snake.body[0].X = winWidth - scl
	}
	if snake.body[0].Y > winHeight-scl {
		snake.body[0].Y = hud.rec.Height
	} else if snake.body[0].Y < hud.rec.Height {
		snake.body[0].Y = winHeight - scl
	}

	snake.prev = snake.body[0]

	if state.timerCount <= 0.05 {
		state.timerCount = 0.05
	}

	if state.timer >= state.timerCount {
		state.timer = 0
		snake.dir = snake.nextDir

		switch snake.dir {
		case up:
			snake.body[0].Y -= scl
		case down:
			snake.body[0].Y += scl
		case right:
			snake.body[0].X += scl
		case left:
			snake.body[0].X -= scl
		}

		for i := range snake.body {
			if i == 0 {
				continue
			}

			if rl.CheckCollisionRecs(rl.NewRectangle(snake.body[0].X, snake.body[0].Y, scl, scl), rl.NewRectangle(snake.body[i].X, snake.body[i].Y, scl, scl)) {
				state.isDead = true
			}

			cur := snake.body[i]
			snake.body[i] = snake.prev
			snake.prev = cur
		}
	}
}

func fruitUpdate(fruit *fruit, snake *snake, state *state, dt float32) {
	if state.isDead {
		return
	}

	fruit.tick += dt
	state.scoreTick += dt
	if rl.CheckCollisionRecs(rl.NewRectangle(snake.body[0].X, snake.body[0].Y, scl, scl), rl.NewRectangle(fruit.pos.X, fruit.pos.Y, scl, scl)) {
		state.scoreColour = rl.Lime
		state.showScoreColour = true
		snake.body = append(snake.body, rl.NewVector2(snake.body[len(snake.body)-1].X, snake.body[len(snake.body)-1].Y))
		fruit.pos = rl.NewVector2(float32(rl.GetRandomValue(0, (winWidth/scl)-1)*scl), float32(rl.GetRandomValue(2, (winHeight/scl)-1)*scl))
		state.timerCount -= 0.05
		fruit.tick = 0
		state.score += state.scoreToGive
		state.scoreToGive = 15
	} else if state.showScoreColour {
		state.scoreColourTimer += dt
	}
	if state.scoreColourTimer >= 0.35 {
		state.showScoreColour = false
		state.scoreColourTimer = 0
		state.scoreColour = rl.White
	}
	if fruit.tick >= 15 {
		fruit.tick = 0
		fruit.pos = rl.NewVector2(float32(rl.GetRandomValue(0, winWidth/scl)*scl), float32(rl.GetRandomValue(2, winHeight/scl)*scl))
		state.scoreToGive = 15
	}
	if state.scoreTick >= 1 {
		state.scoreTick = 0
		state.scoreToGive -= 1
	}
}

func draw(snake *snake, fruit *fruit, hud *hud, state *state) {
	for i := range snake.body {
		rl.DrawRectangleV(snake.body[0], rl.NewVector2(scl, scl), rl.Red)
		if i == 0 {
			continue
		}
		rl.DrawRectangleV(snake.body[i], rl.NewVector2(scl, scl), rl.Maroon)
	}
	rl.DrawRectangleV(fruit.pos, rl.NewVector2(scl, scl), rl.DarkGreen)
	rl.DrawRectangleRec(hud.rec, rl.DarkGray)
	rl.DrawText(fmt.Sprintf("Score: %d", state.score), 10, 2, 30, state.scoreColour)
	rl.DrawText(fmt.Sprintf("+%d", state.scoreToGive), winWidth/2-20, 2, 30, state.scoreColour)
	rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), winWidth-130, 2, 30, rl.White)
}

func restart(snake *snake, state *state, fruit *fruit) {
	snake.body = make([]rl.Vector2, 1)
	snake.nextDir = up
	snake.body[0].X = (winWidth / 2) - scl
	snake.body[0].Y = (winHeight / 2) - scl

	state.timer = 0
	state.timerCount = 0.3
	state.score = 0
	state.isDead = false
	state.scoreToGive = 15
	state.scoreTick = 0
	state.scoreColourTimer = 0
	state.showScoreColour = false

	fruit.pos = rl.NewVector2(float32(rl.GetRandomValue(0, (winWidth/scl)-1)*scl), float32(rl.GetRandomValue(2, (winHeight/scl)-1)*scl))
	fruit.tick = 0
}
