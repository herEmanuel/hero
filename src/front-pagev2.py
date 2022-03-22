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

    manager.update(relogio)
    manager.draw_ui(tela)
    pygame.display.update()
pygame.quit()
