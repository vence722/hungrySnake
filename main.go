package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	CellTypeSnakeHead = iota
	CellTypeSnakeBody
	CellTypeFood

	DirectionLeft = iota
	DirectionUp
	DirectionRight
	DirectionDown

	BoardWidth = 16
)

type Cell struct {
	X int
	Y int
	Type int
}

var (
	SnakeHead *Cell
	SnakeBody []*Cell
	Food *Cell

	Score int
	Direction int
	GameSpeed int

	MovedAfterChangeDirection bool // To avoid bug when changing direction multiple times before next tick
	Playing bool // Controls if the game is continuing
)

func init() {
	rand.Seed(time.Now().Unix())

	startNewGame()
}

func startNewGame() {
	SnakeHead = &Cell{ X: 4, Y: 2, Type: CellTypeSnakeHead}
	SnakeBody = []*Cell{
		{ X: 3, Y: 2, Type: CellTypeSnakeBody},
		{ X: 2, Y: 2, Type: CellTypeSnakeBody},
	}

	Score = 0
	Direction = DirectionRight
	GameSpeed = 100
	Playing = true

	generateFood()
}

func generateFood() {
	foodX := int(rand.Int31n(16))
	foodY := int(rand.Int31n(16))

	if foodX == SnakeHead.X && foodY == SnakeHead.Y {
		generateFood()
		return
	}

	for _, cell := range SnakeBody {
		if foodX == cell.X && foodY == cell.Y {
			generateFood()
			return
		}
	}

	Food = &Cell{X: foodX, Y: foodY, Type: CellTypeFood}
}

func main() {
	fmt.Println("start hungry snake")

	js.Global().Set("getGameStatus", js.FuncOf(getGameStatus))
	js.Global().Set("changeDirection", js.FuncOf(changeDirection))
	js.Global().Set("restartGame", js.FuncOf(restartGame))

	for {
		if Playing {
			tick()
		}
		time.Sleep(time.Duration(GameSpeed) * time.Millisecond)
	}
}

func tick() {
	var lastX = SnakeHead.X
	var lastY = SnakeHead.Y

	// update SnakeHead
	switch Direction {
	case DirectionLeft:
		if SnakeHead.X == 0 {
			SnakeHead.X = BoardWidth - 1
		} else {
			SnakeHead.X -= 1
		}
	case DirectionUp:
		if SnakeHead.Y == 0 {
			SnakeHead.Y = BoardWidth - 1
		} else {
			SnakeHead.Y -= 1
		}
	case DirectionRight:
		if SnakeHead.X == BoardWidth - 1 {
			SnakeHead.X = 0
		} else {
			SnakeHead.X += 1
		}
	case DirectionDown:
		if SnakeHead.Y == BoardWidth - 1 {
			SnakeHead.Y = 0
		} else {
			SnakeHead.Y += 1
		}
	}

	// update SnakeBody
	for _, body := range SnakeBody {
		bodylastX := body.X
		bodylastY := body.Y
		body.X = lastX
		body.Y = lastY
		lastX = bodylastX
		lastY = bodylastY
	}

	// if SnakeHead goes to the Food cell
	// add score by 1, and generate food
	// and then add 1 cell in SnakeBody
	if SnakeHead.X == Food.X && SnakeHead.Y == Food.Y {
		Score += 1

		newBody := &Cell{X:lastX, Y:lastY, Type: CellTypeSnakeBody}
		SnakeBody = append(SnakeBody, newBody)

		generateFood()
	}

	// if SnakeHead touches any cell inside SnakeBody
	// stops the game, and notify for the end of the game
	for _, body := range SnakeBody {
		if SnakeHead.X == body.X && SnakeHead.Y == body.Y {
			Playing = false
			break
		}
	}

	MovedAfterChangeDirection = true
}

func getGameStatus(this js.Value, args []js.Value) interface{} {
	gameStatus := map[string]interface{}{
		"score": Score,
		"blocks": []interface{}{},
	}
	gameStatus["blocks"] = append(gameStatus["blocks"].([]interface{}), map[string]interface{}{"x": SnakeHead.X, "y": SnakeHead.Y})
	for _, body := range SnakeBody {
		gameStatus["blocks"] = append(gameStatus["blocks"].([]interface{}), map[string]interface{}{"x": body.X, "y": body.Y})
	}
	gameStatus["blocks"] = append(gameStatus["blocks"].([]interface{}), map[string]interface{}{"x": Food.X, "y": Food.Y})
	return gameStatus
}

func changeDirection(this js.Value, args []js.Value) interface{} {
	if !MovedAfterChangeDirection {
		return nil
	}
	direction := args[0].String()
	switch direction {
	case "up":
		if Direction != DirectionDown {
			Direction = DirectionUp
		}
	case "left":
		if Direction != DirectionRight {
			Direction = DirectionLeft
		}
	case "down":
		if Direction != DirectionUp {
			Direction = DirectionDown
		}
	case "right":
		if Direction != DirectionLeft{
			Direction = DirectionRight
		}
	}
	MovedAfterChangeDirection = false
	return nil
}

func restartGame(this js.Value, args []js.Value) interface{} {
	startNewGame()
	return nil
}