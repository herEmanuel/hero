import pygame
import pygame_gui
pygame.init()

# Cores
branco = (255, 255, 255)
cinza = (159, 182, 205)
preto = (0, 0, 0)

altura = 720
largura = 1080

tela = pygame.display.set_mode([largura, altura])
fundo = pygame.transform.scale(pygame.image.load('../assets/map.png'), (1280, 720))
clock = pygame.time.Clock()

manager = pygame_gui.UIManager((1080, 720))

start_button = pygame_gui.elements.UIButton(relative_rect=pygame.Rect((415, 500), (250, 50)), text="Jogar!",
                                    manager=manager)
caixa_nome = pygame_gui.elements.UITextEntryLine(relative_rect=pygame.Rect(415, 100, 250, 50),
                                    manager=manager)

class Sprites(pygame.sprite.Sprite):

    def __init__(self,png1, png2, png3, png4):
        pygame.sprite.Sprite.__init__(self)
        self.sprites = []
        self.sprites.append(pygame.image.load(f'../sprites/{png1}'))
        self.sprites.append(pygame.image.load(f'../sprites/{png2}'))
        self.sprites.append(pygame.image.load(f'../sprites/{png3}'))
        self.sprites.append(pygame.image.load(f'../sprites/{png4}'))
        self.atual = 0
        self.image = self.sprites[self.atual]
        self.image = pygame.transform.scale(self.image, (320, 320))

        self.rect = self.image.get_rect()
        self.rect.topleft = 380, 200

    def update(self):
        self.atual += 1
        if self.atual >= len(self.sprites):
            self.atual = 0
        self.image = self.sprites[int(self.atual)]
        self.image = pygame.transform.scale(self.image, (320, 320))

sprites = pygame.sprite.Group()
hero = Sprites('modelo 1 - parado.png', 'modelo 2 - parado.png','modelo 3 - parado.png','modelo 4 - parado.png')
sprites.add(hero)

is_running = True
while is_running:
    relogio = clock.tick(60)/1000
    tela.fill(branco)
    tela.blit(fundo, (0,0))
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            is_running = False
        manager.process_events(event)
        if event.type == pygame_gui.UI_BUTTON_PRESSED:
            if event.ui_element == start_button:
                print('Começa o jogo')
                nickname = caixa_nome.get_text()
                print(nickname)
                sprites.update()
    sprites.draw(tela)
    manager.update(relogio)
    manager.draw_ui(tela)
    pygame.display.update()
pygame.quit()
