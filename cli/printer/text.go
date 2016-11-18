package printer

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/ryanuber/columnize"
	"gitlab.imshealth.com/xfra/layer0/cli/entity"
	"gitlab.imshealth.com/xfra/layer0/cli/printer/table"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"os"
	"time"
)

type TextPrinter struct {
	spinner *spinner.Spinner
}

func (this *TextPrinter) StartSpinner(prefix string) {
	if this.spinner != nil {
		this.spinner.Stop()
	}

	this.spinner = spinner.New(spinner.CharSets[26], 1*time.Second)
	this.spinner.Prefix = prefix
	this.spinner.Start()
}

func (this *TextPrinter) StopSpinner() {
	if this.spinner != nil {
		this.spinner.Stop()
		fmt.Println()
	}
}

func (this *TextPrinter) PrintEntity(e entity.Entity) error {
	return this.PrintEntities([]entity.Entity{e})
}

func (this *TextPrinter) PrintEntities(entities []entity.Entity) error {
	tables := []table.Table{}
	for _, e := range entities {
		tables = append(tables, e.Table())
	}

	return this.printTables(tables)
}

func (this *TextPrinter) printTables(tables []table.Table) error {
	if len(tables) == 0 {
		fmt.Println("You don't have any entities of this type")
		return nil
	}

	// combine tables that have the same headers
	uniqueTables := map[string][]table.Table{}
	for _, table := range tables {
		var key string
		for _, column := range table {
			key += column.Title
		}

		uniqueTables[key] = append(uniqueTables[key], table)
	}

	for _, tables := range uniqueTables {
		formatted := this.formatTables(tables)
		this.Printf("%s\n", formatted)
	}

	return nil
}

// combineTables assumes every table in the list uses the same header
// it will combine the rows in each table and format them into a single, printable table
func (this *TextPrinter) formatTables(tables []table.Table) string {
	var header string
	for _, tableColumns := range tables {
		for _, column := range tableColumns {
			header += fmt.Sprintf("%s |", column.Title)
		}

		break
	}

	rows := []string{header}
	for _, table := range tables {
		rows = append(rows, this.formatRows(table)...)
	}

	return columnize.SimpleFormat(rows)
}

// formatRows will populate each cell in the table with a string
// it returns a list of rows in columnized format
func (this *TextPrinter) formatRows(table table.Table) []string {
	var maxColumnRows int
	for _, column := range table {
		if columnRows := len(column.Rows); columnRows > maxColumnRows {
			maxColumnRows = columnRows
		}
	}

	rows := []string{}
	for i := 0; i < maxColumnRows; i++ {
		var row string

		for _, column := range table {
			if i < len(column.Rows) {
				row += fmt.Sprintf("%s |", column.Rows[i])
			} else {
				row += " |"
			}
		}

		rows = append(rows, row)
	}

	return rows
}

func (this *TextPrinter) PrintLogs(logFiles []*models.LogFile) error {
	for _, logFile := range logFiles {
		fmt.Println(logFile.Name)
		for i := 0; i < len(logFile.Name); i++ {
			fmt.Printf("-")
		}

		fmt.Println()
		for _, line := range logFile.Lines {
			fmt.Println(line)
		}
		fmt.Println()
	}

	return nil
}

func (this *TextPrinter) Printf(format string, tokens ...interface{}) {
	this.StopSpinner()
	fmt.Printf(format, tokens...)
}

func (this *TextPrinter) Fatalf(code int64, format string, tokens ...interface{}) {
	this.Printf(format, tokens...)
	fmt.Println()
	os.Exit(1)
}
