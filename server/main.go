package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// I haven't really tested any of this yet

const (
	serverPort    = 8080
	runningEnvVar = "RUNNING_ENV"

	mapWidth  = 800
	mapHeight = 800

	maxPlayersPerRoom = 8
)

const (
	joinRoomMsg = iota
	createRoomMsg
	dataMsg
	errorMsg
)

var rooms []*Room

var upgrader = websocket.Upgrader{}

func sendError(conn *websocket.Conn, message string) {
	var err map[string]interface{}
	err["type"] = errorMsg
	err["message"] = message
	conn.WriteJSON(&err) //TODO: handle err
	conn.Close()         // not sure
}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Upgrading connection with ", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection with ", r.RemoteAddr)
		return
	}

	handshake(conn)
}

func handshake(conn *websocket.Conn) {
	fmt.Println("Hey im here")
	var message map[string]interface{}
	err := conn.ReadJSON(&message) //TODO: handle error
	if err != nil {
		log.Fatalln("fuck bruh error ", err)
	}

	messageType, ok := message["type"].(float64)
	if !ok {
		log.Println("no type in the message, fck")
		log.Printf("%v\n", message)
		sendError(conn, "Expected the message to have a type")
		return
	}

	playerName, ok := message["name"].(string)
	if !ok {
		sendError(conn, "Expected the message to have the name of the player")
		return
	}

	switch int(messageType) {
	case joinRoomMsg:
		roomCode, ok := message["room"].(string)
		if !ok {
			sendError(conn, "Expected the message to have the room code")
			return
		}

		joinRoom(conn, playerName, roomCode)
	case createRoomMsg:
		createRoom(conn, playerName)
	default:
		// bad msg
	}

	// disconnect
}

func main() {
	if os.Getenv(runningEnvVar) == "production" {
		logFile, err := os.OpenFile("server-logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Could not create the file for logs")
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	http.HandleFunc("/ws", entryPoint)

	log.Println("Starting to listen for new connections")
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}
