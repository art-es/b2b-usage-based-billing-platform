package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
)

const openapiGeneratedPath = "./internal/generated/openapi/server.gen.go"
const handlerTemplatePath = "./scripts/generate_openapi_paths/handler.go.template"
const openapiHandlersDir = "./internal/transport/http/openapi/handlers"

const openapiInterfaceName = "StrictServerInterface"

var logger = log.NewLogger(&log.Options{
	Level:  log.Debug,
	Format: log.FormatText,
})

var handlerTemplate = template.Must(template.ParseFiles(handlerTemplatePath))

type handlerTemplateData struct {
	Method  string
	Package string
	Comment string
}

func main() {
	if err := run(); err != nil {
		logger.Log(log.Error).Set("message", err.Error()).Write()
		os.Exit(1)
	}

	logger.Log(log.Info).Set("message", "done").Write()
}

func run() error {
	file, err := parser.ParseFile(
		token.NewFileSet(),
		openapiGeneratedPath, // parsing file
		nil,                  // not a buffer
		parser.ParseComments,
	)
	if err != nil {
		return fmt.Errorf("parse file %s: %w", openapiGeneratedPath, err)
	}

	head := findStrictServerInterface(file)
	if head == nil {
		return fmt.Errorf("%s not found", openapiInterfaceName)
	}

	err = processInterface(head)
	if err != nil {
		return fmt.Errorf("process: %w", err)
	}

	return nil
}

func findStrictServerInterface(file *ast.File) *ast.TypeSpec {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if typeSpec.Name.Name == openapiInterfaceName {
				return typeSpec
			}
		}
	}

	return nil
}

func processInterface(head *ast.TypeSpec) error {
	intr, ok := head.Type.(*ast.InterfaceType)
	if !ok {
		return fmt.Errorf("%s is not *ast.InterfaceType", openapiInterfaceName)
	}

	if intr.Methods == nil {
		logger.Log(log.Warning).Set("message", openapiInterfaceName+" has no methods").Write()
		return nil
	}

	for _, method := range intr.Methods.List {
		err := processMethod(method)
		if err != nil {
			return fmt.Errorf("process method: %w", err)
		}
	}

	return nil
}

func processMethod(method *ast.Field) error {
	if len(method.Names) == 0 {
		return errors.New("no name")
	}

	methodName := method.Names[0].Name
	packageName := toSnakeCase(methodName)
	comments := make([]string, 0, len(method.Doc.List))
	for _, doc := range method.Doc.List {
		comments = append(comments, doc.Text)
	}

	file, err := createHandlerFile(packageName)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	err = handlerTemplate.Execute(buf, handlerTemplateData{
		Method:  methodName,
		Package: packageName,
		Comment: strings.Join(comments, "\n"),
	})
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format: %w", err)
	}

	_, err = file.Write(out)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	file.Sync()

	logger.Log(log.Info).
		Set("message", "created new endpoint").
		Set("endpoint", methodName).
		Write()

	return nil
}

func createHandlerFile(packageName string) (*os.File, error) {
	packageDir := path.Join(openapiHandlersDir, packageName)

	ok, err := dirExists(packageDir)
	if err != nil {
		return nil, fmt.Errorf("check dir exists: %w", err)
	}
	if ok {
		// handler already exists
		return nil, nil
	}

	err = os.Mkdir(packageDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("mkdir %q: %w", packageDir, err)
	}

	handlerPath := path.Join(packageDir, "handler.go")

	file, err := os.Create(handlerPath)
	if err != nil {
		return nil, fmt.Errorf("create %q file: %w", handlerPath, err)
	}

	return file, nil
}

// Utils

func dirExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("os.Stat: %w", err)
	}

	if !info.IsDir() {
		return false, fmt.Errorf("path %q used for file", dirPath)
	}

	return true, nil
}

func toSnakeCase(in string) string {
	if len(in) == 0 {
		return ""
	}

	out := make([]rune, 1, len(in))
	out[0] = unicode.ToLower(rune(in[0]))

	for _, r := range in[1:] {
		if unicode.IsDigit(r) || unicode.IsLower(r) {
			out = append(out, r)
		} else {
			out = append(out, '_', unicode.ToLower(r))
		}
	}

	return string(out)
}
