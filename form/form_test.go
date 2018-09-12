package form

import (
	"testing"
)

const Form1 string = `C:\Users\Laurent\Golang\src\github.com\lpuig\scopelecformulairecontrol\form\test\PCCDGT.json`

func Test_FormMarshall(t *testing.T) {
	myForm, err := NewFormModelFromFile(Form1)
	if err != nil {
		t.Fatal(err)
	}

	for ic, c := range myForm.Categories {
		t.Logf("%d: %s [%s]\n", ic, c.Title, c.Key)
		for isc, s := range c.SubCategories {
			t.Logf("\t%d-%d: %s", ic, isc, s.String())
			for ifield, field := range s.Fields {
				t.Logf("\t\t%d-%d-%d: %s [%s]", ic, isc, ifield, field.Label, field.Ref)
			}
		}
	}

}
