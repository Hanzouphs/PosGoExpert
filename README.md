# Desafio Client-Server-API em Go

Este projeto aplica conceitos fundamentais da linguagem Go, como servidores HTTP, manipulação de banco de dados SQLite, consumo de APIs externas e, principalmente, o uso de **Contextos** para controle de timeout.

## Requisitos do Desafio

- **Server.go:**
  - Endpoint `/cotacao` na porta `8080`.
  - Consumo da API de cotação (USD-BRL).
  - Timeout de **200ms** para a API de cotação.
  - Persistência dos dados no SQLite com timeout de **10ms**.
  - Retorno do valor do "bid" (valor) para o cliente.

- **Client.go:**
  - Requisição ao servidor local.
  - Timeout de **300ms** para receber a resposta.
  - Gravação da cotação no arquivo `cotacao.txt`.

## Tecnologias Utilizadas

- **Go** (Golang)
- **SQLite** (Driver: `modernc.org/sqlite`)
- **Context Package** (Gerenciamento de limites de tempo)

## Como Executar

1. **Inicie o Servidor:**
   ```bash
   go run server.go
2. **Em Outro Terminal, Inicie o Cliente**
  ```bash
  go run client.go
