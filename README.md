# Snake Go Game

**Explanation of the Code**

```bash
go mod init go.game/m
# Auto install
go mod tidy
# Clean mod cache
go clean -modcache

# Manual install
go get github.com/eiannone/keyboard
go get github.com/fatih/color

# Install Ebiten
# go get github.com/hajimehoshi/ebiten/v2
# go mod tidy
```

### Rules of the Snake Game

The game you're playing is a simple Snake Game. The rules are as follows:

1. Objective:
The goal is to control the snake to eat food that appears on the screen. Each time the snake eats food, it grows longer.
The game continues until the snake collides with itself or the walls.
2. Movement:
The snake is controlled using the arrow keys:
Up Arrow: Moves the snake up.
Down Arrow: Moves the snake down.
Left Arrow: Moves the snake left.
Right Arrow: Moves the snake right.
The snake moves continuously in the direction of the last arrow key pressed.
3. Game Over Conditions:
Wall Collision: If the snake's head moves outside the game board (off the grid), the game is over.
Self Collision: If the snake's head collides with any part of its body, the game is over.
4. Eating Food:
Food: The snake eats food represented by an asterisk (*) on the game board.
Each time the snake eats food, it grows by one segment.
A new piece of food will appear at a random location after the snake eats it.
5. Escaping:
Pressing the Esc key will immediately end the game.
6. Winning:
There is no explicit "win" condition. The goal is to continue playing for as long as possible without colliding with the walls or the snake's own body.
7. Game Restart:
The game will display "Game Over!" when you lose, and you can restart the game by running it again.
The game is designed to be simple but can be modified to add more features like increasing speed, adding a score, or more complex levels. Let me know if you'd like to enhance or adjust the game further!

### Game Mechanics:

* The Game struct holds the snake and food, and the Snake struct holds its body segments and direction.
* The snake moves by shifting each body segment to the position of the segment ahead.
* If the snake head collides with the food, the grow function is called, extending the snake's body by one segment.

### Game Loop and Input:

* The game updates continuously in a loop. The update function checks for collisions with the walls, itself, and food.
* Input is handled in a goroutine that listens for "up", "down", "left", and "right" keys for movement.

### Screen Rendering:

The render function creates a simple grid, places the food and snake parts, and then prints the grid to the screen.

### Issues Addressed

You're right! The food might be hidden under the snake's body if it spawns at the same location as one of the snake's segments, which isn't ideal. We'll fix that by ensuring the food is placed in an empty spot on the board.

Additionally, I love your suggestion to add a win condition based on the available free space for food placement. We can check if there are any empty spots on the board to place the food. If no free spots remain, the player wins the game!

#### Here’s how we can address both issues:

* Ensure Food is Visible:

We'll modify the spawnFood method to place food only on an empty spot. This ensures it doesn't overlap with the snake.

* Add Win Condition:

The game will check if there are any free spots available on the board to place food. If all spots are occupied by the snake or the border, the game will declare the player as the winner.

#### Changes Made:

* Ensuring Food is Visible:

The spawnFood method now checks if the randomly chosen position is occupied by the snake. If it is, it will try again until it finds an empty spot.

* Win Condition:

A new method checkWin is added, which checks if there are any empty spots left on the grid for food. If no empty spots are available (i.e., the snake has occupied all possible spaces except the borders), the game will declare the player as the winner.

#### How It Works:

* Food Visibility: The food is guaranteed not to spawn on the snake’s body and will always be visible.
* Win Condition: If the snake fills up the entire board and there are no free spots left, the game declares the player as the winner.

This makes the game fairer and gives the player the possibility to win! Let me know if you'd like to tweak anything further.

### Run the Game:

* Run the game with go run snake.go.
* Control the snake with "up", "down", "left", and "right" keys to move.

Let me know if you'd like any further modifications or have questions about any part of the code! This example should provide a good base to start exploring Go while adding features or improving performance.

### Empty Spots on the Board

At the beginning of the game, we can calculate the empty spots on the board by subtracting the spots occupied by the snake's initial position and the board's borders from the total board area.

Here's the calculation:

**Board Area**: The board is defined as width x height, so the total area is:
total spots=width×height=20×10=200

**Border Spots**: The borders occupy one row and one column around the edges, so the number of border spots is:
border spots=2×(width+height−2)=2×(20+10−2)=2×28=56

**Initial Snake Spots**: The snake starts with a single point in the center.

**Initial Food Spots**: The food starts with a single point on the board.

**Empty Spots**: Subtract the border and initial snake spots from the total spots to get the number of empty spots:

empty spots = total spots − border spots − initial snake spots - initial food spots = 200 − 56 − 1 − 1 = 142

So, at the beginning, there are **142 empty spots** inside the playable area.
