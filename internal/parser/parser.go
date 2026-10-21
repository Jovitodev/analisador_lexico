// Package parser implementa a gramática Jack por descida recursiva e emite XML.
package parser

import (
	"fmt"
	"strings"

	"jackcompiler/internal/lexer"
)

// SyntaxError identifica o primeiro erro léxico ou sintático da entrada.
type SyntaxError struct {
	Line, Column int
	Message      string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("linha %d, coluna %d: %s", e.Line, e.Column, e.Message)
}

type Parser struct {
	scanner *lexer.Lexer
	current lexer.Token
	xml     strings.Builder
	depth   int
}

func New(scanner *lexer.Lexer) *Parser { return &Parser{scanner: scanner} }

// Parse lê exatamente uma classe. Em caso de erro, não retorna XML parcial.
// O panic interno simplifica a propagação entre regras; somente SyntaxError
// é capturado aqui. Erros de programação continuam sendo propagados.
func (p *Parser) Parse() (output string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if syntax, ok := recovered.(*SyntaxError); ok {
				output, err = "", syntax
			} else {
				panic(recovered)
			}
		}
	}()
	p.current = p.scanner.NextToken()
	p.checkLexicalError()
	p.class()
	if !p.match(lexer.EOF) {
		p.fail("fim do arquivo")
	}
	return p.xml.String(), nil
}

func (p *Parser) peek() lexer.Token { return p.current }

// match consulta a categoria sem consumir o token.
func (p *Parser) match(kind lexer.TokenType) bool { return p.peek().Type == kind }

func (p *Parser) advance() lexer.Token {
	tok := p.current
	if tok.Type == lexer.EOF {
		p.fail("token antes do fim do arquivo")
	}
	tag := "symbol"
	switch tok.Type {
	case lexer.KEYWORD:
		tag = "keyword"
	case lexer.IDENTIFIER:
		tag = "identifier"
	case lexer.INT_CONST:
		tag = "integerConstant"
	case lexer.STRING_CONST:
		tag = "stringConstant"
	}
	escaped := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&apos;").Replace(tok.Literal)
	p.line(fmt.Sprintf("<%s> %s </%s>", tag, escaped, tag))
	p.current = p.scanner.NextToken()
	p.checkLexicalError()
	return tok
}

func (p *Parser) consume(kind lexer.TokenType) lexer.Token {
	if !p.match(kind) {
		p.fail(string(kind))
	}
	return p.advance()
}

func (p *Parser) literal(value string) {
	// Literais de string e identificadores nunca substituem palavras ou símbolos.
	if (p.current.Type != lexer.KEYWORD && string(p.current.Type) != value) || p.current.Literal != value {
		p.fail(fmt.Sprintf("%q", value))
	}
	p.advance()
}

func (p *Parser) is(values ...string) bool {
	for _, value := range values {
		if p.current.Literal == value && (p.match(lexer.KEYWORD) || string(p.current.Type) == value) {
			return true
		}
	}
	return false
}

func (p *Parser) checkLexicalError() {
	if p.match(lexer.ILLEGAL) {
		panic(&SyntaxError{p.current.Line, p.current.Column, fmt.Sprintf("erro léxico: literal inválido %q", p.current.Literal)})
	}
}

func (p *Parser) fail(expected string) {
	found := fmt.Sprintf("%q (%s)", p.current.Literal, p.current.Type)
	if p.match(lexer.EOF) {
		found = "fim do arquivo"
	}
	panic(&SyntaxError{p.current.Line, p.current.Column, fmt.Sprintf("esperado %s; encontrado %s", expected, found)})
}

func (p *Parser) line(text string) { p.xml.WriteString(strings.Repeat("  ", p.depth) + text + "\n") }
func (p *Parser) open(tag string)  { p.line("<" + tag + ">"); p.depth++ }
func (p *Parser) close(tag string) { p.depth--; p.line("</" + tag + ">") }

func (p *Parser) class() {
	p.open("class")
	p.literal("class")
	p.consume(lexer.IDENTIFIER)
	p.literal("{")
	for p.is("static", "field") {
		p.classVarDec()
	}
	for p.is("constructor", "function", "method") {
		p.subroutineDec()
	}
	p.literal("}")
	p.close("class")
}

func (p *Parser) classVarDec() {
	p.open("classVarDec")
	if !p.is("static", "field") {
		p.fail("static ou field")
	}
	p.advance()
	p.typeName()
	p.consume(lexer.IDENTIFIER)
	for p.is(",") {
		p.literal(",")
		p.consume(lexer.IDENTIFIER)
	}
	p.literal(";")
	p.close("classVarDec")
}

// typeName reconhece type; esta regra não recebe tag própria no XML oficial.
func (p *Parser) typeName() {
	if p.is("int", "char", "boolean") {
		p.advance()
		return
	}
	p.consume(lexer.IDENTIFIER)
}

func (p *Parser) subroutineDec() {
	p.open("subroutineDec")
	if !p.is("constructor", "function", "method") {
		p.fail("constructor, function ou method")
	}
	p.advance()
	if p.is("void") {
		p.advance()
	} else {
		p.typeName()
	}
	p.consume(lexer.IDENTIFIER)
	p.literal("(")
	p.parameterList()
	p.literal(")")
	p.subroutineBody()
	p.close("subroutineDec")
}

func (p *Parser) parameterList() {
	p.open("parameterList")
	if !p.is(")") {
		p.typeName()
		p.consume(lexer.IDENTIFIER)
		for p.is(",") {
			p.literal(",")
			p.typeName()
			p.consume(lexer.IDENTIFIER)
		}
	}
	p.close("parameterList")
}

func (p *Parser) subroutineBody() {
	p.open("subroutineBody")
	p.literal("{")
	for p.is("var") {
		p.varDec()
	}
	p.statements()
	p.literal("}")
	p.close("subroutineBody")
}

func (p *Parser) varDec() {
	p.open("varDec")
	p.literal("var")
	p.typeName()
	p.consume(lexer.IDENTIFIER)
	for p.is(",") {
		p.literal(",")
		p.consume(lexer.IDENTIFIER)
	}
	p.literal(";")
	p.close("varDec")
}

func (p *Parser) statements() {
	p.open("statements")
	for p.is("let", "if", "while", "do", "return") {
		p.statement()
	}
	p.close("statements")
}

func (p *Parser) statement() {
	switch p.current.Literal {
	case "let":
		p.letStatement()
	case "if":
		p.ifStatement()
	case "while":
		p.whileStatement()
	case "do":
		p.doStatement()
	case "return":
		p.returnStatement()
	default:
		p.fail("comando let, if, while, do ou return")
	}
}

func (p *Parser) letStatement() {
	p.open("letStatement")
	p.literal("let")
	p.consume(lexer.IDENTIFIER)
	if p.is("[") {
		p.literal("[")
		p.expression()
		p.literal("]")
	}
	p.literal("=")
	p.expression()
	p.literal(";")
	p.close("letStatement")
}

func (p *Parser) ifStatement() {
	p.open("ifStatement")
	p.literal("if")
	p.literal("(")
	p.expression()
	p.literal(")")
	p.literal("{")
	p.statements()
	p.literal("}")
	if p.is("else") {
		p.literal("else")
		p.literal("{")
		p.statements()
		p.literal("}")
	}
	p.close("ifStatement")
}

func (p *Parser) whileStatement() {
	p.open("whileStatement")
	p.literal("while")
	p.literal("(")
	p.expression()
	p.literal(")")
	p.literal("{")
	p.statements()
	p.literal("}")
	p.close("whileStatement")
}

func (p *Parser) doStatement() {
	p.open("doStatement")
	p.literal("do")
	p.subroutineCall()
	p.literal(";")
	p.close("doStatement")
}

func (p *Parser) returnStatement() {
	p.open("returnStatement")
	p.literal("return")
	if !p.is(";") {
		p.expression()
	}
	p.literal(";")
	p.close("returnStatement")
}

// Jack define expression = term (op term)*, sem níveis de precedência binária.
// Parênteses geram expression dentro de term; unários se aplicam a um term.
func (p *Parser) expression() {
	p.open("expression")
	p.term()
	for p.is("+", "-", "*", "/", "&", "|", "<", ">", "=") {
		p.op()
		p.term()
	}
	p.close("expression")
}

func (p *Parser) op() {
	if !p.is("+", "-", "*", "/", "&", "|", "<", ">", "=") {
		p.fail("operador binário")
	}
	p.advance()
}

func (p *Parser) unaryOp() {
	if !p.is("-", "~") {
		p.fail("operador unário")
	}
	p.advance()
}

func (p *Parser) keywordConstant() {
	if !p.is("true", "false", "null", "this") {
		p.fail("true, false, null ou this")
	}
	p.advance()
}

func (p *Parser) term() {
	p.open("term")
	switch {
	case p.match(lexer.INT_CONST), p.match(lexer.STRING_CONST):
		p.advance()
	case p.is("true", "false", "null", "this"):
		p.keywordConstant()
	case p.is("-", "~"):
		p.unaryOp()
		p.term()
	case p.is("("):
		p.literal("(")
		p.expression()
		p.literal(")")
	case p.match(lexer.IDENTIFIER):
		p.consume(lexer.IDENTIFIER)
		if p.is("[") {
			p.literal("[")
			p.expression()
			p.literal("]")
		} else if p.is("(", ".") {
			p.callSuffix()
		}
	default:
		p.fail("termo (constante, identificador, chamada, unário ou expressão entre parênteses)")
	}
	p.close("term")
}

// subroutineCall e seus auxiliares emitem os terminais diretamente, como no
// padrão oficial: não existem tags statement, op, type ou subroutineCall.
func (p *Parser) subroutineCall() {
	p.consume(lexer.IDENTIFIER)
	p.callSuffix()
}

func (p *Parser) callSuffix() {
	if p.is(".") {
		p.literal(".")
		p.consume(lexer.IDENTIFIER)
	}
	p.literal("(")
	p.expressionList()
	p.literal(")")
}

func (p *Parser) expressionList() {
	p.open("expressionList")
	if !p.is(")") {
		p.expression()
		for p.is(",") {
			p.literal(",")
			p.expression()
		}
	}
	p.close("expressionList")
}
