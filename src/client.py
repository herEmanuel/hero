from typing import Tuple
import websocket, json, pygame, threading, math
from pygame import Vector2
from enum import Enum

SERVER = "ws://localhost:8080/ws"
RESOLUTION = (800, 640)
VELOCITY = 5

class Message(Enum):
    JoinRoom = 0
    CreateRoom = 1
    PlayerInfo = 2
    Data = 2
    Error = 3

class Event(Enum):
    Move = 0
    Shoot = 1
    Death = 2
    Respawn = 3
    JoinRoom = 4
    LeaveRoom = 5

class Player(pygame.sprite.Sprite):
    def __init__(self) -> None:
        super().__init__()
        self.pos = Vector2(0, 0)
        self.direction = Vector2(0, 0)

        self.width = 120
        self.height = 120

    def center(self) -> Tuple[int, int]:
        return (self.pos.x + self.width/2, self.pos.y + self.height/2)

class LocalPlayer(Player):
    def __init__(self, ws: websocket.WebSocket) -> None:
        super().__init__()
        self.ws = ws

        self.ogImage = pygame.image.load("../assets/player.png")
        self.image = self.ogImage
        self.rect = self.image.get_rect()   
        self.rect.center = (RESOLUTION[0]/2, RESOLUTION[1]/2)

        self.collisionRect = self.image.get_rect()
        self.collisionRect.center = (RESOLUTION[0]/2, RESOLUTION[1]/2)

        self.lastMousePos = (0, 0)

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
        # print(mousePos)
        if pygame.mouse.get_focused() and mousePos != self.lastMousePos:
            directionVector = Vector2(mousePos[0] - RESOLUTION[0]/2, mousePos[1] - RESOLUTION[1]/2)
            angle = math.atan2(directionVector.y, directionVector.x)*-1
            # print(math.degrees(angle))
            self.image = pygame.transform.rotate(self.ogImage, math.degrees(angle))
            self.rect = self.image.get_rect()   
            self.rect.center = (RESOLUTION[0]/2, RESOLUTION[1]/2)
            self.direction = directionVector.normalize()
            self.lastMousePos = mousePos

            self.triggerMoveEvent(None)

    def triggerMoveEvent(self, directionKey) -> None:
        #TODO: simplify this later
        x, y = 0, 0
        if directionKey is not None:
            if directionKey == pygame.K_w:
                y = -1
            elif directionKey == pygame.K_s:
                y = 1
            elif directionKey == pygame.K_a:
                x = -1
            elif directionKey == pygame.K_d:
                x = 1

        direction = {"dir_x": self.direction.x, "dir_y": self.direction.y}
        event = json.dumps({"type": Message.Data.value, "event_type": Event.Move.value, "x_offset": x, "y_offset": y, "direction": direction})
        self.ws.send(event)

    def shoot(self) -> None:
        shootEvent = json.dumps({"type": Message.Data.value, "event_type": Event.Shoot.value})
        print(shootEvent)
        self.ws.send(shootEvent)

class Bullet(pygame.sprite.Sprite):
    def __init__(self, shooter: Player) -> None:
        super().__init__()

        self.direction = shooter.direction
        self.pos = Vector2(shooter.center()[0], shooter.center()[1])
        self.velocity = 30
        self.friction = 1

        angle = math.atan2(self.direction.y, self.direction.x)*-1
        self.image = pygame.transform.rotate(pygame.image.load("../assets/bullet.png"), math.degrees(angle))
        self.rect = self.image.get_rect()
        self.rect.center = (self.pos.x, self.pos.y)

        self.collisionRect = self.image.get_rect()
        self.collisionRect.center = self.rect.center

    def update(self) -> None:
        self.rect.center = (self.rect.center[0] + self.direction.x * self.velocity, self.rect.center[1] + self.direction.y * self.velocity)
        self.velocity -= self.friction

        if self.velocity <= 0:
            self.kill()
 
class RemotePlayer(Player):
    def __init__(self, playerData: dict) -> None:
        super().__init__()

        self.id = playerData["player_id"]
        self.name = playerData["name"]
        self.pos = Vector2(playerData["x"], playerData["y"])
        self.direction = Vector2(playerData["direction"]["dir_x"], playerData["direction"]["dir_y"])
        self.health = playerData["health"]
        self.kills = playerData["kills"]

        self.ogImage = pygame.image.load("../assets/player.png")
        self.image = self.ogImage

        self.rect = self.image.get_rect()
        self.rect.center = (self.pos.x, self.pos.y)

    def update(self, updatedData: dict) -> None:
        self.pos = Vector2(updatedData["x"], updatedData["y"])
        self.direction = Vector2(updatedData["direction"]["dir_x"], updatedData["direction"]["dir_y"])
        self.health = updatedData["health"]
        self.kills = updatedData["kills"]

        angle = math.atan2(self.direction.y, self.direction.x)*-1
        self.image = pygame.transform.rotate(self.ogImage, math.degrees(angle))
        self.rect = self.image.get_rect()
        self.rect.topleft = (self.pos.x, self.pos.y)

class Camera():
    def __init__(self, x: int, y: int) -> None:
        self.x = x
        self.y = y

    def update(self, playerX: int, playerY: int) -> None:
        self.x = playerX - RESOLUTION[0]/2 + 120/2
        self.y = playerY - RESOLUTION[1]/2 + 120/2

class Game():
    def __init__(self) -> None:
        pygame.init()
        self.screen = pygame.display.set_mode(RESOLUTION)

        self.serverThread = threading.Thread(target=self.websocketThread)
        self.serverThread.start()

        self.player = LocalPlayer(self.ws)
        self.firstPlayerUpdate = True

        self.mainLoop()

    def mainLoop(self) -> None:
        self.remotePlayers = pygame.sprite.Group()
        bullets = pygame.sprite.Group() 
        mapImg = pygame.image.load("../assets/map.png")
        camera = Camera(self.player.pos.x, self.player.pos.y)

        while True:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    pygame.quit()
                    self.ws.close()
                    exit(0)

                if event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_SPACE:
                        self.player.shoot()
                        bullets.add(Bullet(self.player))

            self.player.update()
            bullets.update()
            camera.update(self.player.pos.x, self.player.pos.y)

            self.screen.blit(mapImg, (-camera.x, -camera.y))

            for b in bullets:
                pygame.draw.rect(self.screen, (255, 255, 255), (b.rect.x - camera.x, b.rect.y - camera.y, b.rect.width, b.rect.height))
                # print(b.rect.size)
                self.screen.blit(b.image, (b.rect.x - camera.x, b.rect.y - camera.y))

            pygame.draw.rect(self.screen, (0, 0, 0), self.player.collisionRect)
            self.screen.blit(self.player.image, self.player.rect)
            
            for p in self.remotePlayers:
                self.screen.blit(p.image, (p.rect.x - camera.x, p.rect.y - camera.y))

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

    def onServerMessage(self, conn: websocket.WebSocket, message) -> None:
        serverMsg = json.loads(message)
        print(serverMsg)
        if type(serverMsg) is dict:
            if "players" in serverMsg:
                if self.firstPlayerUpdate:
                    for remotePlayer in serverMsg["players"]:
                        if remotePlayer["player_id"] == self.player.id:
                            continue

                        self.remotePlayers.add(RemotePlayer(remotePlayer))

                    self.firstPlayerUpdate = False
                else:
                    for player in self.remotePlayers:
                        for updatedPlayer in serverMsg["players"]:
                            if updatedPlayer["player_id"] == player.id:
                                player.update(updatedPlayer)
                
            elif serverMsg["type"] == Message.PlayerInfo.value:
                self.player.id = serverMsg["player"]["player_id"]
                self.player.pos = Vector2(serverMsg["player"]["x"], serverMsg["player"]["y"])
                self.player.health = serverMsg["player"]["health"]

            elif serverMsg["event_type"] == Event.JoinRoom.value:
                print(serverMsg["player"])
                self.remotePlayers.add(RemotePlayer(serverMsg["player"]))
                print(len(self.remotePlayers))

            elif serverMsg["event_type"] == Event.LeaveRoom.value:
                for player in self.remotePlayers:
                    if player.id == serverMsg["player_id"]:
                        self.remotePlayers.remove(player)

game = Game()
game.mainLoop()