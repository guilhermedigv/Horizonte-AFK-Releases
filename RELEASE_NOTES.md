# Horizonte AFK 1.6.1

Correção focada na bandeja do Windows e no uso de uma conta adicional durante o modo multi-servidor, preservando o núcleo de conexão, dialogs e updater já validados.

- Corrigida a bandeja do sistema para o ícone do Horizonte AFK permanecer acessível enquanto o aplicativo continua em execução.
- A integração da bandeja agora mantém uma cópia própria do ícone e se registra novamente no Explorer quando necessário.
- Falhas secundárias no menu da bandeja não removem mais o ícone de um aplicativo que continua aberto.
- No modo multi-servidor, o botão de conexão, seleção de servidor e campo de nick permanecem disponíveis.
- É possível conectar uma conta adicional com outro nick sem encerrar as sessões multi-servidor já ativas.
- Dialogs da conta adicional têm prioridade na tela para evitar sobreposição com login/2FA das sessões múltiplas.
- Mantidos tela cheia/F11, painel único de sessões, atualização automática com SHA-256 e núcleo de conexão/dialogs.

Quem estiver na v1.6.0 receberá a v1.6.1 pelo atualizador automático após a publicação.
