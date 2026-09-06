# Horizonte AFK 1.6.15-beta.1

Correção da identidade do nick nas sessões 24/7.

- Corrigida a restauração 24/7 que podia anexar a primeira sessão da Discloud mesmo quando ela pertencia a outro nick.
- O launcher agora restaura somente uma sessão cujo nick, servidor e porta correspondam exatamente aos campos atuais.
- Ao conectar na nuvem, o nick digitado é congelado antes da reconciliação de contas e enviado sem ser substituído por perfil antigo.
- Se a Cloud devolver uma sessão existente com outro nick, o launcher cria uma identidade de sessão isolada e tenta novamente, sem assumir a conta errada.
- O status seguro passa a registrar Nick solicitado e Nick enviado para facilitar a conferência sem expor senha ou 2FA.
- Modo local, Centro de Contas, chat, dialogs, multi-servidor e persistência 24/7 foram preservados.

SHA-256: 943090e36bc9dfac27064642b0b49e486f0161de21b44c5a9bc1ef9d51f51555
