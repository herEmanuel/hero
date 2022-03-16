import websocket, json, pygame, threading
from enum import Enum

class Message(Enum):
    JoinRoom = 0
    CreateRoom = 1
    Data = 2
    Error = 3

class Event(Enum):
    Move = 0
    Shoot = 1

players = []

def handshake(conn: websocket.WebSocket):
    option = input("join or create?")
    if option == "join":
        code = input("Enter room code: ")
        joinRoom = {"type": Message.JoinRoom.value, "name": "test", "room": code}
        res = json.dumps(joinRoom)
        print(res)
        conn.send(res)
    else:
        createRoom = {"type": Message.CreateRoom.value, "name": "test"}
        res = json.dumps(createRoom)
        print(res)
        conn.send(res)

def on_message(conn: websocket.WebSocket, message):
    #TODO: this only works when theres just one more player lol
    global players
    players = [json.loads(message)]

ws = websocket.WebSocketApp("ws://localhost:8080/ws", on_open=handshake, on_message=on_message)

def move_event(x: int, y: int):
    event = json.dumps({"type": Message.Data.value, "event_type": Event.Move.value, "x_offset": x, "y_offset": y})
    ws.send(event)

def mainLoop():
    pygame.init()
    screen = pygame.display.set_mode((800, 800))

    rect = pygame.rect.Rect(100, 100, 50, 50)

    while True:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                return
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_w:
                    move_event(0, -1)
                    rect.y -= 5
                if event.key == pygame.K_a:
                    move_event(-1, 0)
                    rect.x -= 5
                if event.key == pygame.K_s:
                    move_event(0, 1)
                    rect.y += 5
                if event.key == pygame.K_d:
                    move_event(1, 0)
                    rect.x += 5

        screen.fill((255, 255, 255))
        pygame.draw.rect(screen, (255, 0, 0), rect)
        if players:
            print(players)
        for player in players:
            pygame.draw.rect(screen, (0, 255, 0), (player["x"], player["y"], 50, 50))
        pygame.display.flip()

thread = threading.Thread(target=mainLoop)
thread.start()
print("wait for it")
ws.run_forever()

thread.join()
