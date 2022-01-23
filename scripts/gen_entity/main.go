package main

import (
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/AlecAivazis/survey.v1"
)

type templateContext struct {
	EntityNameUpper string
	EntityNameLower string
	EntityVarName   string
}

func main() {
	here, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasSuffix(here, string(filepath.Separator)+`go-cribbly`) {
		log.Fatal(`must be run from go-cribbly root`)
	}

	ctx, err := getTemplateContext()
	if err != nil {
		log.Fatal(err)
	}

	templateDir := filepath.Join(here, `scripts`, `gen_entity`, `templates`)
	tmpl, err := template.ParseGlob(filepath.Join(templateDir, `**`, `*.go.tmpl`))
	if err != nil {
		log.Fatal(err)
	}

	entityRoot := filepath.Join(here, ctx.EntityNameLower)
	err = os.Mkdir(entityRoot, os.ModePerm)
	if err != nil {
		log.Fatal(`error making directory:`, err)
	}

	err = filepath.Walk(templateDir, func(path string, info fs.FileInfo, _ error) error {
		if filepath.Ext(path) != `.tmpl` {
			mirrored := strings.ReplaceAll(path, templateDir, entityRoot)
			if err := os.Mkdir(mirrored, os.ModePerm); os.IsExist(err) || err == nil {
				return nil
			} else {
				return err
			}
		}

		_, file := filepath.Split(path)
		// name := strings.TrimSuffix(file, `.go.tmpl`)
		target := strings.TrimSuffix(strings.ReplaceAll(path, templateDir, entityRoot), `.tmpl`)

		f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		return tmpl.ExecuteTemplate(f, file, ctx)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(entityRoot, func(path string, info fs.FileInfo, _ error) error {
		if filepath.Ext(path) != `.go` {
			return nil
		}

		fset := token.NewFileSet()
		fileNode, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		return format.Node(f, fset, fileNode)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func getTemplateContext() (templateContext, error) {
	var n string
	err := survey.AskOne(&survey.Input{
		Message: `What is the name of the entity (e.g. game, player, etc.)?`,
	}, &n, survey.Required)
	if err != nil {
		return templateContext{}, err
	}
	normalized := strings.ToLower(strings.TrimSpace(n))

	return templateContext{
		EntityNameUpper: strings.Title(normalized),
		EntityNameLower: normalized,
		EntityVarName:   string(normalized[0]),
	}, nil
}
