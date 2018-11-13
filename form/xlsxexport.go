package form

import (
	"fmt"
	"github.com/lpuig/scopelecformulairecontrol/blocklist"
	"github.com/tealeg/xlsx"
	"io"
	"strconv"
	"strings"
)

func (f FormModel) WriteXLS(w io.Writer) (bl blocklist.BlockDescList, err error) {
	bl = blocklist.MakeBlockDescList()
	xs, err := NewXLSSheet()
	if err != nil {
		return
	}
	xs.xlsHeader()

	for ic, c := range f.Categories {
		xs.xlsRow(fmt.Sprintf("%d", ic),
			"Categorie", c.Key, c.Title,
			"",
			"",
			"",
			"",
			"",
		)
		xs.rowColor("00000000", "0000b6ff")
		for isc, s := range c.SubCategories {
			xs.xlsRow(fmt.Sprintf("%d-%d", ic, isc),
				"Sous-Categorie", s.Key, s.Title,
				"",
				"",
				"",
				s.Visibility,
				"",
			)
			xs.rowColor("00000000", "0051cdff")

			writeFields(w, xs, s.Fields, ic, isc, 0, bl)
		}
	}

	err = xs.f.Write(w)
	if err != nil {
		return
	}
	return
}

func writeFields(w io.Writer, xs *XLSSheet, fields []Field, ic, isc, ifield int, bl blocklist.BlockDescList) int {

	for _, field := range fields {
		switch field.Type {
		case string(FT_Block):
			xs.xlsRow("",
				field.Type, field.Name, "v"+strconv.Itoa(field.Version),
				"",
				"",
				"",
				"",
				"",
			)
			xs.rowColor("00000000", "00b2ffb2")
			bl.Add(field.Name, field.Version)

			ifield = writeFields(w, xs, field.BlockFields, ic, isc, ifield, bl)

		default:
			xs.xlsRow(fmt.Sprintf("%d-%d-%d", ic, isc, ifield),
				field.Type, field.Ref, field.Label,
				mandatory(field.IsMandatory),
				readonly(field.IsReadonly),
				field.SectionTitle,
				field.Visibility,
				field.DefaultValue.String(),
			)
			if strings.HasPrefix(field.Ref, "PIDI_") {
				xs.r.Cells[2].GetStyle().Font.Bold = true
			}
			if field.IsMandatory {
				xs.cellColor(4, "00a80000")
			}
			if field.IsReadonly {
				xs.cellColor(5, "00a80000")
			}
			ifield += 1
		}
	}
	return ifield
}

type XLSSheet struct {
	f *xlsx.File
	s *xlsx.Sheet
	r *xlsx.Row
}

func NewXLSSheet() (xs *XLSSheet, err error) {
	xs = &XLSSheet{
		f: xlsx.NewFile(),
	}
	xlsx.SetDefaultFont(11, "Calibri")
	xs.s, err = xs.f.AddSheet("Form")
	if err != nil {
		return
	}
	return
}

func (xs *XLSSheet) xlsHeader() {
	xs.xlsRow(
		"Position",
		"Type",
		"Ref",
		"Label",
		"IsMandatory",
		"IsReadonly",
		"SectionTitle",
		"VisibilityRule",
		"DefaultValue",
	)

	styles := []struct {
		width    float64
		autowrap bool
	}{
		{8, false},
		{17, false},
		{43, false},
		{50, true},
		{12, false},
		{11, false},
		{30, true},
		{50, true},
		{50, true},
	}

	for i, c := range xs.s.Cols {
		c.Width = styles[i].width
		if styles[i].autowrap {
			cs := c.GetStyle()
			cs.Alignment = xlsx.Alignment{WrapText: true}
			cs.ApplyAlignment = true
		}
	}
}

func (xs *XLSSheet) xlsRow(pos, typ, ref, label, mandatory, readonly, section, visibility, defaultvalue string) {
	xs.r = xs.s.AddRow()
	for _, str := range []string{pos, typ, ref, label, mandatory, readonly, section, visibility, defaultvalue} {
		c := xs.r.AddCell()
		c.SetString(str)
	}
}

func (xs *XLSSheet) rowColor(fgcolor, bgcolor string) {
	for _, c := range xs.r.Cells {
		s := c.GetStyle()
		s.Fill = *xlsx.NewFill("solid", bgcolor, "00000000")
		s.ApplyFill = true
		s.Font.Color = fgcolor
	}
}

func (xs *XLSSheet) cellColor(num int, fgcolor string) {
	xs.r.Cells[num].GetStyle().Font.Color = fgcolor
}
