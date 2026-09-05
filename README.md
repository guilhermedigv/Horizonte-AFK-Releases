# Horizonte AFK — Releases

Canal público de distribuição do **Horizonte AFK**.

O usuário final deve manter apenas um arquivo: **`Horizonte_AFK.exe`**.

## Atualização

A partir da linha 1.2.2, o launcher usa estratégia `replace_in_place`: quando uma nova versão é publicada, ela é baixada para a pasta temporária do Windows, validada por SHA-256 e então substitui o mesmo `Horizonte_AFK.exe`. O temporário é removido e o launcher reinicia.

Isso evita acumular arquivos como `v1.2.0.exe`, `v1.2.1.exe`, `v1.2.2.exe` na pasta do usuário.

## Arquivos públicos

- `latest.json` — informa a versão estável mais recente.
- `changelog.json` — histórico público de alterações.
- `updater-config.json` — política pública do atualizador.

## Segurança

Nenhum token ou credencial do GitHub é distribuído no programa. Toda atualização deve ser validada pelo SHA-256 publicado antes da substituição do executável.
