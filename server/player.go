package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

const (
	playerWidth     = 50
	playerHeight    = 50
	maxHealth       = 100
	respawnCooldown = time.Second * 3

	velocity = 5
)

type Spawnpoint struct {
	x, y int
}

var spawnPoints []Spawnpoint = []Spawnpoint{
	{232, 225},
	{1196, 85},
	{2001, 244},
	{2286, 742},
	{1836, 1318},
	{857, 1264},
	{1283, 982},
	{71, 994},
}

type Player struct {
	conn      *websocket.Conn
	Id        int    `json:"player_id"`
	Name      string `json:"name"`
	PosX      int    `json:"x"`
	PosY      int    `json:"y"`
	Health    int    `json:"health"`
	Kills     int    `json:"kills"`
	Direction Vec2f  `json:"direction"`
	Weapon    int

	room      *Room
	events    chan Event
	isLeaving bool
	isDead    bool
}

func getRandomSpawnPos() (int, int) {
	src := rand.NewSource(time.Now().Unix())
	rg := rand.New(src)
	spawn := spawnPoints[rg.Intn(8)]
	return spawn.x, spawn.y
}

func newPlayer(conn *websocket.Conn, name string, id int, room *Room) *Player {
	x, y := getRandomSpawnPos()
	return &Player{
		conn:   conn,
		Id:     id,
		Name:   name,
		PosX:   x,
		PosY:   y,
		Health: maxHealth,
		room:   room,
		events: make(chan Event, maxEvents),
	}
}

func (player *Player) sendInitialInfo() {
	message := make(map[string]interface{})
	message["type"] = playerInfoMsg
	message["player"] = player
	if err := player.conn.WriteJSON(message); err != nil {
		log.Printf("Got an error while writing to %v, %v\n", player.conn.RemoteAddr().String(), err)
		//TODO: aaaa
		player.disconnect()
	}
}

func (player *Player) checkCollision(xOffset, yOffset, x, y, width, height int) bool {
	updatedPlayerX := xOffset * velocity
	updatedPlayerY := yOffset * velocity

	if updatedPlayerX < x+width && updatedPlayerX+playerWidth > x {
		if updatedPlayerY < y+height && updatedPlayerY+playerHeight > y {
			return true
		}
	}

	return false
}

func (player *Player) isWithinTheMap(xOffset, yOffset int) bool {
	if player.PosX+xOffset*velocity > 0 && player.PosX+playerWidth+xOffset*velocity < mapWidth {
		if player.PosY+yOffset*velocity > 0 && player.PosY+playerHeight+yOffset*velocity < mapHeight {
			return true
		}
	}

	return false
}

func (player *Player) processMovement(xOffset, yOffset int, direction Vec2f) {
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

	player.Direction = direction.normalize()

	if player.isWithinTheMap(xOffset, yOffset) {
		for _, obstacle := range mapObstacles {
			if player.checkCollision(xOffset, yOffset, obstacle.posX(), obstacle.posY(), obstacle.width(), obstacle.height()) {
				return
			}
		}

		// TODO: allow players to slowly push others
		for _, pl := range player.room.players {
			if player.checkCollision(xOffset, yOffset, pl.PosX, pl.PosY, playerWidth, playerHeight) {
				return
			}
		}

		player.PosX += xOffset * velocity
		player.PosY += yOffset * velocity
	}
}

func (player *Player) respawn() {
	player.isDead = true

	time.Sleep(respawnCooldown)

	if player.isLeaving {
		return
	}

	x, y := getRandomSpawnPos()
	player.Health = maxHealth
	player.Direction = makeVec2f(0, -1)
	player.PosX = x
	player.PosY = y

	respawnEvent := RespawnOrJoinEvent{
		Type:      dataMsg,
		EventType: respawnEvnt,
		Player:    player,
	}
	player.room.notifyAllPlayers(respawnEvent)
}

func (player *Player) update() {
	for {
		if player.isLeaving {
			return
		}

		if len(player.events) > 0 {
			event := <-player.events
			if err := player.conn.WriteJSON(event); err != nil {
				log.Printf("Got an error while writing to %v, %v\n", player.conn.RemoteAddr().String(), err)
				if !player.isLeaving {
					// TODO: aaaaa
					player.disconnect()
					return
				}
			}
		}

		posUpdate := make(map[string]interface{})
		posUpdate["type"] = dataMsg
		posUpdate["players"] = player.room.players
		if err := player.conn.WriteJSON(posUpdate); err != nil {
			log.Printf("Got an error while writing to %v, %v\n", player.conn.RemoteAddr().String(), err)
			if !player.isLeaving {
				// TODO: aaaaa
				player.disconnect()
				return
			}
		}

		time.Sleep(time.Millisecond * 50) //TODO: ??
	}
}

func (player *Player) receiveInput() {
	for {
		var event GenericEvent

		if err := player.conn.ReadJSON(&event); err != nil {
			log.Printf("Got an error while reading from %v, %v\n", player.conn.RemoteAddr().String(), err)
			player.disconnect()
			return
		}

		if player.isDead && event.EventType != leaveRoomEvnt {
			continue
		}

		event.PlayerId = player.Id

		if event.EventType == leaveRoomEvnt {
			if len(player.room.players) == 1 {
				player.room.isClosing = true
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
	player.isLeaving = true

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

		player = newPlayer(conn, name, room.players[len(room.players)-1].Id+1, room)

		room.lock.Lock()
		room.players = append(room.players, player)
		room.lock.Unlock()

		joinRoomEvent := RespawnOrJoinEvent{
			Type:      dataMsg,
			EventType: joinRoomEvnt,
			Player:    player,
		}
		room.notifyAllOtherPlayers(joinRoomEvent)

		player.sendInitialInfo()
	}

	if player != nil {
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

	player.sendInitialInfo()

	go player.update()
	player.receiveInput()
}
