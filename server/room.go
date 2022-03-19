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

type Obstacle struct {
	x      int
	y      int
	width  int
	height int
}

/* We just have one non-customizable map, so we're just gonna
   hardcode all of the obstacles in it
*/

var mapObstacles []Obstacle = []Obstacle{
	Obstacle{530, 217, 65, 65},
	Obstacle{812, 47, 32, 104},
	Obstacle{301, 381, 90, 58},
	Obstacle{1047, 267, 60, 60},
	Obstacle{944, 507, 69, 84},
	Obstacle{1248, 557, 104, 32},
	Obstacle{1257, 679, 75, 75},
	Obstacle{20, 636, 81, 81},
	Obstacle{91, 692, 73, 73},
	Obstacle{640, 825, 81, 81},
	Obstacle{1997, 24, 65, 65},
	Obstacle{2423, 467, 60, 60},
	Obstacle{2423, 540, 60, 60},
	Obstacle{1248, 1217, 104, 32},
	Obstacle{1320, 1288, 32, 104},
}

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
