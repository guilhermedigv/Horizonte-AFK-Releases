# Horizonte AFK — Releases

Canal público de distribuição do **Horizonte AFK**.

Este repositório contém somente metadados públicos de atualização e, quando publicado, o executável final. O código-fonte e a configuração interna permanecem no repositório privado de manutenção.

## Arquivos públicos

- `latest.json` — informa ao launcher a versão estável mais recente.
- `changelog.json` — histórico público de alterações.

## Segurança do updater

O launcher deve validar o `SHA-256` do executável antes de substituir a versão instalada. Nenhum token ou credencial do GitHub deve ser distribuído junto com o programa.
