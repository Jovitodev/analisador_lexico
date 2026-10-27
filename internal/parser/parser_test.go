package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"jackcompiler/internal/lexer"
)

func parse(source string) (string, error) { return New(lexer.New(source)).Parse() }

// Eventos XML normalizados preservam hierarquia, ordem e conteúdo dos terminais.
// Retiramos somente o espaço de apresentação que envolve cada valor terminal.
func xmlEvents(input string) ([]string, error) {
	decoder := xml.NewDecoder(strings.NewReader(input))
	var events []string
	var text strings.Builder
	terminal := false
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return events, nil
		}
		if err != nil {
			return nil, err
		}
		switch item := token.(type) {
		case xml.StartElement:
			events = append(events, "<"+item.Name.Local+">")
			switch item.Name.Local {
			case "keyword", "symbol", "identifier", "integerConstant", "stringConstant":
				terminal = true
				text.Reset()
			}
		case xml.EndElement:
			if terminal {
				value := strings.TrimSuffix(strings.TrimPrefix(text.String(), " "), " ")
				events = append(events, "text:"+value)
				terminal = false
			}
			events = append(events, "</"+item.Name.Local+">")
		case xml.CharData:
			if terminal {
				text.Write(item)
			}
		}
	}
}

func TestStringWhitespacePreserved(t *testing.T) {
	output, err := parse(`class A { function void f() { do f("  "); return; } }`)
	if err != nil {
		t.Fatal(err)
	}
	events, err := xmlEvents(output)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, event := range events {
		if event == "text:  " {
			found = true
		}
	}
	if !found {
		t.Fatal("espaços do conteúdo da string foram perdidos")
	}
}

func TestOfficialFixtures(t *testing.T) {
	paths, err := filepath.Glob("testdata/*/*.jack")
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 7 {
		t.Fatalf("esperados 7 arquivos oficiais, encontrados %d", len(paths))
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			actual, err := parse(string(source))
			if err != nil {
				t.Fatal(err)
			}
			expected, err := os.ReadFile(strings.TrimSuffix(path, ".jack") + ".xml")
			if err != nil {
				t.Fatal(err)
			}
			actualEvents, err := xmlEvents(actual)
			if err != nil {
				t.Fatal(err)
			}
			expectedEvents, err := xmlEvents(string(expected))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actualEvents, expectedEvents) {
				t.Fatal("hierarquia ou conteúdo XML difere da referência oficial")
			}
		})
	}
}

func TestEmptyClass(t *testing.T) {
	actual, err := parse("class Empty {}")
	if err != nil {
		t.Fatal(err)
	}
	expected := "<class>\n  <keyword> class </keyword>\n  <identifier> Empty </identifier>\n  <symbol> { </symbol>\n  <symbol> } </symbol>\n</class>\n"
	if actual != expected {
		t.Fatalf("XML inesperado:\n%s", actual)
	}
}

func TestGrammarAndEscaping(t *testing.T) {
	source := `class Demo {
static int count;
field Array data;
constructor Demo new(int size, char c, boolean enabled, Demo other) {
var int x, y;
let data = Array.new(size);
let x = 1 + 2 * 3 - 4 / 2;
let y = data[(x + 1)];
if ((x < y) & true | false) { let data[x] = -~x; }
else { while (x > 0) { let x = x - 1; } }
do Output.printString("<&> // /* text */");
do helper(null, this, "");
return this;
}
function void helper(Demo a, Demo b, String s) { return; }
}`
	output, err := parse(source)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := xmlEvents(output); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"&lt;", "&gt;", "&amp;", "<parameterList>", "<expressionList>", "<whileStatement>", "<stringConstant>  </stringConstant>"} {
		if !strings.Contains(output, text) {
			t.Errorf("saída sem %s", text)
		}
	}
	for _, tag := range []string{"<op>", "<statement>", "<type>", "<subroutineCall>"} {
		if strings.Contains(output, tag) {
			t.Errorf("tag incompatível com XML oficial: %s", tag)
		}
	}
}

func TestExpressionShape(t *testing.T) {
	output, err := parse("class A { function int f() { return 1 + 2 * 3; } }")
	if err != nil {
		t.Fatal(err)
	}
	// A gramática oficial mantém todos os operadores binários nesta expressão.
	if strings.Count(output, "<expression>") != 1 || strings.Count(output, "<term>") != 3 {
		t.Fatalf("agrupamento binário incompatível com Jack:\n%s", output)
	}
	output, err = parse("class A { function int f() { return -(1 + 2) * 3; } }")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(output, "<expression>") != 2 {
		t.Fatal("parênteses não criaram expressão interna")
	}
}

func TestInvalidSyntax(t *testing.T) {
	tests := []string{
		"", "class {}", "class A {", "class A {} class B {}",
		"class A { field void x; }",
		"class A { function void f(int) { return; } }",
		"class A { function void f() { let x = ; } }",
		"class A { function void f() { return } }",
		"class A { function void f() { do f(1,); } }",
		"class A { function void f() { do f; } }",
		"class A { function void f() { let a[1 = 2; } }",
		"class A { function void f() { let x = 1 2; } }",
		"class A { function void f() { let x = class; } }",
		"class A { function void f() { return; var int x; } }",
		"class A { @ }", "class A {} /* sem fechamento",
		`class A { function void f() { do f("incompleta); } }`,
		"class A { function int f() { return 32768; } }",
	}
	for i, source := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			output, err := parse(source)
			var syntax *SyntaxError
			if !errors.As(err, &syntax) {
				t.Fatalf("esperado SyntaxError, recebido %v", err)
			}
			if output != "" || syntax.Line < 1 || syntax.Column < 1 {
				t.Fatalf("erro sem posição ou XML parcial: %v", err)
			}
		})
	}
}

func TestErrorPosition(t *testing.T) {
	_, err := parse("class A {\r\n  function void f() {\r\n    return\r\n  }\r\n}")
	var syntax *SyntaxError
	if !errors.As(err, &syntax) || syntax.Line != 4 || syntax.Column != 3 {
		t.Fatalf("posição esperada 4:3, obtida %v", err)
	}
}
