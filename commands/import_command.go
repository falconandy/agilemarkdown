package commands

import (
	"fmt"
	"github.com/mreider/agilemarkdown/backlog"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
	"strings"
)

var ImportCommand = cli.Command{
	Name:      "import",
	Usage:     "Import an existing Pivotal Tracker story",
	ArgsUsage: "CSV_FILE",
	Action: func(c *cli.Context) error {
		if err := checkIsBacklogDirectory(); err != nil {
			fmt.Println(err)
			return nil
		}
		if c.NArg() == 0 {
			fmt.Println("a csv file should be specified")
			return nil
		}
		for _, csvPath := range c.Args() {
			csvPath, err := filepath.Abs(csvPath)
			if err != nil {
				fmt.Printf("The csv file '%s' is wrong: %v\n", csvPath, err)
				continue
			}
			_, err = os.Stat(csvPath)
			if err != nil {
				fmt.Printf("The csv file '%s' is wrong: %v\n", csvPath, err)
				continue
			}
			ext := strings.ToLower(filepath.Ext(csvPath))
			if ext != ".csv" {
				fmt.Printf("The file '%s' should be a CSV file\n", csvPath)
				continue
			}

			csvImporter := backlog.NewCsvImporter(csvPath, ".")
			err = csvImporter.Import()
			if err != nil {
				fmt.Printf("Import of the csv file '%s' failed: %v\n", csvPath, err)
				continue
			}
		}
		return nil
	},
}
