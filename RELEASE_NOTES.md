# Horizonte AFK v1.6.55-beta.1

Corrigida a causa raiz do CONEXÕES E CHAT mostrar DESCONECTADA enquanto a sessão 24/7 continuava online.

A recarga dos perfis falhava no Windows PowerShell e apagava a lista de contas em memória. O refresh perdia a conta antes de consultar a sessão Cloud. A conversão foi corrigida e falhas de leitura agora preservam o último estado válido.

Também foi corrigida a leitura das páginas que concatenava IDs. Vínculos antigos são recuperados por nick, servidor e porta, mantendo contas de servidores diferentes separadas.

Validação concluída antes da publicação: 193,5 segundos com controles WPF e Cloud reais, 1.002 consultas, zero quedas de vínculo e novos logs recebidos. Os 22 testes de regressão passaram. A v1.6.54 apresentou a falha em 5,1 segundos no mesmo teste.

Preservadas as demais funções da versão anterior, incluindo login, reconexão, chat e proteções de identidade. Nenhuma mudança ou reinicialização do backend Cloud foi necessária.
