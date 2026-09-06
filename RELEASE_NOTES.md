# Horizonte AFK 1.6.12-beta.1

Correção do Centro de Contas e novas áreas funcionais.

- Corrigido o Centro de Contas quando dados antigos vinham como System.Object[]: a leitura agora normaliza arrays e valores de servidor/porta antes de exibir os perfis.
- A persistência de contas e regras passa a gravar coleções JSON de forma consistente, inclusive quando existe apenas uma conta.
- A aba CONEXÃO agora abre um painel de chat: texto normal envia chat e /comando envia comando ao servidor pelo mesmo caminho do RakSAMP.
- Por segurança, comandos internos do RakSAMP iniciados por ! são bloqueados, há limite de tamanho e intervalo mínimo entre envios.
- A aba CONFIGURAÇÕES agora permite controlar atualizações automáticas, abertura maximizada e modo privado do painel.
- Dialogs, modo 24/7, multi-servidor, fechamento na bandeja e limpeza dos códigos de cor permanecem preservados.
- Chat no modo 24/7 requer Horizonte AFK Cloud v0.5.3 ou superior; chat local funciona diretamente nesta versão.

SHA-256: 1ca0d2855f89beca10eba7ca8fbba3e976364a93647996102d8ca1525a06508f
