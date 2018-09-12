package main

import (
	"flag"
	"github.com/lpuig/scopelecformulairecontrol/form"
	"log"
	"os"
)

//go:generate go build -o ../formctrl.exe

var filename string
var output string
var outf *os.File

func main() {
	flag.StringVar(&filename, "f", "", "json model file to process")
	flag.StringVar(&output, "o", "", "output file")
	flag.Parse()

	model, err := form.NewFormModelFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if output == "" {
		outf = os.Stdout
	} else {
		outf, err = os.Create(output)
		if err != nil {
			log.Fatal("could not create output file", err)
		}
		defer outf.Close()
	}

	model.OutString(outf)
}
