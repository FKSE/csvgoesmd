package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/qiniu/iconv"
)

func main() {
	app := cli.NewApp()
	//define meta
	app.Name = "csv2md"
	app.Usage = "Convert CSV files to Markdown files"
	app.Author = "Fridolin Koch"
	app.Version = "1.0.0"
	app.Email = "hi@fkse.io"
	//define flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "delimiter, d",
			Value: ",",
			Usage: "The field delimiter. Default value is ,",
		},
		cli.BoolFlag{
			Name:  "no-header",
			Usage: "Don't use first line as header",
		},
		cli.StringFlag{
			Name:  "encoding",
			Value: "utf-8",
			Usage: "The encoding of the input file. Default is utf-8. The output is always utf-8",
		},
	}
	//main
	app.Action = run
	app.Run(os.Args)
}

func run(c *cli.Context) {
	if len(c.Args()) < 2 {
		println("Error: Not enough arguments")
		return
	}
	//try to load file
	csvfile, err := os.Open(c.Args().Get(0))
	if err != nil {
		panic(err)
	}
	var input io.Reader = csvfile
	// convert encoding to utf-8
	if c.String("encoding") != "utf-8" {
		cd, err := iconv.Open("utf-8", "latin1")
		if err != nil {
			panic(err)
		}
		defer cd.Close()
		input = iconv.NewReader(cd, csvfile, 0)
	}
	// read markdown
	reader := csv.NewReader(input)
	reader.Comma = rune(c.String("delimiter")[0])
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	//dertermine output
	var output io.Writer
	if c.Args().Get(1) == "-" {
		output = os.Stdout
	} else {
		output, err = os.Create(c.Args().Get(1))
		if err != nil {
			panic(err)
		}
	}
	// build markdown table
	buildMarkdown(output, records, c.Bool("no-header"))
}

func buildMarkdown(output io.Writer, records [][]string, noHeader bool) {
	// get col lengs
	clen, err := maxColLen(records)
	if err != nil {
		panic(err)
	}
	// Build md table
	for r := 0; r < len(records); r++ {
		for c := 0; c < len(records[r]); c++ {
			format := "| %-" + strconv.Itoa(clen[c]) + "s "
			fmt.Fprintf(output, format, records[r][c])
		}
		fmt.Fprintln(output, "|")
		//make header
		if r == 0 && !noHeader {
			for _, l := range clen {
				fmt.Fprintf(output, "|%s", strings.Repeat("-", l+2))
			}
			fmt.Fprintln(output, "|")
		}
	}
}

func maxColLen(records [][]string) ([]int, error) {
	if len(records) < 1 {
		return nil, errors.New("Invalid array given.")
	}
	lens := make([]int, len(records[0]))
	// Loop over rows
	for r := 0; r < len(records); r++ {
		for c := 0; c < len(records[r]); c++ {
			if l := len(records[r][c]); lens[c] < l {
				lens[c] = l
			}
		}
	}
	return lens, nil
}
