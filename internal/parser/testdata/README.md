# Referências de validação

Arquivos extraídos de `nand2tetris/projects/10` do pacote oficial fornecido pelo
usuário. São material de teste, não implementação do compilador.

Fonte: https://www.nand2tetris.org/project10

- `Square`: Main, Square e SquareGame, com seus XMLs sintáticos oficiais.
- `ExpressionLessSquare`: as mesmas três classes sem expressões complexas.
- `ArrayTest`: exemplo de acesso e atribuição em arrays.

Neste pacote, as referências sintáticas se chamam `Main.xml`, `Square.xml` e
`SquareGame.xml`, equivalentes aos arquivos `*P.xml` mencionados na atividade.
Os arquivos `*T.xml` de tokens não são usados nesta comparação.

Os testes decodificam o XML e comparam hierarquia, ordem e conteúdo dos terminais,
desconsiderando indentação e o espaço de apresentação em volta dos terminais.
