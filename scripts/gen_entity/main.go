package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

func main() {
	here, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasSuffix(here, string(filepath.Separator)+`go-cribbly`) {
		log.Fatal(`must be run from go-cribbly root`)
	}

	var name string
	err = survey.AskOne(&survey.Input{
		Message: `What is the name of the entity (e.g. game, player, etc.)?`,
	}, &name, survey.Required)
	if err != nil {
		log.Fatal(err)
	}

	entityRoot := filepath.Join(here, name)
	model := filepath.Join(entityRoot, `model`)
	network := filepath.Join(entityRoot, `network`)
	repo := filepath.Join(entityRoot, `repository`)
	for _, p := range []string{
		entityRoot,
		model,
		network,
		repo,
	} {
		err := os.Mkdir(p, os.ModePerm)
		if err != nil {
			log.Fatal(`error making directory:`, err)
		}
	}
	for _, f := range []string{
		filepath.Join(model, name),
		filepath.Join(network, `handlers`),
		filepath.Join(repo, `interface`),
		filepath.Join(repo, `memory`),
	} {
		err := createFileWithPackageName(name, f+`.go`)
		if err != nil {
			log.Fatal(`error creating file:`, err)
		}
	}
}

func createFileWithPackageName(packageName, path string) error {
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, `package`, packageName)
	return err
}
