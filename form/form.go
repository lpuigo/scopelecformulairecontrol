package form

import (
	csv "encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"os"
	"strings"
)

type FormModel struct {
	Categories []struct {
		Key           string        `json:"Key"`
		Title         string        `json:"Title"`
		SubCategories []SubCategory `json:"SubCategories"`
	} `json:"Categories"`
	ValidationRules []struct {
		Rule          string `json:"Rule"`
		ErrorMessage  string `json:"ErrorMessage"`
		RedirectField string `json:"RedirectField"`
	} `json:"ValidationRules"`
}

type SubCategory struct {
	Key        string  `json:"Key"`
	Title      string  `json:"Title"`
	Visibility string  `json:"Visibility,omitempty"`
	Fields     []Field `json:"Fields"`
}

func (s SubCategory) String() string {
	return fmt.Sprintf("%s [%s]\n", s.Title, s.Key)
}

type Field struct {
	Ref             string           `json:"Ref"`                       // Référence permettant d'identifier un champ
	Type            string           `json:"Type"`                      // Type du champ
	Label           string           `json:"Label"`                     // Libellé associé au champ
	Label2          string           `json:"Label2,omitempty"`          // Libellé associé au champ (utilisé pour certains champs
	Label3          string           `json:"Label3,omitempty"`          // Libellé associé au champ (utilisé pour certains champs
	Visibility      string           `json:"Visibility,omitempty"`      // Règle de visibilité permettant d'afficher/masquer le champ
	Tooltip         string           `json:"Tooltip,omitempty"`         // Affiche un \"?\" en bout de ligne et affiche le contenu dans un tooltip au click du \"?\"\"
	Placeholder     string           `json:"Placeholder,omitempty"`     // Valeur affiché par défaut (valable que pour les champs libres de type text)
	IsMandatory     bool             `json:"IsMandatory,omitempty"`     // Interdit que le champ soit vide (affiche une *)
	IsMandatoryExpr string           `json:"IsMandatoryExpr,omitempty"` // Idem, mais sur la base d'une expression à évaluer
	MessErrRequired string           `json:"MessErrRequired,omitempty"` // Message d'erreur affiché si le champ est vide
	IsReadonly      bool             `json:"IsReadonly"`                // Le champ n'est pas modifiable
	Href            string           `json:"Href,omitempty"`            // Lien hypertext pour les champs de type link ou linkAppli
	Regexp          string           `json:"Regexp,omitempty"`          // Règle de validation des champs texte sous forme d'expression régulière
	MessErrRegexp   string           `json:"MessErrRegexp,omitempty"`   // Message d'erreur affiché si le champ ne respecte pas Regexp
	SectionTitle    string           `json:"SectionTitle,omitempty"`    // Affiche un titre au dessus d'un champ. La règle de visibilité n'affecte pas le titre de la section.
	DefaultValue    DefaultValueType `json:"DefaultValue,omitempty"`    // Description d'une valeur par défaut. Une seule valeur possible (oneOf pas supporté par swagger v2).
	Quality         struct {         // Information de classification qualité du champ (non géré par le moteur de rendu)
		Ref string `json:"Ref"`
	} `json:"Quality,omitempty"`
	PidiLineTestType      []int      `json:"PidiLineTestType,omitempty"`      // Permet de sélectionner les essais mis à disposition
	PidiLineTestOperator  []string   `json:"PidiLineTestOperator,omitempty"`  // Permet de sélectionner les essais mis à disposition
	CatalogServiceFilter  string     `json:"CatalogServiceFilter,omitempty"`  // Concerne le type catalogInput. Valeur du filtre affecté par défaut (code famille de prestation)
	CatalogMaterialFilter string     `json:"CatalogMaterialFilter,omitempty"` // Concerne le type catalogInput. Valeur du filtre affecté par défaut (code famille de matériel)
	CatalogPriceDisplay   bool       `json:"CatalogPriceDisplay,omitempty"`   // Concerne le type catalogInput. Permet d'indiquer si les prix totaux sont affichés dans le composant côté application mobile
	CatalogSummary        []string   `json:"CatalogSummary,omitempty"`        // Concerne le type catalogInput. Permet d'indiquer la référence des champs dont on souhaite faire le résumé
	OptionList            []struct { // Pour les types choix, dresse la liste des choix possibles
		Visibility string `json:"Visibility,omitempty"` // Règle de visibilité permettant d'afficher/masquer le groupe d'option
		Choices    []struct {
			Value   string `json:"Value"`
			Display string `json:"Display"`
		} `json:"Choices"`
	} `json:"OptionList,omitempty"`
	CascadeOptionList []struct { // Pour les types choix dans les listes en cascade, dresse la liste des choix possibles
		Value   string   `json:"Value"`
		Display string   `json:"Display"`
		Items   []string `json:"Items,omitempty"`
	} `json:"CascadeOptionList,omitempty"`
	Fields   []Field `json:"Fields,omitempty"`   // Pour le type itemlist, dresse la liste des champs groupés et répétables
	MaxItem  int     `json:"MaxItem,omitempty"`  // Nombre maximum d'élément dans le groupe de champs itemlist
	MinLevel int     `json:"MinLevel,omitempty"` // niveau technicien minimum inclus nécessaire pour voir le champs
	MaxLevel int     `json:"MaxLevel,omitempty"` // niveau technicien minimum inclus nécessaire pour voir le champs
}

type DefaultValueType struct {
	Static                   interface{} `json:"Static,omitempty"`                   // Valeur en dur dans le format attendu par le type du champ
	OriginatorWorkOrderQuery string      `json:"OriginatorWorkOrderQuery,omitempty"` // Requête xpath 1.0 appliqué sur l'OT brut provenant du canal donneur d'ordre
	WorkOrderQuery           string      `json:"WorkOrderQuery,omitempty"`           // Requête jsonpath appliqué sur l'OT Planning. Voir le format de l'objet retourné par ApiWorkOrderByIdGetAsync
}

func (v DefaultValueType) String() string {
	if v.Static != nil {
		return fmt.Sprintf("Static: %v", v.Static)
	} else if v.OriginatorWorkOrderQuery != "" {
		return "OriginatorWorkOrderQuery: " + v.OriginatorWorkOrderQuery
	} else if v.WorkOrderQuery != "" {
		return "WorkOrderQuery: " + v.WorkOrderQuery
	}
	return ""
}

type FieldType string

const (
	FT_address            FieldType = "address"
	FT_booleanInput       FieldType = "booleanInput"
	FT_cascadeChoice      FieldType = "cascadeChoice"
	FT_dateInput          FieldType = "dateInput"
	FT_dateTimeInput      FieldType = "dateTimeInput"
	FT_gpsInput           FieldType = "gpsInput"
	FT_itemList           FieldType = "itemList"
	FT_link               FieldType = "link"
	FT_linkAppli          FieldType = "linkAppli"
	FT_mailInput          FieldType = "mailInput"
	FT_multipleChoice     FieldType = "multipleChoice"
	FT_numInput           FieldType = "numInput"
	FT_photoInput         FieldType = "photoInput"
	FT_pjInput            FieldType = "pjInput"
	FT_qrcodeInput        FieldType = "qrcodeInput"
	FT_signatureInput     FieldType = "signatureInput"
	FT_singleChoice       FieldType = "singleChoice"
	FT_singleInput        FieldType = "singleInput"
	FT_telInput           FieldType = "telInput"
	FT_text               FieldType = "text"
	FT_timeInput          FieldType = "timeInput"
	FT_pidiLineTestsInput FieldType = "pidiLineTestsInput"
	FT_catalogInput       FieldType = "catalogInput"
	FT_tamTestsInput      FieldType = "tamTestsInput"
)

func NewFormModelFromFile(file string) (*FormModel, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	myForm := &FormModel{}

	err = json.NewDecoder(f).Decode(myForm)
	if err != nil {
		return nil, fmt.Errorf("could not decode form %v", err)
	}
	return myForm, nil
}

func outString(pos, typ, ref, label, mandatory, readonly string) string {
	return fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", pos, typ, ref, label, mandatory, readonly)
}

func mandatory(b bool) string {
	switch b {
	case true:
		return "Obligatoire"
	default:
		return "Facultatif"
	}
}

func readonly(b bool) string {
	switch b {
	case true:
		return "Lect.Seule"
	default:
		return "Lect.Ecrit"
	}
}

func (f FormModel) WriteString(w io.Writer) {
	for ic, c := range f.Categories {
		w.Write([]byte(outString(fmt.Sprintf("%d", ic),
			"Categorie", c.Key, c.Title,
			"",
			"",
		)))
		for isc, s := range c.SubCategories {
			w.Write([]byte(outString(fmt.Sprintf("%d-%d", ic, isc),
				"Sous-Categorie", s.Key, s.Title,
				"",
				"",
			)))
			for ifield, field := range s.Fields {
				w.Write([]byte(outString(fmt.Sprintf("%d-%d-%d", ic, isc, ifield),
					field.Type, field.Ref, field.Label,
					mandatory(field.IsMandatory),
					readonly(field.IsReadonly),
				)))
			}
		}
	}
}

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

func (f FormModel) WriteXLS(w io.Writer) error {
	xs, err := NewXLSSheet()
	if err != nil {
		return err
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

			for ifield, field := range s.Fields {
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
			}
		}
	}

	err = xs.f.Write(w)
	if err != nil {
		return err
	}
	return nil
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
