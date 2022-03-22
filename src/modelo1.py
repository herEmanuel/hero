import pygame
pygame.init()
tela = pygame.display.set_mode((1080,720))
relogio = pygame.time.Clock()

class Heros(pygame.sprite.Sprite):
    def __init__(self, pos_x, pos_y, png1, png2):
        pygame.sprite.Sprite.__init__(self)
        self.sprites = []
        self.sprites.append(pygame.image.load(f'../sprites/{png1}'))
        self.sprites.append(pygame.image.load(f'../sprites/{png2}'))
        self.atual = 0
        self.image = self.sprites[self.atual]
        self.image = pygame.transform.scale(self.image, (320,320))

        self.rect = self.image.get_rect()
        self.rect.topleft = pos_x, pos_y

    def update(self):
        self.atual += 0.12
        if self.atual >= len(self.sprites):
            self.atual = 0
        self.image = self.sprites[int(self.atual)]
        self.image = pygame.transform.scale(self.image, (320, 320))

todas_sprites = pygame.sprite.Group()
hero = Heros(100, 100, 'modelo 1 - andando.png', 'modelo 1 - andando 2.png')
todas_sprites.add(hero)
hero2 = Heros(420, 100, 'modelo 2 - andando.png', 'modelo 2 - andando 2.png')
todas_sprites.add(hero2)
hero3 = Heros(740, 100, 'modelo 3 - andando.png', 'modelo 3 - andando 2.png')
todas_sprites.add(hero3)
hero4 = Heros(100, 420, 'modelo 4 - andando.png', 'modelo 4 - andando 2.png')
todas_sprites.add(hero4)

is_running = True
while is_running:
    relogio.tick(30)
    tela.fill((255,255,255))
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            is_running = False
    todas_sprites.draw(tela)
    todas_sprites.update()
    pygame.display.flip()
pygame.quit()
