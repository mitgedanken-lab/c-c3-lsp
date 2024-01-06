package c3

//#include "tree_sitter/parser.h"
//TSLanguage *tree_sitter_c3();
import "C"
import (
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	"unsafe"
)

func getParser() *sitter.Parser {
	parser := sitter.NewParser()
	parser.SetLanguage(GetLanguage())

	return parser
}

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_c3())
	return sitter.NewLanguage(ptr)
}

func FindIdentifiers(source string, debug bool) []string {
	parser := getParser()

	sourceCode := []byte(source)
	n := parser.Parse(nil, sourceCode)
	if debug {
		fmt.Print(n.RootNode())
	}

	// Iterate over query results
	variableIdentifiers := FindVariableDeclarations(sourceCode, n)
	functionIdentifiers := FindFunctionDeclarations(sourceCode, n)

	identifiers := append(variableIdentifiers, functionIdentifiers...)

	return identifiers
}

func FindVariableDeclarations(sourceCode []byte, n *sitter.Tree) []string {
	query := `(var_declaration (identifier) @variable_name)`
	q, err := sitter.NewQuery([]byte(query), GetLanguage())
	if err != nil {
		panic(err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n.RootNode())

	var identifiers []string
	found := make(map[string]bool)
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			content := c.Node.Content(sourceCode)
			c.Node.Parent().Type()
			if _, exists := found[content]; !exists {
				found[content] = true
				identifiers = append(identifiers, content)
			}
		}
	}

	return identifiers
}

func FindFunctionDeclarations(sourceCode []byte, n *sitter.Tree) []string {
	query := `(function_declaration name: (identifier) @function_name)`
	q, err := sitter.NewQuery([]byte(query), GetLanguage())
	if err != nil {
		panic(err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n.RootNode())

	var identifiers []string
	found := make(map[string]bool)
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			content := c.Node.Content(sourceCode)
			c.Node.Parent().Type()
			if _, exists := found[content]; !exists {
				found[content] = true
				identifiers = append(identifiers, content)
			}
		}
	}

	return identifiers
}
