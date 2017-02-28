package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func NewGameState(mr MoveRequest) GameState {

	snakes := []*Snake{}
	for _, snake := range mr.Snakes {
		body := []Point{}
		for _, coord := range snake.Coords {
			part := Point{X: coord[0], Y: coord[1]}
			body = append(body, part)
		}
		snakes = append(snakes, &Snake{
			Coords:       body,
			HealthPoints: snake.HealthPoints,
			Id:           snake.Id,
			Name:         snake.Name,
			Taunt:        snake.Taunt,
		})
	}

	foods := []Point{}
	for _, coord := range mr.Food {
		food := Point{X: coord[0], Y: coord[1]}
		foods = append(foods, food)

	}

	gameState := GameState{
		Snakes:  snakes,
		Width:   mr.Width,
		Height:  mr.Height,
		Turn:    mr.Turn,
		Food:    foods,
		winners: []SnakeAI{},
		state:   "running",
		You:     mr.You,
	}

	gameState.generateAStar()
	return gameState
}

func (gameState *GameState) generateAStar() {
	gameState.aStar = map[string]*AStar{}

	start := time.Now()
	for _, snake := range gameState.Snakes {
		gameState.aStar[snake.Id] = NewAStar(gameState, snake.Head())
	}
	if time.Since(start) > 500*time.Millisecond {
		fmt.Printf("Calculated snake A*'s in %v\n", time.Since(start))
	}

}

func (gameState *GameState) MySnake() *Snake {
	snake := gameState.GetSnake(gameState.You)
	return snake
}

func (gameState *GameState) OtherSnakes() []*Snake {
	snakes := []*Snake{}
	for _, snake := range gameState.Snakes {
		if snake.Id != gameState.You {
			snakes = append(snakes, snake)
		}
	}
	return snakes
}

func (gameState *GameState) CountSurroundingWalls(p *Point) int {
	solids := 0
	for _, neighbour := range p.NeighboursWithDiagonals() {
		if gameState.IsSolid(neighbour, "") {
			solids += 1
		}
	}
	return solids
}

func (gameState *GameState) String() string {
	b, err := json.Marshal(gameState)
	if err != nil {
		log.Fatalf("%v", err)
		panic("cant string")
	}
	return string(b)
}

func (gameState *GameState) NextGameState() *GameState {
	nextGameState := GameState{
		Turn:     gameState.Turn + 1,
		SnakeAIs: gameState.SnakeAIs,
		Snakes:   gameState.Snakes,
		Width:    gameState.Width,
		Height:   gameState.Height,
		Food:     gameState.Food,
		winners:  gameState.winners,
		state:    gameState.state,
		You:      gameState.You,
	}

	// get all moves
	moveDirections := map[string]string{}
	for i, snakeAI := range gameState.SnakeAIs {
		snakeId := gameState.Snakes[i].Id
		gameState.You = snakeId
		moveDirections[snakeId] = snakeAI.Move(gameState)
	}

	// extend all snakes
	newHeads := map[string]Point{}
	for _, snake := range nextGameState.Snakes {
		direction := moveDirections[snake.Id]
		newHead := snake.Extend(direction)
		newHeads[snake.Id] = newHead
	}

	// eat or shrink
	foodEaten := []Point{}
	for snakeId, newHead := range newHeads {
		snake := nextGameState.GetSnake(snakeId)
		if nextGameState.FoodAt(&newHead) {
			foodEaten = append(foodEaten, newHead)
			snake.HealthPoints = 100
		} else {
			snake.HealthPoints -= 1
			snake.Coords = snake.Coords[0 : len(snake.Coords)-1]
		}
		if snake.HealthPoints <= 0 {
			nextGameState.KillSnake(snakeId)
		}
	}

	nextGameState.UpdateFood(foodEaten)

	// collision
	for snakeId, newHead := range newHeads {
		if gameState.IsSolid(&newHead, snakeId) {
			nextGameState.KillSnake(snakeId)
		}
	}

	nextGameState.CheckWinStates()
	nextGameState.generateAStar()

	return &nextGameState
}

func (gameState *GameState) UpdateFood(foodEaten []Point) {

	// remove food
	for _, eatenFood := range foodEaten {
		newFoodList := []Point{}
		for _, food := range gameState.Food {
			if !eatenFood.Equals(food) {
				newFoodList = append(newFoodList, food)
			}
		}
		gameState.Food = newFoodList
	}

	// spawn food
	for len(gameState.Food) < len(gameState.Food) {
		gameState.SpawnFood()
	}
}

func (gameState *GameState) CheckWinStates() {
	numSnakes := len(gameState.Snakes)
	if numSnakes == 1 {
		gameState.state = "Won"
		gameState.winners = []SnakeAI{
			gameState.GetHeuristicSnake(gameState.Snakes[0].Id),
		}
	}
	if numSnakes == 0 {
		gameState.state = "Draw"
		gameState.winners = []SnakeAI{}
		for _, snake := range gameState.Snakes {
			heuristicSnake := gameState.GetHeuristicSnake(snake.Id)
			gameState.winners = append(gameState.winners, heuristicSnake)
		}
	}
}

func (gameState *GameState) SpawnFood() {
	emptyPoints := []Point{}
	for x := 0; x < gameState.Width; x += 1 {
		for y := 0; y < gameState.Height; y += 1 {
			p := Point{X: x, Y: y}
			if !gameState.IsSolid(&p, "") && !gameState.FoodAt(&p) {
				emptyPoints = append(emptyPoints, p)
			}
		}
	}

	l := len(emptyPoints)
	if l == 0 {
		return
	}
	p := emptyPoints[rand.Intn(l)]
	gameState.Food = append(gameState.Food, p)
}

func (gameState *GameState) FoodAt(p *Point) bool {
	for _, food := range gameState.Food {
		if food.X == p.X && food.Y == p.Y {
			return true
		}
	}
	return false
}

func (gameState *GameState) IsEmpty(point *Point) bool {
	return !gameState.IsSolid(point, "")
}

func (gameState *GameState) ShortestPathsToFood(from *Point) [][]*Point {
	paths := [][]*Point{}
	for _, food := range gameState.Food {
		println("finding shortest path to food")
		path := gameState.shortestPathTo(from, &food, []*Point{})
		println("found shortest path to food")
		paths = append(paths, path)
	}
	return paths
}

func (gameState *GameState) shortestPathTo(from *Point, to *Point, haveVisited []*Point) []*Point {

	if from.Equals(*to) {
		return []*Point{}
	}

	options := [][]*Point{}
	haveVisited = append(haveVisited, from)
	for _, neighbour := range from.Neighbours() {
		if gameState.IsEmpty(neighbour) {
			visited := false
			for _, visitedPoint := range haveVisited {
				if visitedPoint.Equals(*neighbour) {
					visited = true
					break
				}
			}
			if !visited {
				pathFromNeighbour := gameState.shortestPathTo(neighbour, to, haveVisited)
				if pathFromNeighbour != nil {
					options = append(options, pathFromNeighbour)
				}
			}

		}
	}

	var shortestOption []*Point
	for _, option := range options {
		if shortestOption == nil || len(option) < len(shortestOption) {
			shortestOption = option
		}
	}

	return append(shortestOption, from)
}

func (gameState *GameState) IsSolid(point *Point, ignoreSnakeHead string) bool {
	if point.X < 0 || point.X >= gameState.Width {
		return true
	}
	if point.Y < 0 || point.Y >= gameState.Height {
		return true
	}

	// TODO: take in to account tale shrinks (when no food was eaten)
	for _, snake := range gameState.Snakes {
		for i, coord := range snake.Coords {
			if coord.X == point.X && coord.Y == point.Y {
				if !(i == 0 && snake.Id == ignoreSnakeHead) {
					return true
				}
			}
		}
	}
	return false
}

func (gameState *GameState) KillSnake(snakeId string) {
	newSnakes := []*Snake{}
	for _, snake := range gameState.Snakes {
		if snake.Id != snakeId {
			newSnakes = append(newSnakes, snake)
		}
	}
	gameState.Snakes = newSnakes

	for _, snake := range gameState.SnakeAIs {
		if snakeId == snake.GetId() {
			snake.SetDiedOnTurn(gameState.Turn)
			gameState.losers = append(gameState.losers, snake)
			break
		}
	}
}

func (gameState *GameState) MarkSnakeAsWinner(snakeId string) {
	for i, snake := range gameState.SnakeAIs {
		if gameState.Snakes[i].Id == snakeId {
			snake.SetDiedOnTurn(gameState.Turn)
			gameState.winners = append(gameState.winners, snake)
			break
		}
	}
}

func (gameState *GameState) GetSnake(snakeId string) *Snake {
	for _, snake := range gameState.Snakes {
		if snake.Id == snakeId {
			return snake
		}
	}
	return nil
}

func (gameState *GameState) GetHeuristicSnake(snakeId string) SnakeAI {
	for i, _ := range gameState.SnakeAIs {
		if gameState.Snakes[i].Id == snakeId {
			return gameState.SnakeAIs[i]
		}
	}
	return nil
}
