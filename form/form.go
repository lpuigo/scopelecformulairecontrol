package form

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

	BlockFields []Field `json:"BlockFields,omitempty"` // Pour le type block (editeur de formulaire)
	Name        string  `json:"Name,omitempty"`        // Nom du Block pour les Fields de type Block (editeur de formulaire)
	Version     int     `json:"Version,omitempty"`     // Version du Block pour les Fields de type Block (editeur de formulaire)
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
	FT_Block              FieldType = "block"
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
