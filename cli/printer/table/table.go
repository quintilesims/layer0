package table

type Column struct {
	Title      string
	Rows       []string
	IsSingular bool
}

func NewSingleRowColumn(title, row string) Column {
	return Column{
		IsSingular: true,
		Title:      title,
		Rows:       []string{row},
	}
}

func NewMultiRowColumn(title string, rows []string) Column {
	return Column{
		IsSingular: false,
		Title:      title,
		Rows:       rows,
	}
}

type Table []Column
