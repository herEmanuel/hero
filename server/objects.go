package main

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
	initialBulletSpeed = 20
	friction           = 2
	bulletWidth        = 20
	bulletHeight       = 20
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
}

func newBullet(shooter *Player) *Bullet {
	// this should make the bullets come right from the center of the player
	return &Bullet{x: shooter.PosX + playerWidth/2 - bulletWidth/2,
		y:         shooter.PosY + playerHeight/2 - bulletHeight/2,
		w:         bulletWidth,
		h:         bulletHeight,
		speed:     initialBulletSpeed,
		direction: shooter.Direction,
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

func (b *Bullet) update() {
	b.x += int(b.direction.scale(float64(b.speed)).X)
	b.y += int(b.direction.scale(float64(b.speed)).Y)
	b.speed -= friction
}
