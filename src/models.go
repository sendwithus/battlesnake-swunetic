package main

type GameStartRequest struct {
	GameId string `json:"game_id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type GameStartResponse struct {
	Color   string  `json:"color"`
	HeadUrl *string `json:"head_url,omitempty"`
	Name    string  `json:"name"`
	Taunt   *string `json:"taunt,omitempty"`
}

type MoveRequest struct {
	Board  [][]BoardCell `json:"board"`
	Food   []Point       `json:"food"`
	GameId string        `json:"game_id"`
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Turn   int           `json:"turn"`
	Snakes []Snake       `json:"snakes"`
	You    string        `json:"you"`
}

type MoveResponse struct {
	Move  string  `json:"move"`
	Taunt *string `json:"taunt,omitempty"`
}

type BoardCell struct {
	State string  `json:"state"`
	Snake *string `json:"snake,omitempty"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Coords       []Point `json:"coords"`
	HealthPoints int     `json:"health_points"`
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Taunt        string  `json:"taunt"`
}