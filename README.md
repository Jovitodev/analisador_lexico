# JackCompiler

Projeto acadêmico desenvolvido em **Go (Golang)** para implementar, por etapas,
um compilador da linguagem Jack, baseada no projeto nand2tetris. O primeiro
componente a ser desenvolvido será o analisador léxico.

## Etapa 1 — Estrutura inicial

Esta etapa prepara o módulo Go, o ponto de entrada do programa e o pacote `lexer`.
O pacote define `TokenType`, as cinco categorias de token, a estrutura `Token`
com os campos `Type` e `Value`, as 21 palavras reservadas e os 19 símbolos de Jack.

O analisador léxico ainda não está implementado. Leitura de arquivos `.jack`,
tratamento de comentários, reconhecimento de tokens, geração XML e testes com
os arquivos oficiais serão adicionados nas próximas etapas.

## Estrutura

```text
.
├── cmd/
│   └── jackcompiler/
│       └── main.go
├── lexer/
│   └── token.go
├── tests/
│   └── .gitkeep
├── output/
│   └── .gitkeep
├── go.mod
├── README.md
└── .gitignore
```

`cmd/jackcompiler` contém o programa executável e `lexer` contém as definições
léxicas. `tests` fica reservado aos futuros testes e `output` às saídas geradas.
Os arquivos `.gitkeep` permitem preservar essas pastas vazias no Git.
Pacotes para as próximas etapas, como o parser, poderão ser adicionados quando
forem necessários.

## Executar

Requisito: Go 1.22 ou superior instalado e disponível no PATH.
O projeto utiliza apenas a biblioteca padrão, sem dependências externas.
Execute os comandos na pasta que contém `go.mod`:

```sh
go run ./cmd/jackcompiler
```

Saída esperada:

```text
JackCompiler iniciado — Etapa 1: estrutura inicial do projeto.
```

## Verificar

```sh
go fmt ./...
go vet ./...
go build ./...
```

A Etapa 2 adicionará a implementação do analisador léxico.
