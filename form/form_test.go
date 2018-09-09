package form

import (
	"encoding/json"
	"os"
	"testing"
)

const Form1 string = `C:\Users\Laurent\Golang\src\github.com\lpuig\scopelecformulairecontrol\form\test\PCCDGT.json`

func Test_FormMarshall(t *testing.T) {
	f, err := os.Open(Form1)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var myForm FormModel

	err = json.NewDecoder(f).Decode(&myForm)
	if err != nil {
		t.Fatal("could not decode form", err)
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
