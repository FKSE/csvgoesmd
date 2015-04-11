package main

import (
    "os"
    "github.com/codegangsta/cli"
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
    app.Flags = []cli.Flag {
        cli.StringFlag{
            Name: "delimiter, d",
            Value: ",",
            Usage: "The field delimiter. Default value is ,",
        },
        cli.StringFlag{
            Name: "enclosure, c",
            Value: "\"",
            Usage: "The enclosure of the fields. Default value is \"",
        },
        cli.StringFlag{
            Name: "escape, e",
            Value: "\\",
            Usage: "The escape character. Default value is \\",
        },
        cli.IntFlag{
            Name: "padding, p",
            Value: 1,
            Usage: "Ammount of whitespaces to use for the cell padding, 1 by default",
        },
        cli.BoolFlag{
            Name: "no-header",
            Usage: "Don't use first line as header",
        },
    }
    //main
    app.Action = func(c *cli.Context) {
        if (len(c.Args()) < 2) {

        }

        os.Exit(0)
    }
    app.Run(os.Args)
}
