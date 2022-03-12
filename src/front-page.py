import pygame
import pygame_gui as gui
pygame.init()
# Cores
branco = (255, 255, 255)
cinza = (159, 182, 205)
preto = (0, 0, 0)

altura = 720
largura = 1080
tamanho = 250
meio = (largura-tamanho)/2
# print(meio)

fundo = pygame.display.set_mode([largura, altura])
background = pygame.transform.scale(pygame.image.load("grama background.jpg"),[1080, 720])
botao = pygame.transform.scale(pygame.image.load('start-button.png'), (250, 103))
clock = pygame.time.Clock()

manager = gui.UIManager((1080, 720), 'theme.json')
fonte_padrao = pygame.font.get_default_font()

test_button = gui.elements.UIButton(relative_rect=pygame.Rect((415, 300), (250, 50)), text="Say Hello",
                                    manager=manager)
gui.elements.UITextEntryLine(relative_rect=pygame.Rect(415, 200, 250, 50), manager=manager)
def caixa_nome():
    pygame.draw.rect(background, branco, (415, 100, 250, 50))
nada_escrito = True
nome_cinza = pygame.font.SysFont(fonte_padrao, 50)

# Digitando na caixa de nome:
teclas = [pygame.K_a, pygame.K_b, pygame.K_c, pygame.K_d, pygame.K_e, pygame.K_f, pygame.K_g,
          pygame.K_h, pygame.K_i, pygame.K_j, pygame.K_k, pygame.K_l, pygame.K_m, pygame.K_n, pygame.K_o,
          pygame.K_p, pygame.K_q, pygame.K_r, pygame.K_s, pygame.K_t, pygame.K_u, pygame.K_v, pygame.K_w,
          pygame.K_x, pygame.K_y, pygame.K_z, pygame.K_0, pygame.K_1, pygame.K_2, pygame.K_3, pygame.K_4,
          pygame.K_5, pygame.K_6, pygame.K_7, pygame.K_8, pygame.K_9]
letras = ["a",'b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z',
          '0','1','2','3','4','5','6','7','8','9']
pos_x = 420 # Posição das letras na caixa de nome

loop = True
while loop:
    relogio = clock.tick(60)
    manager.update(relogio)
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            loop = False
        manager.process_events(event)
        if event.type == gui.UI_BUTTON_PRESSED:
            if event.ui_element == test_button:
                print("Hello, World!")
        if event.type == pygame.MOUSEBUTTONDOWN:
            if pygame.mouse.get_pressed()[0]:
                mouse_x = pygame.mouse.get_pos()[0]
                mouse_y = pygame.mouse.get_pos()[1]
                if 415 < mouse_x < 665 and 500 < mouse_y < 603:
                    print('Começa o jogo')
                if 415 <mouse_x<665 and 100< mouse_y< 150:
                    nada_escrito = False
                    print('Digitando...')
                    # Não está funcionando, talvez muitas tarefas
                    for event in pygame.event.get():
                        if event.type == pygame.KEYDOWN:
                            if event.type in teclas:
                                print(event)
                                pos = teclas.index(event.type)
                                texto = pygame.font.SysFont(fonte_padrao, 50).render(f"{letras[pos]}", True, preto)
                                background.blit(texto, (pos_x, 110))
                                pos_x += 15
    fundo.blit(background, (0, 0))
    fundo.blit(botao, (415,500))
    caixa_nome()
    gui.elements.UITextEntryLine(relative_rect=pygame.Rect(415, 200, 250, 50), manager=manager)
    texto = nome_cinza.render("NOME", True, cinza)
    if nada_escrito:
        fundo.blit(texto, (420, 110))
    manager.draw_ui(fundo)
    manager.update(relogio)
    pygame.display.update()
pygame.quit()
