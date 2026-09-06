# Horizonte AFK 1.6.2

Correção estrutural da integração com o Windows, focada na bandeja do sistema e na identidade dos processos do Horizonte AFK, preservando o núcleo de conexão, dialogs e updater já validados.

- Corrigido o comportamento ao clicar no ×: o painel continua oculto com o AFK ativo e o acesso pela bandeja fica mais estável.
- A bandeja não remove e adiciona mais o ícone repetidamente em intervalos curtos.
- O Horizonte AFK passa a observar a mensagem TaskbarCreated do Windows e só refaz o registro do ícone quando o Explorer realmente reconstrói a área de notificação.
- O painel principal e as sessões multi-servidor passam a executar em um host próprio identificado como Horizonte AFK.
- Os processos AFK persistentes deixam de depender de powershell.exe como processo visível no Gerenciador de Tarefas.
- Mantida a possibilidade de conectar outra conta com nick diferente enquanto o modo multi-servidor continua ativo.
- Mantidos painel único, tela cheia/F11, login, 2FA, dialogs, atualização automática com SHA-256 e o núcleo de conexão validado.

Quem estiver na v1.6.1 receberá a v1.6.2 pelo atualizador automático após a publicação.
