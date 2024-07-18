package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	usageText := "slap <path>"
	description := `slap creates either a file or a directory based on the provided path.
It also creates all necessary parent directories.

Examples:
  slap /path/to/file.txt     # Creates/updates a file
  slap /path/to/directory/   # Creates a directory (ends with /) `

	app := &cli.App{
		Name:                 "slap",
		Usage:                "Modern replacement for touch and mkdir commands",
		HideHelpCommand:      true,
		UsageText:            usageText,
		Description:          description,
		Version:              "0.1",
		EnableBashCompletion: true,

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "cd",
				Usage: "cd into the created directory",
			},
		},

		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("expected 1 path argument, but got %d", c.NArg())
			}

			pathArg := c.Args().First()
			if pathArg == "" {
				return fmt.Errorf("path argument is empty")
			}

			err := create(pathArg)
			if err != nil {
				return err
			}

			return nil
		},
	}

	genFishCompletions(app.ToFishCompletion())

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func create(fullPath string) error {
	dirPath := filepath.Dir(fullPath)

	if dirPath == "" {
		dirPath = "./"
	}

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	if strings.HasSuffix(fullPath, "/") {
		return nil
	}

	actionText := "Updated"
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		actionText = "Created"

		file, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()
	}

	fmt.Printf("%s file: %s\n", actionText, fullPath)
	return nil
}

func genFishCompletions(completion string, err error) {
	if err != nil {
		fmt.Printf("Failed to generate fish completion script: %v\n", err)
		return
	}

	completionDir := filepath.Join(os.Getenv("HOME"), ".config", "fish", "completions")
	completionFile := filepath.Join(completionDir, "slap.fish")

	err = os.WriteFile(completionFile, []byte(completion), 0644)
	if err != nil {
		fmt.Printf("Failed to write completion file: %v\n", err)
		return
	}
}
