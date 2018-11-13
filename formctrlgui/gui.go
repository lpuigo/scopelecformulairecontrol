package main

import (
	"fmt"
	"github.com/lpuig/scopelecformulairecontrol/blocklist"
	"github.com/lpuig/scopelecformulairecontrol/form"
	"github.com/lxn/walk"
	"log"
	"os"
	"path/filepath"
	"strings"

	. "github.com/lxn/walk/declarative"
)

type GuiContext struct {
	textEdit  *walk.TextEdit
	blockEdit *walk.TextEdit
}

func main() {
	gc := GuiContext{}
	_, err := MainWindow{
		Title:   "SPRINT Formulaire Controler",
		MinSize: Size{640, 480},
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
				VScroll: true,
			},
			TextEdit{
				AssignTo:  &gc.blockEdit,
				ReadOnly:  true,
				VScroll:   true,
				MaxLength: 100 * 1024,
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
		files, err := walkPath(filename)
		if err != nil {
			gc.textEdit.AppendText("Error:" + err.Error() + "\r\n ")
		}
		for _, file := range files {
			resFilename := outputfile(file)
			bl, err := genXlxsfromForm(file, resFilename)
			if err != nil {
				newText := gc.textEdit.Text() + "\r\nErreur : " + err.Error()
				gc.textEdit.SetText(newText)
				continue
			}
			newText := "Résultat : " + resFilename + "\r\n"
			gc.textEdit.AppendText(newText)

			gc.blockEdit.AppendText(strings.Join(bl.Strings(file+"\t"), "\r\n") + "\r\n")
		}
	}
}

func genXlxsfromForm(jsonfile, xlsxfile string) (bl blocklist.BlockDescList, err error) {
	model, err := form.NewFormModelFromFile(jsonfile)
	if err != nil {
		err = fmt.Errorf("impossible de lire le Formulaire : %v", err)
		return
	}
	outf, err := os.Create(xlsxfile)
	if err != nil {
		err = fmt.Errorf("impossible de créer le fichier résultat : %v", err)
		return
	}
	bl, err = model.WriteXLS(outf)
	if err != nil {
		os.Remove(xlsxfile)
		err = fmt.Errorf("impossible d'écrire le fichier XLSX : %v", err)
		return
	}
	outf.Close()
	return
}

func walkPath(path string) ([]string, error) {
	res := []string{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		res = append(res, path)
		return nil
	})
	return res, err
}
