package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

// I haven't really tested any of this yet

const (
	SERVER_PORT     = 8080
	RUNNING_ENV_VAR = "RUNNING_ENV"
)

const (
	JOIN_ROOM = iota
	CREATE_ROOM
)

type Room struct {
	code    string
	players []net.Conn
}

var rooms []Room

func main() {
	if os.Getenv(RUNNING_ENV_VAR) == "production" {
		logFile, err := os.OpenFile("server-logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Could not create the file for logs")
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	listener, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", SERVER_PORT))
	if err != nil {
		log.Fatalln("Could not create the tcp server, ", err)
	}

	defer listener.Close()

	log.Println("Starting to listen for new connections")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("An error occurred while accepting a connection, ", err)
		}

		go newConnection(conn)
	}
}

func newConnection(conn net.Conn) {
	log.Println("Established a connection with ", conn.RemoteAddr().String())

	var buffer []byte
	_, err := conn.Read(buffer)
	if err != nil {
		// send repeat msg
	}

	var message map[string]interface{}
	err = json.Unmarshal(buffer, &message)
	if err != nil {
		// send repeat msg
	}

	if _, hasKey := message["type"]; !hasKey {
		// send repeat msg
	}
	messageType, ok := message["type"].(int)
	if !ok {
		// send repeat msg
	}

	switch messageType {
	case JOIN_ROOM:
		if _, hasKey := message["roomCode"]; !hasKey {
			// send repeat msg
		}

		roomCode, ok := message["roomCode"].(string)
		if !ok {
			// send repeat msg
		}

		joinRoom(conn, roomCode)
	case CREATE_ROOM:
		createRoom(conn)
	default:
		fmt.Print("oof")
		// send repeat msg
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

func createRoom(creator net.Conn) {
	newRoom := Room{genNewRoomCode(), []net.Conn{creator}}
	rooms = append(rooms, newRoom)
	fmt.Println("created yee")
}

func joinRoom(player net.Conn, roomCode string) {
	for _, room := range rooms {
		if room.code != roomCode {
			continue
		}

		room.players = append(room.players, player)
	}
	fmt.Println("joined yee")
}
