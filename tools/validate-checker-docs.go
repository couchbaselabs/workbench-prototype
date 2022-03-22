// Copyright (C) 2021 Couchbase, Inc.
//
// Use of this software is subject to the Couchbase Inc. License Agreement
// which may be found at https://www.couchbase.com/LA03012021.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func walkDeclarations(comp *ast.CompositeLit, walker func(*ast.CompositeLit)) {
	for _, el := range comp.Elts {
		if kv, ok := el.(*ast.KeyValueExpr); ok {
			if cl, ok := kv.Value.(*ast.CompositeLit); ok {
				walker(cl)
			}
		}
	}
}

func findNamedField(comp *ast.CompositeLit, name string) *ast.KeyValueExpr {
	for _, elt := range comp.Elts {
		if kve, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kve.Key.(*ast.Ident); ok {
				if ident.Name == name {
					return kve
				}
			}
		}
	}
	return nil
}

// documentedCheckerIDRe matches a checker ID (CB12345), preceded either by a hash or a whitespace.
// This is so it matches either AsciiDoc ID syntax (like [#CB12345]), or comments (for checkers documented elsewhere).
var documentedCheckerIDRe = regexp.MustCompile(`[#\s](CB[0-9]+)`)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file := filepath.Join(cwd, "cluster-monitor", "pkg", "values", "checker_defs.go")
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	var declarations *ast.ValueSpec
	for _, decl := range tree.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok {
			if vs, ok := gd.Specs[0].(*ast.ValueSpec); ok {
				if vs.Names[0].Name == "AllCheckerDefs" {
					declarations = vs
					break
				}
			}
		}
	}
	if declarations == nil {
		panic("Could not find declaration")
	}
	vals := declarations.Values[0].(*ast.CompositeLit)
	problems := make([]string, 0)
	knownCheckers := make(map[string]string)
	walkDeclarations(vals, func(cl *ast.CompositeLit) {
		var name string
		if nameField := findNamedField(cl, "Name"); nameField == nil {
			name = "(unknown name)"
		} else {
			switch val := nameField.Value.(type) {
			case *ast.Ident:
				name = val.Name
			case *ast.BasicLit:
				name = val.Value
			default:
				name = "(unknown name)"
			}
		}
		idField := findNamedField(cl, "ID")
		if idField == nil {
			problems = append(problems, fmt.Sprintf("Checker %s has no ID field", name))
			return
		}
		id := strings.Replace(idField.Value.(*ast.BasicLit).Value, `"`, "", 2)
		knownCheckers[id] = name
	})

	docsPath := filepath.Join(cwd, "docs", "modules", "ROOT", "pages", "checkers.adoc")
	docs, err := ioutil.ReadFile(docsPath)
	if err != nil {
		panic(err)
	}

	documentedCheckers := make(map[string]bool)
	matches := documentedCheckerIDRe.FindAllStringSubmatch(string(docs), -1)
	if matches == nil {
		fmt.Println("Found no documented checkers at all (that can't be right!)")
		os.Exit(1)
	}
	for _, match := range matches {
		documentedCheckers[match[1]] = true
	}

	for id := range knownCheckers {
		if _, ok := documentedCheckers[id]; !ok {
			problems = append(problems, fmt.Sprintf("Checker %s is not documented."+
				"If it's documented in couchbaselabs/observability, add a comment with its ID to checkers.adoc.", id))
		}
	}
	for id := range documentedCheckers {
		if _, ok := knownCheckers[id]; !ok {
			problems = append(
				problems,
				fmt.Sprintf("Checker %s is documented but not known.", id),
			)
		}
	}

	if len(problems) == 0 {
		os.Exit(0)
	}
	fmt.Println("Problems: ")
	for _, prob := range problems {
		fmt.Printf("\t%s\n", prob)
	}
	os.Exit(1)
}
