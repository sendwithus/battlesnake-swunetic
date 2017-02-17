package main

import (
	"math/rand"
)

// NOTE: maybe split into multiple files if this gets too big

func NearestFoodHeuristic(gameState *GameState) WeightedDirections {

	var closestFood *Vector
	var food Point

	snake := gameState.MySnake()
	head := snake.Coords[0]
	for _, p := range gameState.Food {
		test := getDistanceBetween(head, p)
		if closestFood == nil {
			closestFood = test
			food = p
		} else if test.Magnitude() < closestFood.Magnitude() {
			closestFood = test
			food = p
		}
	}

	if closestFood == nil {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	if head.Left().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(LEFT)), snake.Id) {
		return []WeightedDirection{{Direction: LEFT, Weight: 50}}
	}
	if head.Right().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(RIGHT)), snake.Id) {
		return []WeightedDirection{{Direction: RIGHT, Weight: 50}}
	}
	if head.Up().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(UP)), snake.Id) {
		return []WeightedDirection{{Direction: UP, Weight: 50}}
	}
	if head.Down().isCloser(&head, &food) && !gameState.IsSolid(head.Add(directionVector(DOWN)), snake.Id) {
		return []WeightedDirection{{Direction: DOWN, Weight: 50}}
	}

	return []WeightedDirection{{Direction: NOOP, Weight: 0}}
}

func GoStraightHeuristic(gameState *GameState) WeightedDirections {

	mySnake := gameState.MySnake()

	if len(mySnake.Coords) <= 1 {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	head := mySnake.Coords[0]
	neck := mySnake.Coords[1]
	directionOfMovement := Point{
		X: head.X - neck.X,
		Y: head.Y - neck.Y,
	}
	allDirections := []string{UP, DOWN, LEFT, RIGHT}

	// try and go straight
	for _, direction := range allDirections {
		if directionOfMovement.Equals(directionVector(direction)) {
			possibleNewHead := head.Add(directionOfMovement)
			if !gameState.IsSolid(possibleNewHead, mySnake.Id) {
				return []WeightedDirection{{Direction: direction, Weight: 50}}
			}
		}
	}
	return []WeightedDirection{{Direction: NOOP, Weight: 0}}
}

func RandomHeuristic(gameState *GameState) WeightedDirections {

	mySnake := gameState.MySnake()
	head := mySnake.Coords[0]
	allDirections := []string{UP, DOWN, LEFT, RIGHT}

	validDirections := []string{}
	for _, direction := range allDirections {
		directionOfMovement := directionVector(direction)
		possibleNewHead := head.Add(directionOfMovement)
		if !gameState.IsSolid(possibleNewHead, mySnake.Id) {
			validDirections = append(validDirections, direction)
		}
	}

	if len(validDirections) == 0 {
		return []WeightedDirection{{Direction: NOOP, Weight: 0}}
	}

	i := rand.Int() % len(validDirections)
	return []WeightedDirection{{Direction: validDirections[i], Weight: 50}}
}
