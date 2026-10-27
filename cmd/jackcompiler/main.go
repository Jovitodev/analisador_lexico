package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"jackcompiler/internal/lexer"
	"jackcompiler/internal/parser"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Erro:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("uso: jackcompiler entrada.jack [saida.xml]")
	}
	input := args[0]
	if filepath.Ext(input) != ".jack" {
		return fmt.Errorf("a entrada deve ter extensão .jack")
	}
	output := strings.TrimSuffix(input, ".jack") + "P.xml"
	if len(args) == 2 {
		output = args[1]
	}
	if filepath.Ext(output) != ".xml" {
		return fmt.Errorf("a saída deve ter extensão .xml")
	}
	// Impede sobrescrever a entrada, inclusive por um hard link ou link simbólico.
	inputInfo, err := os.Stat(input)
	if err != nil {
		return fmt.Errorf("acessar %q: %w", input, err)
	}
	if !inputInfo.Mode().IsRegular() {
		return fmt.Errorf("a entrada deve ser um arquivo regular")
	}
	if outputInfo, err := os.Stat(output); err == nil && os.SameFile(inputInfo, outputInfo) {
		return fmt.Errorf("entrada e saída devem ser arquivos diferentes")
	}
	content, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("ler %q: %w", input, err)
	}
	xml, err := parser.New(lexer.New(string(content))).Parse()
	if err != nil {
		return fmt.Errorf("%s: %w", input, err)
	}
	// Somente escreve quando toda a classe foi validada.
	if err := os.WriteFile(output, []byte(xml), 0644); err != nil {
		return fmt.Errorf("gravar %q: %w", output, err)
	}
	fmt.Fprintf(stdout, "Análise sintática concluída: %s\n", output)
	return nil
}
