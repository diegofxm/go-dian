package common

// Party representa una parte (emisor o cliente)
// IMPORTANTE: El orden de los campos debe seguir el schema UBL 2.1 de DIAN
type Party struct {
	PartyIdentification        PartyIdentification `xml:"cac:PartyIdentification"`
	PartyName                  []PartyName         `xml:"cac:PartyName,omitempty"`
	PhysicalLocation           *PhysicalLocation   `xml:"cac:PhysicalLocation,omitempty"`
	PartyTaxScheme             PartyTaxScheme      `xml:"cac:PartyTaxScheme"`
	PartyLegalEntity           PartyLegalEntity    `xml:"cac:PartyLegalEntity"`
	Contact                    *Contact            `xml:"cac:Contact,omitempty"`
	IndustryClassificationCode string              `xml:"cbc:IndustryClassificationCode,omitempty"`
}

// PartyIdentification representa la identificación de una parte
type PartyIdentification struct {
	ID IDType `xml:"cbc:ID"`
}

// PartyName representa el nombre de una parte
type PartyName struct {
	Name string `xml:"cbc:Name"`
}

// PartyTaxScheme representa el esquema tributario
type PartyTaxScheme struct {
	RegistrationName    string           `xml:"cbc:RegistrationName"`
	CompanyID           IDType           `xml:"cbc:CompanyID"`
	TaxLevelCode        TaxLevelCodeType `xml:"cbc:TaxLevelCode"`
	RegistrationAddress *Address         `xml:"cac:RegistrationAddress,omitempty"`
	TaxScheme           TaxScheme        `xml:"cac:TaxScheme"`
}

// TaxLevelCodeType representa el código de nivel tributario con atributos
type TaxLevelCodeType struct {
	Value    string `xml:",chardata"`
	ListName string `xml:"listName,attr,omitempty"`
}

// PartyLegalEntity representa la entidad legal
type PartyLegalEntity struct {
	RegistrationName            string                       `xml:"cbc:RegistrationName"`
	CompanyID                   IDType                       `xml:"cbc:CompanyID"`
	CorporateRegistrationScheme *CorporateRegistrationScheme `xml:"cac:CorporateRegistrationScheme,omitempty"`
}

// CorporateRegistrationScheme representa el esquema de registro corporativo
type CorporateRegistrationScheme struct {
	ID   string `xml:"cbc:ID,omitempty"`
	Name string `xml:"cbc:Name,omitempty"`
}

// PhysicalLocation representa la ubicación física
type PhysicalLocation struct {
	Address Address `xml:"cac:Address"`
}

// Contact representa información de contacto
type Contact struct {
	Name           string `xml:"cbc:Name,omitempty"`
	Telephone      string `xml:"cbc:Telephone,omitempty"`
	Telefax        string `xml:"cbc:Telefax,omitempty"`
	ElectronicMail string `xml:"cbc:ElectronicMail,omitempty"`
	Note           string `xml:"cbc:Note,omitempty"`
}

// Address representa una dirección
type Address struct {
	ID                   string       `xml:"cbc:ID"`
	CityName             string       `xml:"cbc:CityName"`
	PostalZone           string       `xml:"cbc:PostalZone,omitempty"`
	CountrySubentity     string       `xml:"cbc:CountrySubentity"`
	CountrySubentityCode string       `xml:"cbc:CountrySubentityCode"`
	AddressLine          *AddressLine `xml:"cac:AddressLine,omitempty"`
	Country              Country      `xml:"cac:Country"`
}

// AddressLine representa una línea de dirección
type AddressLine struct {
	Line string `xml:"cbc:Line"`
}

// Country representa un país
type Country struct {
	IdentificationCode string `xml:"cbc:IdentificationCode"`
	Name               string `xml:"cbc:Name,omitempty"`
	LanguageID         string `xml:"languageID,attr,omitempty"`
}

// IDType representa un ID con atributos
type IDType struct {
	Value            string `xml:",chardata"`
	SchemeID         string `xml:"schemeID,attr,omitempty"`
	SchemeName       string `xml:"schemeName,attr,omitempty"`
	SchemeAgencyID   string `xml:"schemeAgencyID,attr,omitempty"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr,omitempty"`
}

// TaxScheme representa un esquema de impuesto
type TaxScheme struct {
	ID   string `xml:"cbc:ID"`
	Name string `xml:"cbc:Name"`
}

// AdditionalAccountIDType representa el tipo de persona con atributos
type AdditionalAccountIDType struct {
	Value          string `xml:",chardata"`
	SchemeID       string `xml:"schemeID,attr,omitempty"`
	SchemeAgencyID string `xml:"schemeAgencyID,attr,omitempty"`
	SchemeName     string `xml:"schemeName,attr,omitempty"`
}
