package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

const (
	width   = 20 // Grid width
	height  = 10 // Grid height
	scale   = 20 // Pixel size per grid cell
)

type Point struct {
	x int
	y int
}

type Snake struct {
	body      []Point
	direction Point
	nextMove  Point // Stores the next direction based on player input
}

type Game struct {
	snake    Snake
	food     Point
	score    int
	winSpots int
	gameOver bool
	moved    bool // Tracks if the snake moved in the current frame
}

func (g *Game) spawnFood(r *rand.Rand) {
	for {
		g.food = Point{
			x: r.Intn(width),
			y: r.Intn(height),
		}

		occupied := false
		for _, part := range g.snake.body {
			if part == g.food {
				occupied = true
				break
			}
		}
		if !occupied {
			break
		}
	}
}

func (s *Snake) move() {
	// Move the body of the snake
	for i := len(s.body) - 1; i > 0; i-- {
		s.body[i] = s.body[i-1]
	}
	// Move the head in the set direction
	s.body[0].x += s.direction.x
	s.body[0].y += s.direction.y
}

func (s *Snake) grow() {
	s.body = append(s.body, s.body[len(s.body)-1])
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	// Process input and update the snake's next direction
	if !g.moved { // Only allow direction change once per move
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.snake.direction != (Point{0, 1}) {
			g.snake.nextMove = Point{0, -1}
			g.moved = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.snake.direction != (Point{0, -1}) {
			g.snake.nextMove = Point{0, 1}
			g.moved = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.snake.direction != (Point{1, 0}) {
			g.snake.nextMove = Point{-1, 0}
			g.moved = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.snake.direction != (Point{-1, 0}) {
			g.snake.nextMove = Point{1, 0}
			g.moved = true
		}
	}

	// Apply the movement only if a valid direction is set
	if g.moved {
		// Update the snake's direction
		g.snake.direction = g.snake.nextMove

		// Move the snake
		g.snake.move()

		// Check for collisions with walls
		head := g.snake.body[0]
		if head.x < 0 || head.x >= width || head.y < 0 || head.y >= height {
			g.gameOver = true
			return nil
		}

		// Check for collisions with itself
		for _, part := range g.snake.body[1:] {
			if part == head {
				g.gameOver = true
				return nil
			}
		}

		// Check for eating food
		if head == g.food {
			g.snake.grow()
			g.spawnFood(rand.New(rand.NewSource(time.Now().UnixNano())))
			g.score++
		}

		// Check for win condition
		if g.checkWin() {
			fmt.Println("You Win!")
			g.gameOver = true
		}

		// Reset the move flag
		g.moved = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Clear screen with black

	// Draw food
	foodRect := ebiten.NewImage(scale, scale)
	foodRect.Fill(color.RGBA{255, 0, 0, 255}) // Red for food
	geoM := ebiten.GeoM{}
	geoM.Translate(float64(g.food.x*scale), float64(g.food.y*scale))
	screen.DrawImage(foodRect, &ebiten.DrawImageOptions{GeoM: geoM})

	// Draw snake
	for _, part := range g.snake.body {
		snakeRect := ebiten.NewImage(scale, scale)
		snakeRect.Fill(color.RGBA{0, 255, 0, 255}) // Green for snake
		geoM := ebiten.GeoM{}
		geoM.Translate(float64(part.x*scale), float64(part.y*scale))
		screen.DrawImage(snakeRect, &ebiten.DrawImageOptions{GeoM: geoM})
	}

	// Display score
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))
	if g.gameOver {
		ebitenutil.DebugPrint(screen, "\nGame Over! Press R to restart.")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return width * scale, height * scale
}

func (g *Game) checkWin() bool {
	emptySpots := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			occupied := false
			for _, part := range g.snake.body {
				if part.x == x && part.y == y {
					occupied = true
					break
				}
			}
			if !occupied {
				emptySpots++
			}
		}
	}
	return emptySpots <= g.winSpots
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	game := &Game{
		snake: Snake{
			body: []Point{{width / 2, height / 2}},
		},
		winSpots: 40, // Default win condition
	}
	game.spawnFood(r)

	ebiten.SetWindowSize(width*scale, height*scale)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Error running game: %v", err)
	}
}
