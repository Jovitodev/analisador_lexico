# JackCompiler — análise léxica e sintática

Projeto acadêmico em **Go**, para a disciplina de Compiladores. O scanner
alimenta um parser por descida recursiva, que verifica a gramática Jack e gera
uma árvore sintática XML compatível com o Projeto 10 do nand2tetris.
Não há geração de código VM nesta etapa.

## Equipe

- **Daniel Nunes Duarte** — matrícula: a preencher pela equipe.
- **João Victor Oliveira** — matrícula: a preencher pela equipe.
- Repositório da dupla: inserir o endereço confirmado antes da entrega.

## Requisitos e execução

Go **1.22 ou superior**, conforme `go.mod`; nenhuma dependência externa.
Execute os comandos na raiz deste repositório, onde estão este README e go.mod.
A subpasta antiga `jackcompiler/`, caso exista, é uma cópia anterior ignorada
pelo Git; o projeto integrado está na raiz e em `internal/`.

```sh
go run ./cmd/jackcompiler caminho/Main.jack
```

O arquivo de saída padrão é **MainP.xml**, na mesma pasta de Main.jack.
É possível escolher o caminho de saída:

```sh
go run ./cmd/jackcompiler "caminho com espaços/Main.jack" "output/MainP.xml"
```

Exemplo reproduzível com o material oficial incluído:

```sh
go run ./cmd/jackcompiler internal/parser/testdata/Square/Main.jack output/MainP.xml
```

A pasta de saída deve existir. Um arquivo XML existente no caminho indicado é
substituído após a análise completa. Entrada e saída devem ser arquivos distintos.
O programa recebe exatamente um arquivo `.jack` e uma saída `.xml` opcional;
não processa pastas. Erros são apresentados em stderr, com código de saída não
zero. Erros léxicos/sintáticos incluem linha e coluna. Não se grava XML parcial
quando a análise falha, nem se anuncia sucesso nesse caso.

Para compilar um executável na pasta atual:

```sh
go build -o jackcompiler.exe ./cmd/jackcompiler
```

No Windows:

```powershell
.\jackcompiler.exe "caminho\Main.jack" "output\MainP.xml"
```

Em Linux/macOS, use `go build -o jackcompiler ./cmd/jackcompiler` e execute
`./jackcompiler caminho/Main.jack output/MainP.xml`.

## Organização

```text
cmd/jackcompiler/          Interface de linha de comando e testes de integração
internal/lexer/            Scanner reutilizado, tokens e testes léxicos
internal/parser/           Parser por descida recursiva e testes sintáticos
internal/parser/testdata/  Fontes Jack e XMLs oficiais de referência
output/                   Saídas locais geradas (ignoradas pelo Git)
go.mod                    Módulo jackcompiler, somente biblioteca padrão
```

O scanner mantém `New(input)` e `NextToken() Token`; `Token` contém `Type`,
`Literal`, `Line` e agora `Column`. A integração acrescentou posições e a rejeição
de strings/comentários não fechados e inteiros fora de 0..32767. Colunas são
contadas em bytes, a partir de 1; tabs ocupam uma posição. LF, CRLF e CR são
reconhecidos como quebras de linha.

O parser recebe o scanner por `parser.New(lexer.New(source)).Parse()` e retorna
XML ou o primeiro erro. Usa `advance`, `peek`, `match` e `consume` para percorrer
os tokens, com funções para as regras gramaticais de classes, declarações,
parâmetros, variáveis, comandos, expressões, termos e chamadas.

## Gramática e XML

São aceitos `static`, `field`, `constructor`, `function`, `method`, declarações
locais, `let` (incluindo arrays), `if`/`else`, `while`, `do` e `return`.
Expressões incluem constantes, arrays, chamadas locais/qualificadas, os nove
operadores binários e os unários `-` e `~`.

A gramática Jack define `expression = term (op term)*`: o XML preserva essa
sequência, sem impor a precedência aritmética de outras linguagens. Parênteses
criam expressões internas e operadores unários se ligam ao próximo termo.
A sequência binária conserva a ordem da esquerda para a direita; o parser
não calcula valores nem gera código.

As tags seguem o formato oficial: `class`, `classVarDec`, `subroutineDec`,
`parameterList`, `subroutineBody`, `varDec`, `statements`, `letStatement`,
`ifStatement`, `whileStatement`, `doStatement`, `returnStatement`, `expression`,
`term` e `expressionList`, além dos cinco tipos de terminal.
Regras como `type`, `op`, `unaryOp`, `keywordConstant`, `statement` e
`subroutineCall` têm funções auxiliares, mas **não recebem tags extras**,
pois isso alteraria a hierarquia dos XMLs oficiais. O agrupamento dos comandos
é representado por `<statements>` e pelas tags específicas de cada comando.
Os caracteres `&`, `<`, `>` e aspas são escapados na saída XML.

A análise é sintática: não verifica declaração de variáveis, resolução de
chamadas, tipos, retorno obrigatório ou outras regras semânticas.

## Testes e referências

```sh
go fmt ./...
go test ./...
go test ./... -cover
go vet ./...
go build ./...
```

Os testes comparam o XML decodificado, preservando hierarquia, ordem e conteúdo
dos terminais, e ignorando indentação e espaços de apresentação. Os fixtures
foram extraídos do ZIP nand2tetris fornecido pela equipe, sem usar código de
implementação de repositórios de referência.

Referência: https://www.nand2tetris.org/project10

O pacote fornecido nomeia as referências sintáticas como `Main.xml`,
`Square.xml` e `SquareGame.xml`; elas correspondem aos `*P.xml` citados na
atividade. Os arquivos `*T.xml` representam tokens e não são referências da
árvore sintática. As saídas geradas usam o sufixo P para evitar sobrescrever
as referências.

Status da validação local em 21/10/2026: **7 de 7 arquivos oficiais aprovados**
na comparação estrutural automatizada: Main, Square e SquareGame de `Square`;
os três arquivos de `ExpressionLessSquare`; e Main de `ArrayTest`.
`go fmt ./...`, `go vet ./...`, `go build ./...` e os testes passaram com o
SDK oficial Go 1.27.3, baixado em `.tools/` e verificado por SHA-256.
A diretiva `go 1.22` do módulo foi preservada; o teste local utilizou 1.27.3.
Não foi executado o TextComparer Java; a validação utiliza `encoding/xml` do Go.

Se o Go não estiver no PATH, a cópia local pode ser usada no PowerShell:

```powershell
.\.tools\go\bin\go.exe test ./...
.\.tools\go\bin\go.exe run ./cmd/jackcompiler internal/parser/testdata/Square/Main.jack output/MainP.xml
```

Além das referências oficiais, há testes para classe vazia, expressões,
parênteses/unários, escape XML, erros léxicos/sintáticos, posições e fluxo da
linha de comando. Os testes da interface usam diretórios temporários.

## Desafios da unidade

A integração exigiu preservar a API do scanner, introduzir posições de coluna
e impedir que erros léxicos fossem interpretados como fim normal da entrada.
Outro cuidado foi distinguir a gramática Jack da precedência aritmética de
outras linguagens e manter a hierarquia XML oficial, inclusive listas vazias.
A comparação automatizada dos eventos XML permite verificar essas decisões
sem depender da indentação escolhida.

## Antes da entrega

Complete as matrículas e o endereço do repositório. Confirme o histórico real
da dupla (a atividade exige pelo menos 15 commits coerentes) e execute novamente
os testes no ambiente da equipe. Esta implementação não cria commits nem altera
o histórico. Para o ZIP da entrega, inclua os fontes e os fixtures, excluindo
`.git/`, `.tools/`, binários, saídas geradas e a cópia antiga `jackcompiler/`.
