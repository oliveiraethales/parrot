// Package parrot provides a linter that detects comments which merely restate the obvious.
package parrot

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "parrot",
	Doc:  "detects comments that parrot what the code already says",
	Run:  run,
}

// Threshold for identifier overlap ratio (0.0 - 1.0)
// If this fraction of comment words match code identifiers, flag it.
const overlapThreshold = 0.4

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		lineNodes := make(map[int]ast.Node)

		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			switch n.(type) {
			case *ast.FuncDecl, *ast.AssignStmt, *ast.IfStmt, *ast.ReturnStmt, *ast.ExprStmt:
				line := pass.Fset.Position(n.Pos()).Line
				lineNodes[line] = n
			}
			return true
		})

		for _, cg := range file.Comments {
			for _, c := range cg.List {
				commentLine := pass.Fset.Position(c.Pos()).Line

				if node, ok := lineNodes[commentLine+1]; ok {
					identifiers := extractIdentifiers(node)
					if len(identifiers) > 0 && isParrotComment(c.Text, identifiers) {
						pass.Reportf(c.Pos(), "comment parrots the code: consider removing or adding insight")
						continue
					}
				}

				if node, ok := lineNodes[commentLine]; ok {
					identifiers := extractIdentifiers(node)
					if len(identifiers) > 0 && isParrotComment(c.Text, identifiers) {
						pass.Reportf(c.Pos(), "comment parrots the code: consider removing or adding insight")
					}
				}
			}
		}
	}

	return nil, nil
}

func extractIdentifiers(n ast.Node) map[string]bool {
	ids := make(map[string]bool)

	switch n.(type) {
	case *ast.IfStmt:
		ids["error"] = true
		ids["check"] = true
	case *ast.ReturnStmt:
		ids["return"] = true
	}

	ast.Inspect(n, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.Ident:
			originalName := x.Name
			name := strings.ToLower(originalName)
			if len(name) > 2 && !isBoringIdent(name) {
				ids[name] = true
				for _, word := range splitIdentifier(originalName) {
					word = strings.ToLower(word)
					if len(word) > 2 {
						ids[word] = true
					}
				}
			}
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				for _, word := range tokenizeText(x.Value) {
					if len(word) > 3 {
						ids[word] = true
					}
				}
			}
		}
		return true
	})

	return ids
}

func isParrotComment(comment string, codeIdents map[string]bool) bool {
	text := strings.TrimPrefix(comment, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		return false
	}

	if idx := strings.Index(text, "//"); idx > 0 {
		text = strings.TrimSpace(text[:idx])
	}

	words := tokenizeText(text)
	if len(words) < 2 {
		return false
	}

	matches := 0
	meaningfulWords := 0

	for _, word := range words {
		if isFillerWord(word) {
			continue
		}
		meaningfulWords++
		if codeIdents[word] {
			matches++
		}
	}

	if meaningfulWords == 0 {
		return false
	}

	ratio := float64(matches) / float64(meaningfulWords)
	return ratio >= overlapThreshold
}

func tokenizeText(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else if current.Len() > 0 {
			words = append(words, current.String())
			current.Reset()
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

func splitIdentifier(s string) []string {
	if strings.Contains(s, "_") {
		return strings.Split(s, "_")
	}

	var words []string
	var current strings.Builder

	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			if current.Len() > 0 {
				words = append(words, strings.ToLower(current.String()))
				current.Reset()
			}
		}
		current.WriteRune(r)
	}
	if current.Len() > 0 {
		words = append(words, strings.ToLower(current.String()))
	}

	return words
}

func isFillerWord(word string) bool {
	fillers := map[string]bool{
		"a": true, "an": true, "the": true, "is": true, "are": true,
		"if": true, "then": true, "else": true, "when": true, "will": true,
		"to": true, "for": true, "of": true, "in": true, "on": true,
		"and": true, "or": true, "not": true, "this": true, "that": true,
		"it": true, "we": true, "be": true, "as": true, "with": true,
		"from": true, "by": true, "at": true, "do": true, "does": true,
		"has": true, "have": true, "here": true, "there": true, "was": true,
		"handles": true, "handle": true, "processing": true, "process": true,
	}
	return fillers[word]
}

func isBoringIdent(name string) bool {
	boring := map[string]bool{
		"err": true, "nil": true, "true": true, "false": true,
		"int": true, "string": true, "bool": true, "error": true,
		"ctx": true, "ok": true, "i": true, "j": true, "k": true,
	}
	return boring[name]
}
