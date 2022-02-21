package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chaosmanagement/chewbyte/pkg/chewbyte"
	"github.com/urfave/cli/v2"
)

// TODO schema validation
// TODO CSV importer/exporter
// TODO ascii-table exporter
// TODO add readme, license, gitignore and so on

const outputPathKey = "output-path"
const outputFormatKey = "output-format"

func main() {
	app := &cli.App{
		Name:                   "chewbyte",
		HelpName:               "chewbyte",
		Usage:                  "Chewbyte is a tool for converting JSON/YAML/Jsonnet documents into YAML/JSON ones. Comes with Jsonnet superpowers. Available as a CLI and Go library.",
		ArgsUsage:              "[path to input file (or files if merged)]",
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    outputPathKey,
				Aliases: []string{"o"},
				Usage:   "path for the resulting file",
			},
			&cli.StringFlag{
				Name:    outputFormatKey,
				Aliases: []string{"f"},
				Value:   "yaml",
				Usage:   "desired format for the resulting file",
				EnvVars: []string{"CHEWBYTE_FORMAT", "FORMAT"},
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				cli.ShowAppHelp(c)
				return nil
			}

			var err error
			var data interface{}

			if c.NArg() == 1 {
				data, err = chewbyte.ImportFile(c.Args().Get(0))
			} else {
				data, err = chewbyte.ImportFiles(c.Args().Slice(), chewbyte.UniqueJoin)
			}

			if err != nil {
				return err
			}

			format := chewbyte.YAML
			if c.IsSet(outputFormatKey) {
				format = chewbyte.DetectFormat("file." + c.String(outputFormatKey))
			} else if c.IsSet(outputPathKey) {
				format = chewbyte.DetectFormat(c.String(outputPathKey))
			}

			if c.IsSet(outputPathKey) {
				err = chewbyte.ExportFile(data, c.String(outputPathKey), format)
				if err != nil {
					return err
				}
			} else {
				str, err := chewbyte.ExportStr(data, format)
				if err != nil {
					return err
				}
				fmt.Print(str)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
