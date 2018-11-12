package main

import (
	"fmt"
	"github.com/lpuig/scopelecformulairecontrol/form"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type GuiContext struct {
	textEdit *walk.TextEdit
}

func main() {
	gc := GuiContext{}
	_, err := MainWindow{
		Title:   "SPRINT Formulaire Controler",
		MinSize: Size{320, 240},
		Layout:  VBox{},
		OnDropFiles: func(files []string) {
			gc.GenXLSX(files)
		},
		Children: []Widget{
			Label{Text: "Glisser un formulaire JSON ici ..."},
			TextEdit{
				AssignTo: &gc.textEdit,
				ReadOnly: true,
				//Text:      "Glisser un formulaire JSON ici ...",
				//Alignment: AlignCenter,
			},
		},
	}.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func outputfile(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + ".xlsx"
}

func (gc GuiContext) GenXLSX(filenames []string) {
	for _, filename := range filenames {
		output := outputfile(filename)
		err := genXlxsfromForm(filename, output)
		if err != nil {
			newText := gc.textEdit.Text() + "\r\nErreur : " + err.Error()
			gc.textEdit.SetText(newText)
			continue
		}
		newText := gc.textEdit.Text() + "\r\n" + "Résultat : " + output
		gc.textEdit.SetText(newText)
	}
}

func genXlxsfromForm(jsonfile, xlsxfile string) error {
	model, err := form.NewFormModelFromFile(jsonfile)
	if err != nil {
		return fmt.Errorf("impossible de lire le Formulaire : %v", err)
	}

	outf, err := os.Create(xlsxfile)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier résultat : %v", err)
	}
	err = model.WriteXLS(outf)
	if err != nil {
		os.Remove(xlsxfile)
		return fmt.Errorf("impossible d'écrire le fichier XLSX : %v", err)
	}
	outf.Close()
	return nil
}
