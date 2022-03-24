package main

const (
	moveEvnt = iota
	shootEvnt
	deathEvnt
	respawnEvnt
	joinRoomEvnt
	leaveRoomEvnt
)

type Event interface {
	eventType() int
	playerId() int
}

type GenericEvent struct {
	Type      int `json:"type"`
	EventType int `json:"event_type"`

	XOffset   int   `json:"x_offset"`
	YOffset   int   `json:"y_offset"`
	Direction Vec2f `json:"direction"`

	PlayerId int `json:"player_id"`
}

func (ge GenericEvent) eventType() int {
	return ge.EventType
}

func (ge GenericEvent) playerId() int {
	return ge.PlayerId
}

type DeathEvent struct {
	Type      int `json:"type"`
	EventType int `json:"event_type"`

	PlayerId int `json:"player_id"`
	KillerId int `json:"killer_id"`
}

func (de DeathEvent) eventType() int {
	return de.EventType
}

func (de DeathEvent) playerId() int {
	return de.PlayerId
}

type RespawnOrJoinEvent struct {
	Type      int `json:"type"`
	EventType int `json:"event_type"`

	Player *Player `json:"player"`
}

func (de RespawnOrJoinEvent) eventType() int {
	return de.EventType
}

func (de RespawnOrJoinEvent) playerId() int {
	return de.Player.Id
}
