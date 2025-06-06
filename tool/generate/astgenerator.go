package main

import (
	"os"
	"strings"
	"text/template"
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

type {{.BaseName}}[R any] interface {
	Accept(v Visitor[R]) R
}

type Visitor[R any] interface {
{{- range .Types }}
	Visit{{.Class}}(expr *{{.Class}}[R]) R
{{- end }}
}

{{range .Types}}
type {{.Class}}[R any] struct {
{{- range .Fields }}
	{{ .Name }} {{ .Type }}
{{- end }}
}

func (n *{{.Class}}[R]) Accept(v Visitor[R]) R {
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
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	t.Execute(file, data)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func main() {
	defineAst("ast", "Expr", []string{
		"Binary   : Expr[R] left, lexer.Token operator, Expr[R] right",
		"Grouping : Expr[R] expression",
		"Literal  : any value",
		"Unary    : lexer.Token operator, Expr[R] right",
	})
}
