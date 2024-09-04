package model

import (
	"cmdkeep/file"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

type Command struct {
	Template string
	LastArgs []string
}

type Model struct {
	Last     string
	LastArgs []string
	Commands map[string]*Command
}

func NewModel() *Model {
	return &Model{
		Commands: map[string]*Command{},
		LastArgs: []string{},
	}
}

func (m *Model) AddCommand(key string, command *Command) {
	m.Commands[key] = command
}

func NewCommand(template string) *Command {
	strings.NewReplacer()
	return &Command{
		Template: template,
		LastArgs: []string{},
	}
}

func ReadModel() *Model {
	path, err := GetModelPath()
	if err != nil {
		readFailure(err)
	}

	if _, err := os.Stat(path); err != nil {
		return NewModel()
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		readFailure(err)
	}

	var m Model
	if err = json.Unmarshal(bytes, &m); err != nil {
		return NewModel()
	}
	return &m
}

func GetModelPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(userConfigDir, "CmdKeep", "model.json"), nil
}

func WriteModel(m *Model) {
	modelJson, err := json.Marshal(m)
	if err != nil {
		writeFailure(err)
	}

	modelPath, err := GetModelPath()
	if err != nil {
		writeFailure(err)
	}

	if err = file.WriteToFile(modelPath, string(modelJson)); err != nil {
		writeFailure(err)
	}
}

func readFailure(err error) {
	fmt.Fprintf(os.Stderr, "Failed to read model: %v\n", err)
	os.Exit(1)
}

func writeFailure(err error) {
	fmt.Fprintf(os.Stderr, "Failed to write model: %v\n", err)
	os.Exit(1)
}
