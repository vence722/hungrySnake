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
)

func init() {
	rand.Seed(time.Now().Unix())

	SnakeHead = &Cell{ X: 4, Y: 2, Type: CellTypeSnakeHead}
	SnakeBody = []*Cell{
		{ X: 3, Y: 2, Type: CellTypeSnakeBody},
		{ X: 2, Y: 2, Type: CellTypeSnakeBody},
	}

	Direction = DirectionDown
	GameSpeed = 100

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

	for {
		tick()
		time.Sleep(time.Duration(GameSpeed) * time.Millisecond)
	}
}

func tick() {
	var lastX, lastY int

	switch Direction {
	case DirectionLeft:
		lastX = SnakeHead.X
		lastY = SnakeHead.Y

		if SnakeHead.X == 0 {
			SnakeHead.X = BoardWidth
		} else {
			SnakeHead.X -= 1
		}

		for _, body := range SnakeBody {
			bodylastX := body.X
			bodylastY := body.Y
			body.X = lastX
			body.Y = lastY
			lastX = bodylastX
			lastY = bodylastY
		}
	case DirectionUp:
		lastX = SnakeHead.X
		lastY = SnakeHead.Y

		if SnakeHead.Y == 0 {
			SnakeHead.Y = BoardWidth
		} else {
			SnakeHead.Y -= 1
		}

		for _, body := range SnakeBody {
			bodylastX := body.X
			bodylastY := body.Y
			body.X = lastX
			body.Y = lastY
			lastX = bodylastX
			lastY = bodylastY
		}
	case DirectionRight:
		lastX = SnakeHead.X
		lastY = SnakeHead.Y

		if SnakeHead.X == BoardWidth - 1 {
			SnakeHead.X = 0
		} else {
			SnakeHead.X += 1
		}

		for _, body := range SnakeBody {
			bodylastX := body.X
			bodylastY := body.Y
			body.X = lastX
			body.Y = lastY
			lastX = bodylastX
			lastY = bodylastY
		}
	case DirectionDown:
		lastX = SnakeHead.X
		lastY = SnakeHead.Y

		if SnakeHead.Y == BoardWidth - 1 {
			SnakeHead.Y = 0
		} else {
			SnakeHead.Y += 1
		}

		for _, body := range SnakeBody {
			bodylastX := body.X
			bodylastY := body.Y
			body.X = lastX
			body.Y = lastY
			lastX = bodylastX
			lastY = bodylastY
		}
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