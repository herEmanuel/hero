package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	playerWidth  = 50
	playerHeight = 50

	velocity = 5
)

type Player struct {
	conn   *websocket.Conn
	id     int
	Name   string `json:"name"`
	PosX   int    `json:"x"`
	PosY   int    `json:"y"`
	Weapon int

	room *Room
}

func newPlayer(conn *websocket.Conn, name string, id int, room *Room) *Player {
	// src := rand.NewSource(time.Now().Unix())
	// rg := rand.New(src)
	//TODO: adjust this later
	return &Player{
		conn:   conn,
		id:     id,
		Name:   name,
		PosX:   100,
		PosY:   100,
		Weapon: 0,
		room:   room,
	}
}

func (player *Player) processMovement(xOffset, yOffset int) {
	log.Println("at process movement")
	if xOffset > 0 {
		xOffset = 1
	} else if xOffset < 0 {
		xOffset = -1
	}

	if yOffset > 0 {
		yOffset = 1
	} else if yOffset < 0 {
		yOffset = -1
	}

	if player.PosX+xOffset*velocity > 0 && player.PosX+playerWidth+xOffset*velocity < mapWidth {
		if player.PosY+yOffset*velocity > 0 && player.PosY+playerHeight+yOffset*velocity < mapHeight {
			player.PosX += xOffset * velocity
			player.PosY += yOffset * velocity
		}
	}

	// TODO: check for collisions
}

func (player *Player) update() {
	for {
		//TODO: send all the player info in 1 message?
		for _, pl := range player.room.players {
			if pl == player {
				continue
			}

			err := player.conn.WriteJSON(&pl)
			if err != nil {
				log.Fatalln("fuck bruh error ", err)
			}

			time.Sleep(time.Millisecond * 50)
		}
	}
}

func (player *Player) receiveInput() {
	for {
		var event Event
		err := player.conn.ReadJSON(&event)
		if err != nil {
			log.Fatalln("fuck bruh error ", err)
		}
		event.playerId = player.id

		player.room.events <- event
	}
}

func joinRoom(conn *websocket.Conn, name, roomCode string) {
	log.Println("Joining room haha ", roomCode)
	var player *Player

	for _, room := range rooms {
		log.Println("a")
		if room.code != roomCode {
			continue
		}

		if len(room.players) >= maxPlayersPerRoom {
			// repeat all the handshake somehow idk
			return
		}

		player = newPlayer(conn, name, room.players[len(room.players)-1].id+1, room)
		room.players = append(room.players, player)
	}

	if player != nil {
		/* one goroutine for receiving user's input, and another for
		   sending the state of the game... ooof
		*/
		go player.receiveInput()
		player.update()
	}

	log.Println("got here, oh shit")
}

func createRoom(conn *websocket.Conn, playerName string) {
	player := newPlayer(conn, playerName, 0, nil)
	room := newRoom(player)
	player.room = room
	rooms = append(rooms, room)
	log.Println("new room code: ", room.code)
	go room.mainLoop()

	go player.receiveInput()
	player.update()
}
