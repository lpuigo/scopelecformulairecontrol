package main

import (
	"flag"
	"github.com/lpuig/scopelecformulairecontrol/form"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go build -o ../formctrl.exe
//go:generate bash make32.sh

var filenames []string
var output string
var outf *os.File

func outputfile(filename string) string {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	newExt := ".xlsx"
	if ext == "" {
		base = base + newExt
	} else {
		base = strings.Replace(base, ext, newExt, 1)
	}
	return filepath.Join(dir, base)
}

func main() {
	flag.Parse()
	filenames = flag.Args()

	for _, filename := range filenames {
		model, err := form.NewFormModelFromFile(filename)
		if err != nil {
			log.Print("could not parse Model Form :", err)
			continue
		}

		output := outputfile(filename)

		outf, err = os.Create(output)
		if err != nil {
			log.Print("could not create output file", err)
			continue
		}
		err = model.WriteXLS(outf)
		if err != nil {
			log.Print("could not generate XLS", err)
			os.Remove(output)
		} else {
			log.Print("produced ", output)
			outf.Close()
		}
	}
}
