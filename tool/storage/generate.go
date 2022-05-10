package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

	entity = strings.TrimSpace(entity)
	data := storageGenData{
		Entity:      entity,
		EntityLower: untitle(entity),
	}

	basePath := filepath.Join(outPath, snake(entity)+`s`)

	ok, err := exists(basePath)
	if err != nil {
		panic(err)
	}

	if !ok {
		err = os.Mkdir(basePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	interfacePath := filepath.Join(basePath, snake(entity)+`s.go`)
	s3Path := filepath.Join(basePath, `s3.go`)
	s3TestPath := filepath.Join(basePath, `s3_test.go`)

	tmplMap := map[string]string{
		interfacePath: interfaceTmpl,
		s3Path:        s3Tmpl,
		s3TestPath:    s3TestTmpl,
	}

	for p, t := range tmplMap {
		ok, err := exists(p)
		if err != nil {
			panic(err)
		}

		overwrite := true
		if ok {
			survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("%s exists. Overwrite?", p),
			}, &overwrite, survey.Required)
		}

		if !overwrite {
			continue
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

func exists(p string) (bool, error) {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func untitle(s string) string {
	done := false
	return strings.Map(func(r rune) rune {
		if done {
			return r
		}
		done = true
		return rune(strings.ToLower(string(r))[0])
	}, s)
}

var wordBoundaryPattern = regexp.MustCompile(`([A-Z])`)

// snake takes a CamelCase word and transforms it to snake_case
func snake(s string) string {
	return wordBoundaryPattern.ReplaceAllStringFunc(untitle(s), func(s string) string {
		return `_` + strings.ToLower(s)
	})
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
