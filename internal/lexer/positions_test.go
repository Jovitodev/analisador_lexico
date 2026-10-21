package lexer

import "testing"

func TestPositions(t *testing.T) {
	l := New("// comentário\r\n/* bloco\ncontinua */\r\n  let x = 1;")
	tok := l.NextToken()
	if tok.Literal != "let" || tok.Line != 4 || tok.Column != 3 {
		t.Fatalf("posição errada: %+v", tok)
	}
	tok = l.NextToken()
	if tok.Literal != "x" || tok.Column != 7 {
		t.Fatalf("posição errada: %+v", tok)
	}
}

func TestMalformedTokens(t *testing.T) {
	for _, source := range []string{"\"sem fechamento", "\"linha\nseguinte\"", "/* aberto", "32768", "999999999999999999999999999999999999", "\x00"} {
		if tok := New(source).NextToken(); tok.Type != ILLEGAL {
			t.Errorf("%q: %+v", source, tok)
		}
	}
}

func TestStringAndEOF(t *testing.T) {
	l := New("\"// /* <&>\" 32767")
	if tok := l.NextToken(); tok.Type != STRING_CONST || tok.Literal != "// /* <&>" {
		t.Fatal(tok)
	}
	if tok := l.NextToken(); tok.Type != INT_CONST || tok.Literal != "32767" {
		t.Fatal(tok)
	}
	for i := 0; i < 3; i++ {
		if tok := l.NextToken(); tok.Type != EOF {
			t.Fatal(tok)
		}
	}
}
