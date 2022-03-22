package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	maxEvents = 1024
)

/* We just have one non-customizable map, so we're just gonna
   hardcode all of the obstacles in it
*/

var mapObstacles []Obstacle = []Obstacle{
	{530, 217, 65, 65},
	{812, 47, 32, 104},
	{301, 381, 90, 58},
	{1047, 267, 60, 60},
	{944, 507, 69, 84},
	{1248, 557, 104, 32},
	{1257, 679, 75, 75},
	{20, 636, 81, 81},
	{91, 692, 73, 73},
	{640, 825, 81, 81},
	{1997, 24, 65, 65},
	{2423, 467, 60, 60},
	{2423, 540, 60, 60},
	{1248, 1217, 104, 32},
	{1320, 1288, 32, 104},
}

type Room struct {
	code    string
	lock    sync.Mutex
	players []*Player
	objects []Object

	chat   chan string
	events chan Event

	closing bool
}

func newRoom(creator *Player) *Room {
	return &Room{
		code:    genNewRoomCode(),
		players: []*Player{creator},
		chat:    make(chan string),
		events:  make(chan Event, maxEvents),
	}
}

func (room *Room) getPlayer(id int) *Player {
	for _, player := range room.players {
		if player.Id == id {
			return player
		}
	}

	return nil
}

func (room *Room) notifyAllPlayers(event Event) {
	for _, player := range room.players {
		player.events <- event
	}
}

func (room *Room) handleEvent() {
	event := <-room.events

	player := room.getPlayer(event.PlayerId)
	if player == nil {
		// prob the event from a player that left the room
		return
	}

	switch event.EventType {
	case move:
		player.processMovement(event.XOffset, event.YOffset, event.Direction)
	case shoot:
		bullet := newBullet(player)
		room.objects = append(room.objects, bullet)
		room.notifyAllPlayers(event)
	case leaveRoom:
		if room.closing {
			for i, r := range globalState.rooms {
				if r == room {
					globalState.lock.Lock()
					globalState.rooms = append(globalState.rooms[:i], globalState.rooms[i+1:]...)
					globalState.lock.Unlock()
				}
			}

			log.Println("Closing room ", room.code)
			return
		}

		room.notifyAllPlayers(event)
	default:
		log.Println("Received invalid event message from player ", event.PlayerId)
		player.disconnect() // TODO: fix this
	}
}

func (room *Room) mainLoop() {
	// TODO: execute at x ticks?
	for {
		room.handleEvent()
		for _, obj := range room.objects {
			obj.update()
		}
	}
}

func genNewRoomCode() string {
	chars := "abcdefghijklmnopqrstuvwxyz123456789"
	source := rand.NewSource(time.Now().Unix())
	rg := rand.New(source)

	var roomCode string

	for i := 0; i < 8; i++ {
		roomCode += string(chars[rg.Intn(9)])
	}

	return roomCode
}
