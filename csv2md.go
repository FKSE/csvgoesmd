package csvgoesmd

import (
    "os"
    "strings"
    "io/ioutil"
    "github.com/codegangsta/cli"
    "github.com/qiniu/iconv"
    "fmt"
    "strconv"
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
        cli.StringFlag{
            Name: "encoding",
            Value: "utf-8",
            Usage: "The encoding of the input file. Default is utf-8. The output is always utf-8",
        },
    }
    //main
    app.Action = func(c *cli.Context) {
        if len(c.Args()) < 2 {
            println("Error: Not enough arguments")

            return;
        }
        //try to load file
        bytes, err := ioutil.ReadFile(c.Args()[0])
        //check for error
        if err != nil {
            panic(err);
        }
        csvContent := string(bytes);
        // convert encoding to utf-8
        if c.String("encoding") != "utf-8" {
            cd, err := iconv.Open("utf-8", "latin1",)
            if err != nil {
                panic(err)
            }
            defer cd.Close()
            //convert encoding
            csvContent = cd.ConvString(csvContent);
        }
        //columns
        columns := ParseCsv(csvContent, c.String("delimiter"))

        // build markdown table
        markdown := BuildMarkdown(columns, c.Bool("no-header"))
        //convert string to bytes
        bytes = []byte(markdown)
        //write file
        err = ioutil.WriteFile(c.Args()[1], bytes, 0644)
        //check for error
        if err != nil {
            panic(err);
        }
        os.Exit(1)
    }
    app.Run(os.Args)
}

func ParseCsv(content string, delimiter string) [][]string {
    //split by line
    lines := strings.Split(content, "\n")
    //columns
    var columns [][]string;
    //loop over lines
    for rowIndex, line := range lines {
        //split line in fields
        fields := strings.Split(strings.Trim(line, " \r\n"), delimiter)
        //init columns
        if columns == nil {
            columns = make([][]string, len(fields))
            //init columns
            for i := 0; i < len(fields); i++ {
                columns[i] = make([]string, len(lines))
            }
        }
        //loop over fields
        for columnIndex, field := range fields {
            columns[columnIndex][rowIndex] = field
        }
    }

    return columns
}

func BuildMarkdown(columns [][]string, noHeader bool) string {
    markdown := ""
    header := ""
    for row := 0; row < len(columns[0]); row++ {
        for column := 0; column < len(columns); column++ {
            len := MaxLength(columns[column])
            format := "| %-" + strconv.Itoa(len) + "s "
            markdown += fmt.Sprintf(format, columns[column][row])
            //make header
            if row == 0 {
                header += "|" + strings.Repeat("-", len+2)
            }
        }
        markdown += "|\n"
        //add header separation
        if row == 0 && !noHeader {
            markdown += header + "|\n"
        }
    }
    return markdown
}

func MaxLength(slice []string) int {
    var length int = 0
    for _, elm := range slice {
        if len(elm) > length {
            length = len(elm)
        }
    }
    return length
}