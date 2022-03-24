package main

import "log"

// Random stuff on the map, e.g. bullets, obstacles, or power-ups
type Object interface {
	objType() int
	posX() int
	posY() int
	width() int
	height() int
	update()
}

// Types of objects
const (
	obstacleObj = iota
	bulletObj
)

const (
	initialBulletSpeed = 30
	friction           = 1
	bulletWidth        = 62
	bulletHeight       = 19
	bulletDamage       = 20
)

type Obstacle struct {
	x, y, w, h int
}

func (o Obstacle) objType() int {
	return obstacleObj
}

func (o Obstacle) posX() int {
	return o.x
}

func (o Obstacle) posY() int {
	return o.y
}

func (o Obstacle) width() int {
	return o.w
}

func (o Obstacle) height() int {
	return o.h
}

func (o Obstacle) update() {
}

type Bullet struct {
	x, y, w, h int
	speed      int
	direction  Vec2f
	shooter    *Player
	room       *Room
}

func newBullet(shooter *Player) *Bullet {
	// this should make the bullets come right from the center of the player
	return &Bullet{x: shooter.PosX + playerWidth/2 - bulletWidth/2,
		y:         shooter.PosY + playerHeight/2 - bulletHeight/2,
		w:         bulletWidth,
		h:         bulletHeight,
		speed:     initialBulletSpeed,
		direction: shooter.Direction,
		shooter:   shooter,
		room:      shooter.room,
	}
}

func (b Bullet) objType() int {
	return bulletObj
}

func (b Bullet) posX() int {
	return b.x
}

func (b Bullet) posY() int {
	return b.y
}

func (b Bullet) width() int {
	return b.w
}

func (b Bullet) height() int {
	return b.h
}

func (b *Bullet) vanish() {
	for i, obj := range b.room.objects {
		if bullet, ok := obj.(*Bullet); !ok || bullet != b {
			continue
		}

		b.room.lock.Lock()
		b.room.objects = append(b.room.objects[:i], b.room.objects[i+1:]...)
		b.room.lock.Unlock()
	}
}

func (b *Bullet) update() {
	b.x += int(b.direction.scale(float64(b.speed)).X)
	b.y += int(b.direction.scale(float64(b.speed)).Y)
	b.speed -= friction

	if b.speed <= 0 {
		b.vanish()
	}

	for _, player := range b.room.players {
		if b.x < player.PosX+playerWidth && b.x+bulletWidth > player.PosX {
			if b.y < player.PosY+playerHeight && b.y+bulletHeight > player.PosY {
				log.Printf("shot player %d\n", player.Id)
				if player.isDead || player == b.shooter {
					continue
				}

				b.vanish()

				player.Health -= bulletDamage
				if player.Health <= 0 {
					b.shooter.Kills += 1

					deathEvent := DeathEvent{
						Type:      dataMsg,
						EventType: deathEvnt,
						PlayerId:  player.Id,
						KillerId:  b.shooter.Id,
					}
					b.room.notifyAllPlayers(deathEvent)

					go player.respawn()
				}
			}
		}
	}
}
