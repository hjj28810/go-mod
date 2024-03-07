package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type ProcessCmd struct {
	headers          []string
	columns          []string
	data             []map[string]interface{}
	DefaultSheetName string
	Path             string
	Error            error
	style            []func(currentSheet string, f *excelize.File) error
}

func ExcelProcess(val interface{}) (p *ProcessCmd) {
	p = &ProcessCmd{DefaultSheetName: "sheet1"}
	if reflect.TypeOf(val).Kind() != reflect.Slice {
		p.Error = fmt.Errorf("val is not slice")
		return p
	}
	arrBytes, err := json.Marshal(val)
	if err != nil {
		p.Error = err
		return
	}
	var data = make([]map[string]interface{}, 0)
	err = json.Unmarshal(arrBytes, &data)
	if err != nil {
		p.Error = err
		return
	}
	p.data = data
	return p
}

func (p *ProcessCmd) Sheet(sheet string) *ProcessCmd {
	if p.Error != nil {
		return p
	}
	p.DefaultSheetName = sheet
	if sheet == "" {
		p.Error = fmt.Errorf("sheet name is empty")
		return p
	}
	return p
}
func (p *ProcessCmd) Style(f func(currentSheet string, f *excelize.File) error) *ProcessCmd {
	if p.Error != nil {
		return p
	}
	p.style = append(p.style, f)

	return p
}

func (p *ProcessCmd) SavePath(path string) *ProcessCmd {
	if p.Error != nil {
		return p
	}
	if path == "" {
		p.Error = fmt.Errorf("path name is empty")
		return p
	}
	p.Path = path
	return p
}
func (p *ProcessCmd) Headers(headers ...string) *ProcessCmd {
	if p.Error != nil {
		return p
	}
	if len(headers) == 0 {
		p.Error = fmt.Errorf("headers  is empty")
		return p
	}
	p.headers = headers
	return p
}
func (p *ProcessCmd) Columns(columns ...string) *ProcessCmd {
	if p.Error != nil {
		return p
	}
	if len(columns) == 0 {
		p.Error = fmt.Errorf("columns  is empty")
		return p
	}
	p.columns = columns
	return p
}

func (p *ProcessCmd) ToExcel() *ProcessCmd {
	if p.Error != nil {
		return p
	}
	f := excelize.NewFile()
	index, _ := f.NewSheet(p.DefaultSheetName)
	f.SetActiveSheet(index)
	rowNumber := 1

	for i := range p.headers {
		colName, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			p.Error = err
			return p
		}

		err = f.SetCellValue(p.DefaultSheetName, colName+cast.ToString(rowNumber), p.headers[i])
		if err != nil {
			p.Error = err
			return p
		}

	}

	rowNumber++
	for k := range p.data {
		columnNumber := 1
		if len(p.columns) == 0 {

			for j := range p.data[k] {

				colName, err := excelize.ColumnNumberToName(columnNumber)
				if err != nil {
					p.Error = err
					return p
				}

				err = f.SetCellValue(p.DefaultSheetName, colName+cast.ToString(rowNumber), p.data[k][j])
				if err != nil {
					p.Error = err
					return p
				}
				columnNumber++
			}
		} else {
			for j := range p.columns {
				colName, err := excelize.ColumnNumberToName(columnNumber)
				if err != nil {
					p.Error = err
					return p
				}

				err = f.SetCellValue(p.DefaultSheetName, colName+cast.ToString(rowNumber), p.data[k][p.columns[j]])
				if err != nil {
					p.Error = err
					return p
				}
				columnNumber++
			}
		}

		rowNumber++
	}
	for i := range p.style {
		if p.style[i] != nil {
			err := p.style[i](p.DefaultSheetName, f)
			if err != nil {
				p.Error = err
				return p
			}
		}
	}
	if f.Path == "" {
		p.Path = "demo.xlsx"
	}
	err := f.SaveAs(p.Path)
	if err != nil {
		fmt.Println(err)
	}
	return p
}

type ProcessCmdStream struct {
	headers          []string
	columns          []string
	data             []map[string]interface{}
	DefaultSheetName string
	Path             string
	Error            error
	styleFunc        func(styleMap map[string]int, x int, y int, val interface{}) (int, error)
	newStyleMap      map[string]int
	newStyleFunc     []KeyStyle
	colWidth         []*ColWidth
	writer           io.Writer
}

type ColWidth struct {
	mix   int
	max   int
	width float64
}

func ExcelProcessStream(val interface{}) (p *ProcessCmdStream) {
	p = &ProcessCmdStream{DefaultSheetName: "sheet1", newStyleMap: make(map[string]int, 4)}
	if reflect.TypeOf(val).Kind() != reflect.Slice {
		p.Error = fmt.Errorf("val is not slice")
		return p
	}
	arrBytes, err := json.Marshal(val)
	if err != nil {
		p.Error = err
		return
	}
	var data = make([]map[string]interface{}, 0)
	err = json.Unmarshal(arrBytes, &data)
	if err != nil {
		p.Error = err
		return
	}
	if len(data) == 0 {
		p.Error = errors.New("data is empty")
		return p
	}
	p.data = data
	return p
}

func (p *ProcessCmdStream) SavePath(path string) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	if path == "" {
		p.Error = fmt.Errorf("path name is empty")
		return p
	}
	p.Path = path
	return p
}
func (p *ProcessCmdStream) Headers(header ...string) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}

	p.headers = append(p.headers, header...)
	return p
}
func (p *ProcessCmdStream) Columns(columns ...string) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	if len(columns) == 0 {
		p.Error = fmt.Errorf("columns  is empty")
		return p
	}
	p.columns = columns
	return p
}

func (p *ProcessCmdStream) Sheet(sheet string) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	p.DefaultSheetName = sheet
	if sheet == "" {
		p.Error = fmt.Errorf("sheet name is empty")
		return p
	}
	return p
}
func (p *ProcessCmdStream) Style(styleFunc func(styleMap map[string]int, x int, y int, val interface{}) (int, error)) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	p.styleFunc = styleFunc
	return p
}
func (p *ProcessCmdStream) Writer(writer io.Writer) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	p.writer = writer
	return p
}
func (p *ProcessCmdStream) SetColWidth(min, max int, width float64) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	p.colWidth = append(p.colWidth, &ColWidth{mix: min, max: max, width: width})
	return p
}

type KeyStyle struct {
	Key   string
	Style interface{}
}

func (p *ProcessCmdStream) NewStyle(newStyleFunc ...KeyStyle) *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	p.newStyleFunc = newStyleFunc
	return p
}
func (p *ProcessCmdStream) ToExcel() *ProcessCmdStream {
	if p.Error != nil {
		return p
	}
	f := excelize.NewFile()
	if p.newStyleFunc != nil {
		for _, va := range p.newStyleFunc {
			styleId, err := f.NewStyle(va.Style.(*excelize.Style))
			if err != nil {
				p.Error = err
				return p
			}
			p.newStyleMap[va.Key] = styleId
		}
	}
	stream, err := f.NewStreamWriter(p.DefaultSheetName)
	if err != nil {
		p.Error = err
		return p
	}
	for i := range p.colWidth {
		err = stream.SetColWidth(p.colWidth[i].mix, p.colWidth[i].max, p.colWidth[i].width)
		if err != nil {
			p.Error = err
			return p
		}
	}
	y := 0
	if len(p.headers) > 0 {
		headerRaw := make([]interface{}, 0, len(p.headers))
		for i := range p.headers {
			cell := &excelize.Cell{
				Value: p.headers[i],
			}
			if p.styleFunc != nil {
				styleId, err := p.styleFunc(p.newStyleMap, i, 0, p.headers[i])
				if err != nil {
					p.Error = err
					return p
				}
				cell.StyleID = styleId
			}
			headerRaw = append(headerRaw, cell)
		}
		y++
		err = stream.SetRow("A1", headerRaw)
		if err != nil {
			p.Error = err
			return p
		}
	}

	for k := range p.data {
		x := 0
		var rowValues []interface{}
		if len(p.columns) == 0 {
			rowValues = make([]interface{}, 0, len(p.data[k]))
			for j := range p.data[k] {
				cell := &excelize.Cell{
					Value: p.data[k][j],
				}
				if p.styleFunc != nil {
					styleId, err := p.styleFunc(p.newStyleMap, x, y, p.data[k][j])
					if err != nil {
						p.Error = err
						return p
					}
					cell.StyleID = styleId
				}
				rowValues = append(rowValues, cell)
				x++
			}
		} else {
			rowValues = make([]interface{}, 0, len(p.columns))
			for j := range p.columns {
				cell := &excelize.Cell{
					Value: p.data[k][p.columns[j]],
				}
				if p.styleFunc != nil {
					styleId, err := p.styleFunc(p.newStyleMap, x, y, p.data[k][p.columns[j]])
					if err != nil {
						p.Error = err
						return p
					}
					cell.StyleID = styleId
				}
				rowValues = append(rowValues, cell)
				x++
			}
		}
		cellName, err := excelize.CoordinatesToCellName(1, y+1)
		if err != nil {
			p.Error = err
			return p
		}
		p.Error = stream.SetRow(cellName, rowValues)
		if p.Error != nil {
			return p
		}
		y++
	}
	err = stream.Flush()
	if err != nil {
		p.Error = err
		return p
	}

	if p.writer == nil {
		if p.Path == "" {
			p.Path = "demo.xlsx"
		}

		p.Error = f.SaveAs(p.Path)
		return p
	}
	p.Error = f.Write(p.writer)
	return p
}

func ExportExcel(list interface{}, fileName string, headers, columns []string) error {
	return ExcelProcessStream(list).
		Headers(headers...).
		Columns(columns...).
		SavePath(fileName).SetColWidth(1, 11, 32).
		NewStyle(KeyStyle{
			Key:   "style1",
			Style: &excelize.Style{Font: &excelize.Font{Family: "Microsoft YaHei UI", Size: 12}},
		}, KeyStyle{
			Key: "style2",
			Style: &excelize.Style{Font: &excelize.Font{Family: "Microsoft YaHei UI", Size: 12}, Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"#E8E8E8"},
			}},
		}).Style(func(styleMap map[string]int, x int, y int, val interface{}) (int, error) {
		if y%2 == 0 {
			return styleMap["style2"], nil
		} else {
			return styleMap["style1"], nil
		}
	}).ToExcel().Error
}
