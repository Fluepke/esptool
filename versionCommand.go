package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type VersionInfo struct {
	Name        string
	Description string
	Version     string
	Authors     []string
}

func GetVersionInfo() *VersionInfo {
	return &VersionInfo{
		Name:        "Esptool",
		Description: "ESP32 flashing utility written in GoLang",
		Version:     version,
		Authors:     []string{"@fluepke"},
	}
}

func (v *VersionInfo) String() string {
	return fmt.Sprintf("%s\n\n%s\nVersion: %s\nAuthors: %s",
		v.Name,
		v.Description,
		v.Version,
		strings.Join(v.Authors, ", "),
	)
}

func versionCommand(jsonOutput bool) error {
	if jsonOutput {
		prettyJson, err := json.MarshalIndent(GetVersionInfo(), "", "  ")
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(prettyJson)
		return err
	} else {
		_, err := fmt.Print(GetVersionInfo().String())
		return err
	}
}
