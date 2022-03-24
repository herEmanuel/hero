import pygame
import pygame_gui

pygame.init()

# Cores
branco = (255, 255, 255)
cinza_claro = (159, 182, 205)
preto = (0, 0, 0)

# Definições Iniciais
altura = 720
largura = 1080
tela = pygame.display.set_mode([largura, altura])
fundo = pygame.transform.scale(pygame.image.load('../assets/map.png'), (1280, 720))
clock = pygame.time.Clock()

manager = pygame_gui.UIManager((1080, 720))

# Botões
create_button = pygame_gui.elements.UIButton(relative_rect=pygame.Rect((415, 610), (124, 50)), text="Criar",
                                    manager=manager)
join_button = pygame_gui.elements.UIButton(relative_rect=pygame.Rect((541, 610), (124, 50)), text="Entrar",
                                    manager=manager)
previous_button = pygame_gui.elements.UIButton(relative_rect=pygame.Rect((415, 420), (124, 50)), text="< Anterior",
                                    manager=manager)
next_button = pygame_gui.elements.UIButton(relative_rect=pygame.Rect((541, 420), (124, 50)), text="Próximo >",
                                    manager=manager)
caixa_nome = pygame_gui.elements.UITextEntryLine(relative_rect=pygame.Rect(415, 80, 250, 50),
                                    manager=manager)
caixa_code = pygame_gui.elements.UITextEntryLine(relative_rect=pygame.Rect(415, 550, 250, 50),
                                    manager=manager)

# Labels

def size(self, text):
    return pygame.font.Font.size(self, text)

pygame.font.init()
fonte_default = pygame.font.get_default_font()
nome_txt = pygame.font.SysFont(fonte_default, 30)
code_txt = pygame.font.SysFont(fonte_default, 30)
nome_size = size(nome_txt, "Nickname:")
code_size = size(code_txt, "Código da Sala:")
# print(nome_size, code_size)


class Sprites(pygame.sprite.Sprite):

    def __init__(self, pngs):
        pygame.sprite.Sprite.__init__(self)
        self.sprites = []
        for png in pngs:
            self.sprites.append(pygame.image.load(f'../sprites/{png}'))
        self.atual = 0
        self.image = self.sprites[self.atual]
        self.image = pygame.transform.scale(self.image, (320, 320))

        self.rect = self.image.get_rect()
        self.rect.topleft = 380, 150

    def update(self, anterior):
        if anterior:
            self.atual -= 1
            if self.atual < 0:
                self.atual = len(self.sprites)-1
            self.image = self.sprites[self.atual]
            self.image = pygame.transform.scale(self.image, (320, 320))
        else:
            self.atual += 1
            if self.atual >= len(self.sprites):
                self.atual = 0
            self.image = self.sprites[int(self.atual)]
            self.image = pygame.transform.scale(self.image, (320, 320))

pngs = ['modelo 1 - parado.png', 'modelo 2 - parado.png','modelo 3 - parado.png','modelo 4 - parado.png']
sprites = pygame.sprite.Group()
hero = Sprites(pngs)
sprites.add(hero)

is_running = True
while is_running:
    relogio = clock.tick(60)/1000
    tela.fill(branco)
    tela.blit(fundo, (0,0))
    caixa_nome.set_text_length_limit(20)

    # Labels
    texto_n = nome_txt.render("Nickname:", True, (preto))
    tela.blit(texto_n, (488, 60))
    texto_c = code_txt.render("Código da Sala:", True, (preto))
    tela.blit(texto_c, (462, 529))

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            is_running = False
        manager.process_events(event)

        # Funcionalidade dos botões
        if event.type == pygame_gui.UI_BUTTON_PRESSED:
            if event.ui_element == create_button or event.ui_element == join_button:
                print('Começa o jogo')
                room_code = caixa_code.get_text()
                nickname = caixa_nome.get_text()
                print(nickname, room_code)
            if event.ui_element == next_button:
                sprites.update(False)
            if event.ui_element == previous_button:
                sprites.update(True)

    # Atualizar tela
    sprites.draw(tela)
    manager.update(relogio)
    manager.draw_ui(tela)
    pygame.display.update()
pygame.quit()
