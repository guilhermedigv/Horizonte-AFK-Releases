# Horizonte AFK 1.6.6-beta.1

Reintrodução segura do modo Cloud 24/7.

- O botão 24/7 funciona dentro do próprio Horizonte AFK, sem redirecionar para navegador.
- A sessão RakSAMP é criada na Discloud; login, 2FA e dialogs são respondidos pelo EXE.
- Corrige a regressão da beta anterior que podia iniciar sem mostrar nenhuma janela.
- O script PowerShell/WPF final agora é analisado e carregado em teste real antes da build ser aceita.
- Foi adicionado fallback seguro: uma falha na integração Cloud não impede o modo local de abrir.
- O backend Discloud precisa estar na versão 0.4.0 com direct_client=true para ativar o modo direto.

SHA-256: 
