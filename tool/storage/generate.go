package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	modelPath     = `internal/model`
	tmplPath      = `tool/storage/*.go.tmpl`
	interfaceTmpl = `interface.go.tmpl`
	s3Tmpl        = `s3.go.tmpl`
	s3TestTmpl    = `s3_test.go.tmpl`
	outPath       = `internal/persistence/`
)

var tmpl = template.Must(template.ParseGlob(tmplPath))

type storageGenData struct {
	Entity      string
	EntityLower string
}

func main() {
	modelTypes, err := getModelTypes()
	if err != nil {
		panic(err)
	}

	var entity string
	err = survey.AskOne(&survey.Select{
		Message: `Enter the name of the entity for which you want to generate storage`,
		Options: modelTypes,
	}, &entity, survey.Required)
	if err != nil {
		panic(err)
	}

	entity = strings.TrimSpace(strings.ToLower(entity))

	data := storageGenData{
		Entity:      cases.Title(language.AmericanEnglish).String(entity),
		EntityLower: entity,
	}

	basePath := filepath.Join(outPath, entity+`s`)
	_, err = os.Stat(basePath)
	switch {
	case err == nil:
		panic(`path already exists: ` + basePath)
	case !os.IsNotExist(err):
		panic(err)
	}

	err = os.Mkdir(basePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	interfacePath := filepath.Join(basePath, entity+`s.go`)
	s3Path := filepath.Join(basePath, `s3.go`)
	s3TestPath := filepath.Join(basePath, `s3_test.go`)

	tmplMap := map[string]string{
		interfacePath: interfaceTmpl,
		s3Path:        s3Tmpl,
		s3TestPath:    s3TestTmpl,
	}

	for p, t := range tmplMap {
		_, err := os.Stat(p)
		if err == nil {
			panic(`file already exists: ` + p)
		}

		f, err := os.OpenFile(p, os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = tmpl.ExecuteTemplate(f, t, data)
		if err != nil {
			panic(err)
		}
	}
}

func getModelTypes() ([]string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, modelPath, func(fi fs.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		return nil, errors.New(`only expected to parse one package`)
	}

	var pkg *ast.Package
	for _, pkg = range pkgs {
	}

	var res []string
	for _, f := range pkg.Files {
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				if !hasIDField(structType) {
					continue
				}
				res = append(res, typSpec.Name.Name)
			}
		}
	}
	return res, nil
}

func hasIDField(s *ast.StructType) bool {
	for _, fld := range s.Fields.List {
		fieldTyp, ok := fld.Type.(*ast.Ident)
		if !ok || fieldTyp.Name != `string` {
			continue
		}
		for _, n := range fld.Names {
			if n.Name == `ID` {
				return true
			}
		}
	}
	return false
}
