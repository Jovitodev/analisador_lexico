package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "Classe com espaços.jack")
	if err := os.WriteFile(input, []byte("class Demo {}"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, explicit := range []bool{false, true} {
		output := strings.TrimSuffix(input, ".jack") + "P.xml"
		args := []string{input}
		if explicit {
			output = filepath.Join(dir, "saida.xml")
			args = append(args, output)
		}
		var stdout bytes.Buffer
		if err := run(args, &stdout); err != nil {
			t.Fatal(err)
		}
		content, err := os.ReadFile(output)
		if err != nil || !strings.Contains(string(content), "<class>") {
			t.Fatalf("XML ausente: %v", err)
		}
		if !strings.Contains(stdout.String(), output) {
			t.Fatal("sucesso sem caminho de saída")
		}
	}
	for _, args := range [][]string{
		nil, {input, "a.xml", "extra"}, {"a.txt"}, {filepath.Join(dir, "inexistente.jack")},
		{input, "saida.jack"}, {dir}, {input, filepath.Join(dir, "ausente", "saida.xml")},
	} {
		if err := run(args, &bytes.Buffer{}); err == nil {
			t.Errorf("esperado erro para %v", args)
		}
	}
}

func TestSyntaxErrorDoesNotOverwriteOutput(t *testing.T) {
	for _, source := range []string{"class Bad {", "class Bad { @ }"} {
		dir := t.TempDir()
		input, output := filepath.Join(dir, "Bad.jack"), filepath.Join(dir, "BadP.xml")
		if err := os.WriteFile(input, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(output, []byte("anterior"), 0600); err != nil {
			t.Fatal(err)
		}
		var stdout bytes.Buffer
		if err := run([]string{input}, &stdout); err == nil {
			t.Fatal("esperado erro")
		}
		content, err := os.ReadFile(output)
		if err != nil || string(content) != "anterior" || stdout.Len() != 0 {
			t.Fatal("erro sobrescreveu saída ou anunciou sucesso")
		}
	}
}
