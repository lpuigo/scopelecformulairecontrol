package form

import (
	"encoding/csv"
	"fmt"
	"io"
)

func outCSV(pos, typ, ref, label, mandatory, readonly, section string) []string {
	return []string{pos, typ, ref, label, mandatory, readonly, section}
}

func (f FormModel) WriteCSV(w io.Writer) {
	csvw := csv.NewWriter(w)
	csvw.Comma = ';'
	defer csvw.Flush()

	for ic, c := range f.Categories {
		csvw.Write(outCSV(fmt.Sprintf("%d", ic),
			"Categorie", c.Key, c.Title,
			"",
			"",
			"",
		))
		for isc, s := range c.SubCategories {
			csvw.Write(outCSV(fmt.Sprintf("%d-%d", ic, isc),
				"Sous-Categorie", s.Key, s.Title,
				"",
				"",
				"",
			))
			for ifield, field := range s.Fields {
				csvw.Write(outCSV(fmt.Sprintf("%d-%d-%d", ic, isc, ifield),
					field.Type, field.Ref, field.Label,
					mandatory(field.IsMandatory),
					readonly(field.IsReadonly),
					field.SectionTitle,
				))
			}
		}
	}
}
