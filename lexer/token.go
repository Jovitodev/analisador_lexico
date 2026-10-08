// Package lexer reúne as definições usadas pelo futuro analisador léxico de Jack.
package lexer

// TokenType representa uma das cinco categorias de tokens da linguagem Jack.
type TokenType string

const (
	Keyword         TokenType = "keyword"
	Symbol          TokenType = "symbol"
	IntegerConstant TokenType = "integerConstant"
	StringConstant  TokenType = "stringConstant"
	Identifier      TokenType = "identifier"
)

// Token representa a categoria e o valor textual de um token.
type Token struct {
	Type  TokenType
	Value string
}

// Keywords contém as 21 palavras reservadas da linguagem Jack.
// Esta lista será usada pelo analisador léxico nas próximas etapas.
var Keywords = [...]string{
	"class", "constructor", "function", "method", "field", "static", "var",
	"int", "char", "boolean", "void", "true", "false", "null", "this",
	"let", "do", "if", "else", "while", "return",
}

// Symbols contém os 19 símbolos válidos da linguagem Jack.
var Symbols = [...]rune{
	'{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/',
	'&', '|', '<', '>', '=', '~',
}
