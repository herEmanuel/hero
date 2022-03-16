package main

const (
	move = iota
	shoot
)

type Event struct {
	Type      int `json:"type"`
	EventType int `json:"event_type"`

	XOffset int `json:"x_offset"`
	YOffset int `json:"y_offset"`

	playerId int
}
