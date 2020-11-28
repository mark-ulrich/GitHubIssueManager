package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Settings struct {
	EditorCommand string
}

const (
	defaultSettingsFilename = "GitHubIssueManager.conf"
)

// Load settings from a given configuration file
func loadSettings(settings *Settings, filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineno := 0
	for {
		lineno++

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) == 0 {
					break
				}
			} else {
				return err
			}
		}

		line = strings.TrimSpace(line)

		// Ignore comments
		if line[0] == '#' {
			continue
		}

		pair := strings.Split(line, "=")
		if len(pair) == 1 {
			return fmt.Errorf("Reading settings file %s: Line %d: Malformed\n", filename, lineno)
		}
		key, value := strings.ToLower(strings.TrimSpace(pair[0])), strings.TrimSpace(pair[1])
		fmt.Printf("key='%s'\nvalue='%s'", key, value)
		switch key {
		case "editor":
			settings.EditorCommand = value
		}

	}

	return nil
}
