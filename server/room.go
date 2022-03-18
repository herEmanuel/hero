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

type Room struct {
	code    string
	lock    sync.Mutex
	players []*Player

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
		if player.id == id {
			return player
		}
	}

	return nil
}

func (room *Room) mainLoop() {
	for {
		event := <-room.events

		player := room.getPlayer(event.playerId)
		if player == nil {
			log.Printf("Could not find get player %d. This should never happen.\n", event.playerId)
			continue
			// disconnect the player
		}

		switch event.EventType {
		case move:
			player.processMovement(event.XOffset, event.YOffset)
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

			//TODO: notify other players
		default:
			//error
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
