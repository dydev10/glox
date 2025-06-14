package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

type Field struct {
	Name string
	Type string
}

type AstType struct {
	Class  string
	Fields []Field
}

const astTemplate = `package {{.PackageName}}

import "github.com/dydev10/glox/lexer"

type {{.BaseName}} interface {
	Accept(v Visitor{{.BaseName}}[any]) (any, error)
}

type Visitor{{.BaseName}}[R any] interface {
{{- range .Types }}
	Visit{{.Class}}(expr *{{.Class}}) (R, error)
{{- end }}
}

{{range .Types}}
type {{.Class}} struct {
{{- range .Fields }}
	{{ .Name }} {{ .Type }}
{{- end }}
}

func (n *{{.Class}}) Accept(v Visitor{{$.BaseName}}[any]) (any, error) {
	return v.Visit{{.Class}}(n)
}
{{end}}
`

func defineAst(outputDir, baseName string, types []string) {
	var astTypes []AstType

	for _, t := range types {
		parts := strings.Split(t, ":")
		className := strings.TrimSpace(parts[0])
		fieldsList := strings.Split(strings.TrimSpace(parts[1]), ",")

		var fields []Field
		for _, f := range fieldsList {
			fieldParts := strings.Fields(strings.TrimSpace(f))
			fieldType := fieldParts[0]
			fieldName := capitalize(fieldParts[1])
			fields = append(fields, Field{
				Name: fieldName,
				Type: fieldType,
			})
		}

		astTypes = append(astTypes, AstType{
			Class:  className,
			Fields: fields,
		})
	}

	data := struct {
		PackageName string
		BaseName    string
		Types       []AstType
	}{
		PackageName: "ast",
		BaseName:    baseName,
		Types:       astTypes,
	}

	t := template.Must(template.New("ast").Parse(astTemplate))

	os.MkdirAll(outputDir, os.ModePerm)
	filename := outputDir + "/" + strings.ToLower(baseName) + ".go"

	var fileBuf bytes.Buffer

	err := t.Execute(&fileBuf, data)
	if err != nil {
		panic(err)
	}

	formatted, err := imports.Process(filename, fileBuf.Bytes(), nil)
	if err != nil {
		panic(err)
	}

	os.WriteFile(filename, formatted, 0644)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func main() {
	defineAst("ast", "Expr", []string{
		"Assign   : *lexer.Token name, Expr value",
		"Binary   : Expr left, *lexer.Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : any value",
		"Unary    : *lexer.Token operator, Expr right",
		"Variable : *lexer.Token name",
	})

	defineAst("ast", "Stmt", []string{
		"Block			: []Stmt statements",
		"Expression	: Expr expression",
		"If         : Expr condition, Stmt thenBranch, Stmt elseBranch",
		"Print      : Expr expression",
		"Var        : *lexer.Token name, Expr initializer",
	})
}
