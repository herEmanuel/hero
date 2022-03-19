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
	Health int    `json:"health"`
	Kills  int    `json:"kills"`
	Weapon int

	room    *Room
	leaving bool
}

func getSpawnPosition() (int, int) {
	return 0, 0
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
		if player.leaving {
			return
		}

		//TODO: send all the player info in 1 message?
		for _, pl := range player.room.players {
			if player.leaving {
				return
			}

			if pl == player {
				continue
			}

			if err := player.conn.WriteJSON(&pl); err != nil {
				log.Printf("Got an error while writing to %v, %v\n", player.conn.RemoteAddr().String(), err)
				player.disconnect()
			}

			time.Sleep(time.Millisecond * 50) //TODO: ??
		}
	}
}

func (player *Player) receiveInput() {
	for {
		var event Event

		if err := player.conn.ReadJSON(&event); err != nil {
			log.Printf("Got an error while reading from %v, %v\n", player.conn.RemoteAddr().String(), err)
			player.disconnect()
		}

		event.playerId = player.id

		if event.EventType == leaveRoom {
			if len(player.room.players) == 1 {
				player.room.closing = true
				player.room.events <- event
				return
			}

			player.room.events <- event

			player.disconnect()
			return
		}

		player.room.events <- event
	}
}

func (player *Player) disconnect() {
	player.leaving = true

	for i, pl := range player.room.players {
		if pl == player {
			player.room.lock.Lock()
			player.room.players = append(player.room.players[:i], player.room.players[i+1:]...)
			player.room.lock.Unlock()
		}
	}
}

func joinRoom(conn *websocket.Conn, name, roomCode string) {
	log.Println("Joining room haha ", roomCode)
	var player *Player

	for _, room := range globalState.rooms {
		if room.code != roomCode {
			continue
		}

		if len(room.players) >= maxPlayersPerRoom {
			sendError(conn, fullRoomErr)
			return
		}

		player = newPlayer(conn, name, room.players[len(room.players)-1].id+1, room)

		room.lock.Lock()
		room.players = append(room.players, player)
		room.lock.Unlock()
	}

	if player != nil {
		/* one goroutine for receiving user's input, and another for
		   sending the state of the game... ooof
		*/
		go player.update()
		player.receiveInput()
		//BIG TODO: does ReadJSON keep returning an error if the connection has been closed already?
		return
	}

	sendError(conn, invalidRoomCodeErr)
}

func createRoom(conn *websocket.Conn, playerName string) {
	player := newPlayer(conn, playerName, 0, nil)
	room := newRoom(player)
	player.room = room

	globalState.lock.Lock()
	globalState.rooms = append(globalState.rooms, room)
	globalState.lock.Unlock()

	log.Printf("New room %s created by %s: \n", room.code, conn.RemoteAddr().String())

	go room.mainLoop()

	go player.update()
	player.receiveInput()
}
