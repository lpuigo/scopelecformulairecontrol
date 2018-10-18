package form2

type FormModel struct {
	Categories      []Category       `json:"Categories"`
	ValidationRules []ValidationRule `json:"ValidationRules"`
}

type Category struct {
	Key           string        `json:"Key"`
	SubCategories []SubCategory `json:"SubCategories"`
	Title         string        `json:"Title"`
}

type ValidationRule struct {
	Rule          string `json:"Rule"`
	ErrorMessage  string `json:"ErrorMessage"`
	RedirectField string `json:"RedirectField"`
}

type SubCategory struct {
	Fields     []Field `json:"Fields"`
	Title      string  `json:"Title"`
	Visibility string  `json:"Visibility"`
}

type Field struct {
	Ref                  string        `json:"Ref"`
	Type                 string        `json:"Type"`
	Label                string        `json:"Label"`
	Visibility           string        `json:"Visibility"`
	Tooltip              string        `json:"Tooltip,omitempty"`
	IsMandatory          bool          `json:"IsMandatory,omitempty"`
	IsReadonly           bool          `json:"IsReadonly,omitempty"`
	DefaultValue         DefaultValue  `json:"DefaultValue,omitempty"`
	Href                 string        `json:"Href,omitempty"`
	CatalogPriceDisplay  bool          `json:"CatalogPriceDisplay,omitempty"`
	CatalogSummary       []string      `json:"CatalogSummary,omitempty"`
	OptionList           []ChoiceGroup `json:"OptionList,omitempty"`
	PidiLineTestOperator []interface{} `json:"PidiLineTestOperator,omitempty"`
	PidiLineTestType     []int64       `json:"PidiLineTestType,omitempty"`
	SectionTitle         string        `json:"SectionTitle,omitempty"`
}

type ChoiceGroup struct {
	Choices    []Choice `json:"Choices"`
	Visibility string   `json:"Visibility"`
}

type Choice struct {
	Display string `json:"Display"`
	Value   string `json:"Value"`
}

type DefaultValue struct {
	OriginatorWorkOrderQuery string      `json:"OriginatorWorkOrderQuery"`
	Static                   interface{} `json:"Static"`
	WorkOrderQuery           string      `json:"WorkOrderQuery,omitempty"` // Requête jsonpath appliqué sur l'OT Planning. Voir le format de l'objet retourné par ApiWorkOrderByIdGetAsync
}
