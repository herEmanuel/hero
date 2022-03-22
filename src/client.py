import websocket, json, pygame, threading, math
from pygame import Vector2
from enum import Enum

SERVER = "ws://localhost:8080/ws"
VELOCITY = 5

class Message(Enum):
    JoinRoom = 0
    CreateRoom = 1
    Data = 2
    Error = 3

class Event(Enum):
    Move = 0
    Shoot = 1

class Player(pygame.sprite.Sprite):
    def __init__(self, ws: websocket.WebSocket) -> None:
        super().__init__()

        self.ws = ws
        self.direction = Vector2(1, 0)
        self.pos = Vector2(0, 0)
        
        self.ogImage = pygame.image.load("../assets/start-button.png")
        self.image = self.ogImage

        self.rect = self.image.get_rect()

    def update(self) -> None:
        pressedKeys = pygame.key.get_pressed()
        if pressedKeys[pygame.K_w]:
            self.pos.y -= VELOCITY
            self.triggerMoveEvent(pygame.K_w)
        if pressedKeys[pygame.K_s]:
            self.pos.y += VELOCITY
            self.triggerMoveEvent(pygame.K_s)
        if pressedKeys[pygame.K_a]:
            self.pos.x -= VELOCITY
            self.triggerMoveEvent(pygame.K_a)
        if pressedKeys[pygame.K_d]:
            self.pos.x += VELOCITY
            self.triggerMoveEvent(pygame.K_d)

        mousePos = pygame.mouse.get_pos()   
        angle = math.atan2(mousePos[0] - self.pos.y, mousePos[1] - self.pos.x)
        print(math.degrees(angle))
        self.image = pygame.transform.rotate(self.ogImage, math.degrees(angle))
        self.rect = self.image.get_rect()
        self.rect.center = (self.pos.x, self.pos.y)

    def triggerMoveEvent(self, directionKey: int) -> None:
        pass

class Game():
    def __init__(self) -> None:
        pygame.init()
        self.screen = pygame.display.set_mode((800, 800))
        self.playerData = []

        self.serverThread = threading.Thread(target=self.websocketThread)
        self.serverThread.start()

        self.player = Player(self.ws)

        self.mainLoop()

    def mainLoop(self) -> None:
        while True:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    pygame.quit()
                    return

            self.player.update()

            self.screen.fill((255, 255, 255))
            self.screen.blit(self.player.image, self.player.rect)

            # if players:
            #     print(players)
            # for player in players:
            #     pygame.draw.rect(self.screen, (0, 255, 0), (player["x"], player["y"], 50, 50))

            pygame.display.flip()

    def websocketThread(self) -> None:
        self.ws = websocket.WebSocketApp(SERVER, on_open=self.performHandshake, on_message=self.onServerMessage)
        self.ws.run_forever()

    def performHandshake(self, conn: websocket.WebSocket) -> None:
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

    def onServerMessage(self) -> None:
        pass

# def on_message(conn: websocket.WebSocket, message):
#     #TODO: this only works when theres just one more player lol
#     global players
#     players = json.loads(message)

# def move_event(x: int, y: int):
#     event = json.dumps({"type": Message.Data.value, "event_type": Event.Move.value, "x_offset": x, "y_offset": y})
#     ws.send(event)

game = Game()
game.mainLoop()