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
		Name:            "slap",
		Usage:           "Modern replacement for touch and mkdir commands",
		HideHelpCommand: true,
		UsageText:       usageText,
		Description:     description,

		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return fmt.Errorf("expected 1 path argument, but got %d", c.NArg())
			}

			pathArg := c.Args().First()
			if pathArg == "" {
				return fmt.Errorf("path argument is empty")
			}

			return create(pathArg)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func create(fullPath string) error {
	dirPath := filepath.Dir(fullPath)
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
