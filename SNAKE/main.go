package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

type point struct {
	x int
	y int
}

var (
	snake   []point
	food    point
	width   = 20
	height  = 20
	dir     = point{1, 0} // Initial direction (right)
	running = true
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Initialize the game
	initGame()

	// Game loop
	for running {
		render()
		handleInput()
		update()
		time.Sleep(100 * time.Millisecond) // Control game speed
	}
}

func initGame() {
	// Initialize the snake in the middle of the screen
	snake = []point{{width / 2, height / 2}}
	spawnFood()
}

func spawnFood() {
	// Generate random food position
	rand.Seed(time.Now().UnixNano())
	food = point{rand.Intn(width), rand.Intn(height)}
}

func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw the snake
	for _, p := range snake {
		termbox.SetCell(p.x, p.y, 'O', termbox.ColorGreen, termbox.ColorDefault)
	}

	// Draw the food
	termbox.SetCell(food.x, food.y, '*', termbox.ColorRed, termbox.ColorDefault)

	// Draw the borders
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, '#', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, height-1, '#', termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < height; y++ {
		termbox.SetCell(0, y, '#', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(width-1, y, '#', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.Flush()
}

func handleInput() {
	// Non-blocking input handling
	select {
	case ev := <-termbox.PollEvent():
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyArrowUp:
				if dir != (point{0, 1}) {
					dir = point{0, -1}
				}
			case termbox.KeyArrowDown:
				if dir != (point{0, -1}) {
					dir = point{0, 1}
				}
			case termbox.KeyArrowLeft:
				if dir != (point{1, 0}) {
					dir = point{-1, 0}
				}
			case termbox.KeyArrowRight:
				if dir != (point{-1, 0}) {
					dir = point{1, 0}
				}
			case termbox.KeyEsc:
				running = false
			}
		}
	default:
		// No input
	}
}

func update() {
	// Move the snake
	head := snake[len(snake)-1]
	newHead := point{head.x + dir.x, head.y + dir.y}

	// Check for collision with walls
	if newHead.x <= 0 || newHead.x >= width-1 || newHead.y <= 0 || newHead.y >= height-1 {
		running = false
		return
	}

	// Check for collision with itself
	for _, p := range snake {
		if newHead == p {
			running = false
			return
		}
	}

	// Check if food is eaten
	if newHead == food {
		snake = append(snake, newHead)
		spawnFood()
	} else {
		// Move snake forward
		snake = append(snake, newHead)
		snake = snake[1:]
	}
}
	