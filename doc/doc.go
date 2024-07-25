package doc

import (
	"fmt"
	"os"
	"strings"
)

const LevelTitle = 1
const LevelSection = 2
const LevelSubSection = 3

type MarkDownDoc struct {
	builder *strings.Builder
}

func NewMarkDown() *MarkDownDoc {
	m := new(MarkDownDoc)
	m.builder = new(strings.Builder)
	return m
}

func (m *MarkDownDoc) WriteLevel1Title(content string) *MarkDownDoc {
	m.WriteTitle(content, 1)
	return m
}

func (m *MarkDownDoc) write(content string) {
	m.builder.WriteString(content)
}

func (m *MarkDownDoc) GetTitle(content string, level int) string {
	return strings.Repeat("#", level) + " " + content
}

func (m *MarkDownDoc) WriteTitle(content string, level int) *MarkDownDoc {
	m.write(m.GetTitle(content, level))
	m.Writeln()
	return m
}

func (m *MarkDownDoc) WriteWordLine(content string) *MarkDownDoc {
	m.Write(content)
	m.Writeln()
	return m
}

func (m *MarkDownDoc) Write(content string) *MarkDownDoc {
	m.write(content)
	return m
}

func (m *MarkDownDoc) Writeln() *MarkDownDoc {
	m.write("\n")
	return m
}

func (m *MarkDownDoc) WriteLines(lines int) *MarkDownDoc {
	for i := 0; i < lines; i++ {
		m.Writeln()
	}
	return m
}

func (m *MarkDownDoc) WriteJson(content string) *MarkDownDoc {
	m.WriteMultiCode(content, "json")
	return m
}

func (m *MarkDownDoc) GetMultiCode(content, t string) string {
	return fmt.Sprintf("``` %s\n%s\n```\n", t, content)
}

func (m *MarkDownDoc) WriteMultiCode(content, t string) *MarkDownDoc {
	m.write(m.GetMultiCode(content, t))
	return m

}

func (m *MarkDownDoc) WriteCodeLine(content string) *MarkDownDoc {
	m.WriteCode(content)
	m.Writeln()
	return m
}

func (m *MarkDownDoc) GetCode(content string) string {
	return fmt.Sprintf("`%s`", content)
}

func (m *MarkDownDoc) WriteCode(content string) *MarkDownDoc {
	m.write(m.GetCode(content))
	return m
}

func (m *MarkDownDoc) GetTable(t *Table) string {
	return t.String()
}

func (m *MarkDownDoc) WriteTable(t *Table) *MarkDownDoc {
	m.write(m.GetTable(t))
	return m
}

func (m *MarkDownDoc) Export(filename string) error {
	return os.WriteFile(filename, []byte(m.builder.String()), os.ModePerm)
}

func (m *MarkDownDoc) GetLink(desc, url string) string {
	return fmt.Sprintf("[%s](%s)", desc, url)
}

func (m *MarkDownDoc) WriteLink(desc, url string) *MarkDownDoc {
	m.write(m.GetLink(desc, url))
	return m
}

func (m *MarkDownDoc) WriteLinkLine(desc, url string) *MarkDownDoc {
	m.WriteLink(desc, url)
	m.WriteLines(2)
	return m
}

func (m *MarkDownDoc) String() string {
	return m.builder.String()
}

func (m *MarkDownDoc) GetList(content string) string {
	return "- " + content
}

func (m *MarkDownDoc) WriteList(content string) *MarkDownDoc {
	m.write(m.GetList(content))
	m.Writeln()
	return m
}

type Table struct {
	body [][]string
}

func (t *Table) SetTitle(col int, content string) *Table {
	t.body[0][col] = content
	return t
}
func (t *Table) SetContent(row, col int, content string) *Table {
	row = row + 2
	t.body[row][col] = content
	return t
}

func (t *Table) String() string {
	var buffer strings.Builder
	for _, row := range t.body {
		buffer.WriteString("|")
		for _, col := range row {
			buffer.WriteString(col)
			buffer.WriteString("|")
		}
		buffer.WriteString("\n")

	}
	return buffer.String()
}

func NewTable(row, col int) *Table {
	t := new(Table)
	row = row + 2
	t.body = make([][]string, row)
	for i := 0; i < row; i++ {
		t.body[i] = make([]string, col)
		if i == 1 {
			for j := 0; j < col; j++ {
				t.body[i][j] = "----"
			}
		}
	}
	return t
}
