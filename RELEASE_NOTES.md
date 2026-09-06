# Horizonte AFK 1.6.9-beta.1

Hotfix para a regressão de conexão da 1.6.8-beta.1.

- Corrigida regressão da 1.6.8-beta.1 que podia impedir até a conexão local antes de iniciar o RakSAMP.
- Falhas ao registrar/sincronizar uma conta no Centro de Contas não interrompem mais o fluxo de conexão local.
- Em caso de nova falha ao conectar, o painel e o log agora mostram a causa real em vez da mensagem genérica.
- Mantidas as correções de persistência 24/7, reconciliação de contas e fila multi-servidor da versão anterior.
- A versão foi elevada para 1.6.9-beta.1 para impedir o updater antigo de interpretar betas do mesmo 1.6.8 como equivalentes e regredir o executável.

SHA-256: 337e466a2fdabe6b7d2dae140837c19c7ce29698d0e685cfd0e21568f37a2571
