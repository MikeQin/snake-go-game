package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// Create a title caser for the default language (Und)
	caser = cases.Title(language.Und)
)

const (
	width  = 20
	height = 10
)

type Point struct {
	x int
	y int
}

type Snake struct {
	body      []Point
	direction Point
}

type Game struct {
	snake     Snake
	food      Point
	score     int
	winSpots  int
}

func (g *Game) spawnFood(r *rand.Rand) {
	for {
		g.food = Point{
			x: r.Intn(width-2) + 1,
			y: r.Intn(height-2) + 1,
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
	for i := len(s.body) - 1; i > 0; i-- {
		s.body[i] = s.body[i-1]
	}
	s.body[0].x += s.direction.x
	s.body[0].y += s.direction.y
}

func (s *Snake) grow() {
	s.body = append(s.body, s.body[len(s.body)-1])
}

func (g *Game) update() bool {
	head := g.snake.body[0]
	newX := head.x + g.snake.direction.x
	newY := head.y + g.snake.direction.y

	if newX < 0 || newX >= width || newY < 0 || newY >= height {
		return false
	}

	g.snake.move()

	for _, part := range g.snake.body[1:] {
		if part == g.snake.body[0] {
			return false
		}
	}

	if g.snake.body[0] == g.food {
		g.snake.grow()
		g.spawnFood(rand.New(rand.NewSource(time.Now().UnixNano())))
		g.score++
	}

	return true
}

func (g *Game) render(level string) {
	clearScreen()

	// Print the score and level
	color.Set(color.FgYellow)
	fmt.Printf("Score: %d  Level: %s\n", g.score, caser.String(level))
	color.Unset()

	// Create the game grid with borders
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] { 
			if i == 0 || i == height-1 {
				grid[i][j] = '='
			} else if  j == 0 || j == width-1 {
				grid[i][j] = '|'
			} else {
				grid[i][j] = ' '
			}
		}
	}

	// Draw food
	if g.food.x >= 0 && g.food.x < width && g.food.y >= 0 && g.food.y < height {
		grid[g.food.y][g.food.x] = '*'
	}

	// Draw snake
	for _, part := range g.snake.body {
		if part.x >= 0 && part.x < width && part.y >= 0 && part.y < height {
			grid[part.y][part.x] = 'O'
		}
	}

	// Print the grid
	for _, row := range grid {
		for _, cell := range row {
			if cell == '*' {
				color.Set(color.FgRed)
				fmt.Print(string(cell))
				color.Unset()
			} else if cell == 'O' {
				color.Set(color.FgGreen)
				fmt.Print(string(cell))
				color.Unset()
			} else {
				fmt.Print(string(cell))
			}
		}
		fmt.Println()
	}
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") 
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		clearScreenLinux()
	}
}

func clearScreenLinux() {
	cmd := exec.Command("clear") 
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (g *Game) checkWin() bool {
	emptySpots := 0
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
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

func printRules() {
	color.Set(color.FgCyan)
	fmt.Println("Welcome to Snake Game!")
	fmt.Println("Rules:")
	fmt.Println("1. Use arrow keys to control the snake:")
	fmt.Println("  ↑ - Move up")
	fmt.Println("  ↓ - Move down")
	fmt.Println("  ← - Move left")
	fmt.Println("  → - Move right")
	fmt.Println("2. Eat the red food (*) to grow and increase your score.")
	fmt.Println("3. Avoid running into walls or your own tail.")
	fmt.Println("4. The game ends if you hit the wall or yourself.")
	fmt.Println("5. You win by reaching the target empty spots left.")
	fmt.Println("Press ESC to quit the game during the game.")
	fmt.Println("Press any key to start the game...")
	color.Unset()
}

func chooseDifficultyLevel() string {
	reader := bufio.NewReader(os.Stdin)
	color.Set(color.FgMagenta)
	fmt.Println("Choose difficulty level\n[1 - easy (default), 2 - medium, 3 - hard]:")
	color.Unset()
	difficulty, _ := reader.ReadString('\n')
	difficulty = strings.TrimSpace(strings.ToLower(difficulty))

	switch difficulty {
	case "1", "2", "3":
		return difficulty
	default:
		fmt.Println("Defaulting to level 1...")
		return "1"
	}
}

func levelToWinSpots(level string) int {
	switch level {
	case "1": // "easy"
		return 80
	case "2": // "medium"
		return 60
	case "3":  // "hard"
		return 40
	default:
		return 80
	}
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	level := chooseDifficultyLevel()
	game := Game{
		snake: Snake{
			body: []Point{
				{width / 2, height / 2},
			},
			direction: Point{0, 0},
		},
		winSpots: levelToWinSpots(level),
	}
	game.spawnFood(r)

	err := keyboard.Open()
	if err != nil {
		fmt.Println("Error initializing keyboard:", err)
		return
	}
	defer keyboard.Close()

	printRules()
	keyboard.GetKey()
	game.render(level)

	for {
		if game.snake.direction != (Point{0, 0}) {
			if !game.update() {
				fmt.Println("Game Over!")
				break
			}
			if game.checkWin() {
				fmt.Println("You Win!")
				break
			}
			game.render(level)
		}

		_, key, _ := keyboard.GetKey()
		switch key {
		case keyboard.KeyArrowUp:
			if game.snake.direction != (Point{0, 1}) {
				game.snake.direction = Point{0, -1}
			}
		case keyboard.KeyArrowDown:
			if game.snake.direction != (Point{0, -1}) {
				game.snake.direction = Point{0, 1}
			}
		case keyboard.KeyArrowLeft:
			if game.snake.direction != (Point{1, 0}) {
				game.snake.direction = Point{-1, 0}
			}
		case keyboard.KeyArrowRight:
			if game.snake.direction != (Point{-1, 0}) {
				game.snake.direction = Point{1, 0}
			}
		case keyboard.KeyEsc:
			fmt.Println("Game Over!")
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}

