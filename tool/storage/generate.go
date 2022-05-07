package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
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
	var entity string
	err := survey.AskOne(&survey.Input{
		Message: `Enter the name of the entity for which you want to generate storage`,
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
