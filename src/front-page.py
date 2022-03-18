import pygame
import pygame_gui as gui
pygame.init()
# Cores
branco = (255, 255, 255)
cinza = (159, 182, 205)
preto = (0, 0, 0)

class HeroIoWindow(gui.elements.UIWindow):
    altura = 720
    largura = 1080
    tamanho_botoes = 250

    def __init__(self, manager):
        super.init(pygame.Rect((0,0),(1080,720)), manager, window_display_title="Hero.io", object_id="#hero_window")

        altura = 720
        largura = 1080

        self.fundo = pygame.display.set_mode([largura, altura])
        self.background = pygame.transform.scale(pygame.image.load("grama background.jpg"),[1080, 720])
        #self.botao = pygame.transform.scale(pygame.image.load('start-button.png'), (250, 103))
        self.clock = pygame.time.Clock()

        manager = gui.UIManager((1080, 720))
        fonte_padrao = pygame.font.get_default_font()

        start_button = gui.elements.UIButton(
            relative_rect=pygame.Rect((415, 500), (250, 50)),
            text="Jogar!", manager=manager)

        name_input = gui.elements.UITextEntryLine(
            relative_rect=pygame.Rect(415, 100, 250, 50),
            manager=manager).set_text_length_limit(16)

class Player:

    def __init__(self, manager, nickname):
        self.nickname = nickname.get_text()

name_input = gui.elements.UITextEntryLine(
    relative_rect=pygame.Rect(415, 100, 250, 50),
    manager=manager).set_text_length_limit(16)

loop = True
while loop:
    HeroIoWindow()
    relogio = clock.tick(60)/1000

    for event in pygame.event.get():
        manager.process_events(event)
        if event.type == pygame.QUIT:
            loop = False
        if event.type == gui.UI_BUTTON_PRESSED:
            if event.ui_element == start_button:
                print('Começa o jogo')
                Player1 = Player(manager, name_input)
                print(Player1.nickname)
        if event.type == pygame.MOUSEBUTTONDOWN:
            if pygame.mouse.get_pressed()[0]:
                mouse_x = pygame.mouse.get_pos()[0]
                mouse_y = pygame.mouse.get_pos()[1]
                #if 415 < mouse_x < 665 and 500 < mouse_y < 603:
    fundo.blit(background, (0, 0))
    manager.update(relogio)
    manager.draw_ui(fundo)
    pygame.display.update()
pygame.quit()
