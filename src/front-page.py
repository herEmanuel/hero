import pygame
import pygame_gui
pygame.init()
# Cores
branco = (255, 255, 255)
cinza = (159, 182, 205)
preto = (0, 0, 0)

altura = 720
largura = 1080

fundo = pygame.display.set_mode([largura, altura])
background = pygame.transform.scale(pygame.image.load("../sprites/grama background.jpg"), [1080, 720])
clock = pygame.time.Clock()

manager = pygame_gui.UIManager((1080, 720))

start_button = pygame_gui.elements.UIButton(relative_rect=pygame.Rect((415, 500), (250, 50)), text="Jogar!",
                                    manager=manager)
caixa_nome = pygame_gui.elements.UITextEntryLine(relative_rect=pygame.Rect(415, 100, 250, 50),
                                    manager=manager)


class Heros(pygame.sprite.Sprite):

    def __init__(self, png1, png2):
        pygame.sprite.Sprite.__init__(self)
        self.sprites = []
        self.sprites.append(pygame.image.load(f'../sprites/{png1}'))
        self.sprites.append(pygame.image.load(f'../sprites/{png2}'))
        self.atual = 0
        self.image = self.sprites[self.atual]
        self.image = pygame.transform.scale(self.image, (320,320))

        self.rect = self.image.get_rect()
        self.rect.topleft = 380, 200

    def update(self):
        self.atual += 0.05
        if self.atual >= len(self.sprites):
            self.atual = 0
        self.image = self.sprites[int(self.atual)]
        self.image = pygame.transform.scale(self.image, (320, 320))

todas_sprites = pygame.sprite.Group()
hero = Heros('modelo 1 - andando.png', 'modelo 1 - andando 2.png')
todas_sprites.add(hero)

is_running = True
while is_running:
    relogio = clock.tick(60)/1000

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            is_running = False
        manager.process_events(event)
        if event.type == pygame_gui.UI_BUTTON_PRESSED:
            if event.ui_element == start_button:
                print('Começa o jogo')
                user = caixa_nome.get_text()
                print(user)
    fundo.blit(background, (0, 0))

    manager.update(relogio)
    manager.draw_ui(fundo)
    todas_sprites.draw(fundo)
    todas_sprites.update()
    pygame.display.update()
pygame.quit()
