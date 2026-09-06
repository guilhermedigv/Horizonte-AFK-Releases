# Horizonte AFK 1.6.14-beta.1

Correção do envio por ENTER na aba Conexão e Chat.

- Corrigido o fechamento inesperado do EXE ao pressionar ENTER no campo da aba CONEXÃO E CHAT.
- A tecla ENTER agora é interceptada antes da ação padrão da janela e envia pelo mesmo fluxo seguro de chat/comando.
- Removida a dependência do handler de teclado em uma variável local que já não existia quando a tecla era pressionada.
- O botão ENVIAR continua com o mesmo comportamento e as correções de chat/dialog da 1.6.13 foram preservadas.

SHA-256: a6474a1873c269c595a649dc33c9f9ea9985204a47061f3ef8d22d93f8d6023e
