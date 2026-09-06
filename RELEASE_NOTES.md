Horizonte AFK 1.4.0

- Novo visual e ícone inspirados no Horizonte, integrados ao painel.
- Mantidos o componente de conexão e todos os dialogs da versão 1.3.0 validada.
- Corrigida a tela de atualização que falhava antes de baixar o executável.
- Atualização no início e aviso durante o uso, com verificação SHA-256, substituição do mesmo arquivo e reinício.
- Substituição com cópia de segurança e restauração caso o reinício falhe.
- Executável completo: também prepara o componente validado em uma instalação nova.

Use somente **Horizonte_AFK.exe** como arquivo permanente.

As builds de teste 1.3.x e a prévia 1.4.0 possuem um defeito no atualizador já instalado. Nesses casos, é necessária uma única substituição manual do EXE, ou executar `Repair-HorizonteAFK.ps1 -TargetPath 'caminho do seu EXE'`. Depois disso, a versão corrigida mantém o fluxo automático. O endpoint público e o formato do manifesto foram preservados.

Validação: testes de integridade, painel WPF, atualização por versão/hash, formatos normal/gzip/partes, rejeição de download inválido e restauração após falha. O motor de conexão permanece byte a byte igual à base validada; esta publicação não representa um novo teste de login em servidor real.
