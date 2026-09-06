# Horizonte AFK 1.6.3

Correção emergencial da regressão de abertura introduzida na v1.6.2.

- Corrigido o host próprio do Horizonte AFK para executar o painel como arquivo .ps1 real, preservando $PSScriptRoot e os caminhos internos usados pelo aplicativo.
- Corrigido o caso em que o Horizonte_AFK.exe abria e encerrava o host sem exibir nenhuma janela.
- Adicionado fallback automático: se o host próprio não funcionar em um Windows específico, o launcher abre o painel pelo Windows PowerShell tradicional em vez de ficar sem interface.
- O host agora grava um log interno de erro quando falhar, facilitando diagnóstico sem deixar o usuário sem resposta.
- Mantidas as correções de bandeja da v1.6.2, o modo multi-servidor, conta adicional, tela cheia, login, 2FA, dialogs e updater com SHA-256.
- Núcleo de conexão e RakSAMP preservados.

A v1.6.2 consegue receber a v1.6.3 pelo atualizador que roda antes da abertura do painel.
