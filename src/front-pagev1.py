import pygame
import pygame_gui as gui
pygame.init()
# Cores
branco = (255, 255, 255)
cinza = (159, 182, 205)
preto = (0, 0, 0)

altura = 720
largura = 1080
tamanho_botoes = 250

fundo = pygame.display.set_mode([largura, altura])
background = pygame.transform.scale(pygame.image.load("grama background.jpg"),[1080, 720])
botao = pygame.transform.scale(pygame.image.load('start-button.png'), (250, 103))
clock = pygame.time.Clock()

manager = gui.UIManager((1080, 720))
fonte_padrao = pygame.font.get_default_font()

start_button = gui.elements.UIButton(relative_rect=pygame.Rect((415, 500), (250, 50)), text="Jogar!",
                                    manager=manager)
user = gui.elements.UITextEntryLine(relative_rect=pygame.Rect(415, 100, 250, 50),
                                    manager=manager).get_text()

loop = True
while loop:
    relogio = clock.tick(60)/1000

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            loop = False
        manager.process_events(event)
        if event.type == gui.UI_BUTTON_PRESSED:
            if event.ui_element == start_button:
                print('Começa o jogo')
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