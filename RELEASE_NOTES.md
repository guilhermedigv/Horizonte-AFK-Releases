# Horizonte AFK 1.6.8-beta.1

Atualização de correção para o modo 24/7, contas e multi-servidor.

- Centro de Contas volta a listar e reconciliar contas salvas com sessões Cloud ativas.
- Sessões 24/7 continuam na Discloud ao fechar ou sair do launcher; somente DESCONECTAR/ENCERRAR envia stop explícito.
- Multi-servidor agora processa login, 2FA e dialogs em fila, um servidor por vez.
- Ao concluir ou falhar uma sessão da fila, o próximo servidor é iniciado automaticamente.
- State.json antigo de sessão encerrada não bloqueia mais a fila multi-servidor.
- Atualização do Centro de Contas foi limitada para evitar consultas excessivas à Cloud, mantendo o botão Atualizar imediato.
- Compatível com Horizonte AFK Cloud Discloud DIRECT v0.5.2.

SHA-256: 81d6bc3ea2f8fa05adb6c32cdcd62d95cafe9436ee21399017d134b793730ac6
