package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	serverPort    = 8080
	runningEnvVar = "RUNNING_ENV"

	mapWidth  = 2560
	mapHeight = 1440

	maxPlayersPerRoom = 16
)

// Message types
const (
	joinRoomMsg = iota
	createRoomMsg
	dataMsg
	errorMsg
)

const (
	invalidRoomCodeErr = iota
	fullRoomErr
)

type State struct {
	rooms []*Room
	lock  sync.Mutex
}

var globalState State

var upgrader = websocket.Upgrader{}

func sendError(conn *websocket.Conn, errorCode int) {
	message := map[string]int{"type": errorMsg, "code": errorCode}
	if err := conn.WriteJSON(&message); err != nil {
		if _, ok := err.(*websocket.CloseError); !ok {
			// (probably, i think) an error related to the unmarshlling of the json
			conn.Close() // TODO: close the conn gracefully
		}
	}
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
	for {
		log.Println("at handshake")
		var message map[string]interface{}
		if err := conn.ReadJSON(&message); err != nil {
			log.Printf("Got an error while reading from %v, %v\n", conn.RemoteAddr().String(), err)

			if _, ok := err.(*websocket.CloseError); !ok {
				// (probably, i think) an error related to the unmarshlling of the json
				conn.Close() // TODO: close the conn gracefully
			}
		}

		messageType, ok := message["type"].(float64)
		if !ok {
			conn.Close() // TODO: close the conn gracefully
			return
		}

		playerName, ok := message["name"].(string)
		if !ok {
			conn.Close() // TODO: close the conn gracefully
			return
		}

		switch int(messageType) {
		case joinRoomMsg:
			roomCode, ok := message["room"].(string)
			if !ok {
				conn.Close() // TODO: close the conn gracefully
				return
			}

			joinRoom(conn, playerName, roomCode)

			/*
				If we returned from joinRoom, that either means that the room code
				is invalid, the player left the room, or the room is full,
				so we just repeat the handshake process
			*/

			continue
		case createRoomMsg:
			createRoom(conn, playerName)
			continue
		default:
			conn.Close() // TODO: close the conn gracefully
			return
		}
	}
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
