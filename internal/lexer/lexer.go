package lexer

import "strconv"

// Lexer percorre a entrada byte a byte, preservando a API do scanner original.
type Lexer struct {
	input                  string
	position, readPosition int
	ch                     byte
	line, column           int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.ch == '\n' || (l.ch == '\r' && l.peekChar() != '\n') {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	l.position = l.readPosition
	if l.position >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.position]
	}
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) atEnd() bool { return l.position >= len(l.input) }

func (l *Lexer) skipWhitespace() {
	for !l.atEnd() && (l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r') {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	// Iteração evita acumular chamadas recursivas para comentários consecutivos.
	for {
		l.skipWhitespace()
		tok := Token{Line: l.line, Column: l.column}
		if l.atEnd() {
			tok.Type = EOF
			return tok
		}
		start := l.position
		if l.ch == '/' && l.peekChar() == '/' {
			for !l.atEnd() && l.ch != '\n' && l.ch != '\r' {
				l.readChar()
			}
			continue
		}
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar()
			l.readChar()
			for !l.atEnd() && !(l.ch == '*' && l.peekChar() == '/') {
				l.readChar()
			}
			if l.atEnd() {
				tok.Type, tok.Literal = ILLEGAL, l.input[start:l.position]
				return tok
			}
			l.readChar()
			l.readChar()
			continue
		}
		if l.ch == '"' {
			l.readChar()
			for !l.atEnd() && l.ch != '"' && l.ch != '\n' && l.ch != '\r' {
				l.readChar()
			}
			if l.atEnd() || l.ch != '"' {
				tok.Type, tok.Literal = ILLEGAL, l.input[start:l.position]
				return tok
			}
			tok.Type, tok.Literal = STRING_CONST, l.input[start+1:l.position]
			l.readChar()
			return tok
		}
		if isLetter(l.ch) {
			for !l.atEnd() && (isLetter(l.ch) || isDigit(l.ch)) {
				l.readChar()
			}
			tok.Literal = l.input[start:l.position]
			tok.Type = LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			for !l.atEnd() && isDigit(l.ch) {
				l.readChar()
			}
			tok.Type, tok.Literal = INT_CONST, l.input[start:l.position]
			value, err := strconv.Atoi(tok.Literal)
			if err != nil || value > 32767 {
				tok.Type = ILLEGAL
			}
			return tok
		}
		tok.Literal = string(l.ch)
		switch l.ch {
		case '{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/', '&', '|', '<', '>', '=', '~':
			tok.Type = TokenType(tok.Literal)
		default:
			tok.Type = ILLEGAL
		}
		l.readChar()
		return tok
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }
