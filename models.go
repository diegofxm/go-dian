package dian

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Invoice representa una factura electrónica UBL 2.1
type Invoice struct {
	XMLName  xml.Name `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`
	XmlnsCbc string   `xml:"xmlns:cbc,attr"`
	XmlnsExt string   `xml:"xmlns:ext,attr"`
	XmlnsCac string   `xml:"xmlns:cac,attr"`
	XmlnsSts string   `xml:"xmlns:sts,attr"`

	UBLExtensions *UBLExtensions `xml:"ext:UBLExtensions,omitempty"`

	UBLVersionID         string               `xml:"cbc:UBLVersionID"`
	CustomizationID      string               `xml:"cbc:CustomizationID"`
	ProfileID            string               `xml:"cbc:ProfileID"`
	ProfileExecutionID   string               `xml:"cbc:ProfileExecutionID"`
	ID                   string               `xml:"cbc:ID"`
	UUID                 UUIDType             `xml:"cbc:UUID"`
	IssueDate            string               `xml:"cbc:IssueDate"`
	IssueTime            string               `xml:"cbc:IssueTime"`
	DueDate              string               `xml:"cbc:DueDate,omitempty"`
	InvoiceTypeCode      string               `xml:"cbc:InvoiceTypeCode"`
	Note                 []string             `xml:"cbc:Note,omitempty"`
	DocumentCurrencyCode DocumentCurrencyType `xml:"cbc:DocumentCurrencyCode"`
	LineCountNumeric     int                  `xml:"cbc:LineCountNumeric"`

	InvoicePeriod           *InvoicePeriod          `xml:"cac:InvoicePeriod,omitempty"`
	BillingReference        []BillingReference      `xml:"cac:BillingReference,omitempty"`
	AccountingSupplierParty AccountingSupplierParty `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty AccountingCustomerParty `xml:"cac:AccountingCustomerParty"`
	TaxRepresentativeParty  *TaxRepresentativeParty `xml:"cac:TaxRepresentativeParty,omitempty"`
	Delivery                *Delivery               `xml:"cac:Delivery,omitempty"`
	DeliveryTerms           *DeliveryTerms          `xml:"cac:DeliveryTerms,omitempty"`
	PaymentMeans            []PaymentMeans          `xml:"cac:PaymentMeans,omitempty"`
	PaymentTerms            []PaymentTerms          `xml:"cac:PaymentTerms,omitempty"`
	PrepaidPayment          []PrepaidPayment        `xml:"cac:PrepaidPayment,omitempty"`
	TaxTotal                []TaxTotal              `xml:"cac:TaxTotal"`
	LegalMonetaryTotal      LegalMonetaryTotal      `xml:"cac:LegalMonetaryTotal"`
	InvoiceLines            []InvoiceLine           `xml:"cac:InvoiceLine"`
}

// Validate valida los campos requeridos de la factura
func (i *Invoice) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("ID de factura es requerido")
	}
	if i.IssueDate == "" {
		return fmt.Errorf("fecha de emisión es requerida")
	}
	if i.AccountingSupplierParty.Party.PartyTaxScheme.CompanyID.Value == "" {
		return fmt.Errorf("NIT del emisor es requerido")
	}
	if i.AccountingCustomerParty.Party.PartyTaxScheme.CompanyID.Value == "" {
		return fmt.Errorf("NIT del cliente es requerido")
	}
	if len(i.InvoiceLines) == 0 {
		return fmt.Errorf("debe haber al menos una línea de factura")
	}
	return nil
}

// UUIDType representa el UUID con atributos
type UUIDType struct {
	Value      string `xml:",chardata"`
	SchemeID   string `xml:"schemeID,attr"`
	SchemeName string `xml:"schemeName,attr"`
}

// DocumentCurrencyType representa la moneda del documento con atributos
type DocumentCurrencyType struct {
	Value          string `xml:",chardata"`
	ListAgencyID   string `xml:"listAgencyID,attr,omitempty"`
	ListAgencyName string `xml:"listAgencyName,attr,omitempty"`
	ListID         string `xml:"listID,attr,omitempty"`
}

// InvoicePeriod representa el período de facturación
type InvoicePeriod struct {
	StartDate string `xml:"cbc:StartDate"`
	EndDate   string `xml:"cbc:EndDate"`
}

// BillingReference representa referencias a documentos previos
type BillingReference struct {
	InvoiceDocumentReference InvoiceDocumentReference `xml:"cac:InvoiceDocumentReference"`
}

// InvoiceDocumentReference representa una referencia a factura
type InvoiceDocumentReference struct {
	ID                  string   `xml:"cbc:ID"`
	UUID                UUIDType `xml:"cbc:UUID,omitempty"`
	IssueDate           string   `xml:"cbc:IssueDate,omitempty"`
	DocumentDescription string   `xml:"cbc:DocumentDescription,omitempty"`
}

// AccountingSupplierParty representa el emisor de la factura
type AccountingSupplierParty struct {
	AdditionalAccountID AdditionalAccountIDType `xml:"cbc:AdditionalAccountID"`
	Party               Party                   `xml:"cac:Party"`
}

// AccountingCustomerParty representa el cliente
type AccountingCustomerParty struct {
	AdditionalAccountID AdditionalAccountIDType `xml:"cbc:AdditionalAccountID"`
	Party               Party                   `xml:"cac:Party"`
}

// AdditionalAccountIDType representa el tipo de persona con atributos
type AdditionalAccountIDType struct {
	Value      string `xml:",chardata"`
	SchemeName string `xml:"schemeName,attr,omitempty"`
}

// TaxRepresentativeParty representa el representante fiscal
type TaxRepresentativeParty struct {
	PartyIdentification PartyIdentification `xml:"cac:PartyIdentification"`
}

// Party representa una parte (emisor o cliente)
type Party struct {
	PartyName                  []PartyName         `xml:"cac:PartyName,omitempty"`
	PhysicalLocation           *PhysicalLocation   `xml:"cac:PhysicalLocation,omitempty"`
	PartyTaxScheme             PartyTaxScheme      `xml:"cac:PartyTaxScheme"`
	PartyLegalEntity           PartyLegalEntity    `xml:"cac:PartyLegalEntity"`
	Contact                    *Contact            `xml:"cac:Contact,omitempty"`
	IndustryClassificationCode string              `xml:"cbc:IndustryClassificationCode,omitempty"`
	PartyIdentification        PartyIdentification `xml:"cac:PartyIdentification"`
}

// PhysicalLocation representa la ubicación física
type PhysicalLocation struct {
	Address Address `xml:"cac:Address"`
}

// Delivery representa información de entrega
type Delivery struct {
	DeliveryAddress *Address       `xml:"cac:DeliveryAddress,omitempty"`
	DeliveryParty   *DeliveryParty `xml:"cac:DeliveryParty,omitempty"`
}

// DeliveryParty representa la parte que realiza la entrega
type DeliveryParty struct {
	PartyName        []PartyName       `xml:"cac:PartyName,omitempty"`
	PhysicalLocation *PhysicalLocation `xml:"cac:PhysicalLocation,omitempty"`
	PartyTaxScheme   PartyTaxScheme    `xml:"cac:PartyTaxScheme"`
	PartyLegalEntity PartyLegalEntity  `xml:"cac:PartyLegalEntity"`
	Contact          *Contact          `xml:"cac:Contact,omitempty"`
}

// DeliveryTerms representa los términos de entrega
type DeliveryTerms struct {
	ID                     string            `xml:"cbc:ID,omitempty"`
	SpecialTerms           string            `xml:"cbc:SpecialTerms,omitempty"`
	LossRiskResponsibility string            `xml:"cbc:LossRiskResponsibilityCode,omitempty"`
	DeliveryLocation       *DeliveryLocation `xml:"cac:DeliveryLocation,omitempty"`
}

// PaymentMeans representa el medio de pago
type PaymentMeans struct {
	ID               string `xml:"cbc:ID"`
	PaymentMeansCode string `xml:"cbc:PaymentMeansCode"`
	PaymentDueDate   string `xml:"cbc:PaymentDueDate,omitempty"`
}

// PaymentTerms representa las condiciones de pago
type PaymentTerms struct {
	PaymentMeansID string `xml:"cbc:PaymentMeansID,omitempty"`
	PaymentDueDate string `xml:"cbc:PaymentDueDate,omitempty"`
	Note           string `xml:"cbc:Note,omitempty"`
}

// PrepaidPayment representa un pago anticipado
type PrepaidPayment struct {
	ID            string     `xml:"cbc:ID"`
	PaidAmount    AmountType `xml:"cbc:PaidAmount"`
	ReceivedDate  string     `xml:"cbc:ReceivedDate,omitempty"`
	PaidDate      string     `xml:"cbc:PaidDate,omitempty"`
	InstructionID string     `xml:"cbc:InstructionID,omitempty"`
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

// AmountType representa un monto con moneda
type AmountType struct {
	Value      float64 `xml:"-"`
	CurrencyID string  `xml:"currencyID,attr"`
}

// MarshalXML implementa xml.Marshaler para evitar notación científica
func (a AmountType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias AmountType
	aux := struct {
		*Alias
		Value string `xml:",chardata"`
	}{
		Alias: (*Alias)(&a),
		Value: fmt.Sprintf("%.2f", a.Value),
	}
	return e.EncodeElement(aux, start)
}

// Quantity representa una cantidad con unidad de medida
type Quantity struct {
	Value    float64 `xml:"-"`
	UnitCode string  `xml:"unitCode,attr"`
}

// MarshalXML implementa xml.Marshaler para evitar notación científica
func (q Quantity) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Alias Quantity
	aux := struct {
		*Alias
		Value string `xml:",chardata"`
	}{
		Alias: (*Alias)(&q),
		Value: fmt.Sprintf("%.4f", q.Value),
	}
	return e.EncodeElement(aux, start)
}

// TaxTotal representa el total de impuestos
type TaxTotal struct {
	TaxAmount   AmountType    `xml:"cbc:TaxAmount"`
	TaxSubtotal []TaxSubtotal `xml:"cac:TaxSubtotal"`
}

// TaxSubtotal representa un subtotal de impuesto
type TaxSubtotal struct {
	TaxableAmount AmountType  `xml:"cbc:TaxableAmount"`
	TaxAmount     AmountType  `xml:"cbc:TaxAmount"`
	TaxCategory   TaxCategory `xml:"cac:TaxCategory"`
}

// TaxCategory representa una categoría de impuesto
type TaxCategory struct {
	Percent   float64   `xml:"cbc:Percent"`
	TaxScheme TaxScheme `xml:"cac:TaxScheme"`
}

// TaxScheme representa un esquema de impuesto
type TaxScheme struct {
	ID   string `xml:"cbc:ID"`
	Name string `xml:"cbc:Name"`
}

// LegalMonetaryTotal representa el total monetario legal
type LegalMonetaryTotal struct {
	LineExtensionAmount AmountType  `xml:"cbc:LineExtensionAmount"`
	TaxExclusiveAmount  AmountType  `xml:"cbc:TaxExclusiveAmount"`
	TaxInclusiveAmount  AmountType  `xml:"cbc:TaxInclusiveAmount"`
	PrepaidAmount       *AmountType `xml:"cbc:PrepaidAmount,omitempty"`
	PayableAmount       AmountType  `xml:"cbc:PayableAmount"`
}

// InvoiceLine representa una línea de factura
type InvoiceLine struct {
	ID                    string               `xml:"cbc:ID"`
	InvoicedQuantity      Quantity             `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount   AmountType           `xml:"cbc:LineExtensionAmount"`
	FreeOfChargeIndicator *bool                `xml:"cbc:FreeOfChargeIndicator,omitempty"`
	Delivery              *InvoiceLineDelivery `xml:"cac:Delivery,omitempty"`
	AllowanceCharge       []AllowanceCharge    `xml:"cac:AllowanceCharge,omitempty"`
	TaxTotal              []TaxTotal           `xml:"cac:TaxTotal,omitempty"`
	DocumentReference     []DocumentReference  `xml:"cac:DocumentReference,omitempty"`
	PricingReference      *PricingReference    `xml:"cac:PricingReference,omitempty"`
	Item                  Item                 `xml:"cac:Item"`
	Price                 Price                `xml:"cac:Price"`
}

// InvoiceLineDelivery representa la entrega de una línea
type InvoiceLineDelivery struct {
	DeliveryLocation *DeliveryLocation `xml:"cac:DeliveryLocation,omitempty"`
}

// DeliveryLocation representa la ubicación de entrega
type DeliveryLocation struct {
	ID IDType `xml:"cbc:ID"`
}

// AllowanceCharge representa un descuento o cargo
type AllowanceCharge struct {
	ID                      string     `xml:"cbc:ID"`
	ChargeIndicator         bool       `xml:"cbc:ChargeIndicator"`
	AllowanceChargeReason   string     `xml:"cbc:AllowanceChargeReason,omitempty"`
	MultiplierFactorNumeric float64    `xml:"cbc:MultiplierFactorNumeric,omitempty"`
	Amount                  AmountType `xml:"cbc:Amount"`
	BaseAmount              AmountType `xml:"cbc:BaseAmount,omitempty"`
}

// DocumentReference representa una referencia a un documento
type DocumentReference struct {
	ID               string `xml:"cbc:ID"`
	IssueDate        string `xml:"cbc:IssueDate,omitempty"`
	DocumentTypeCode string `xml:"cbc:DocumentTypeCode,omitempty"`
	DocumentType     string `xml:"cbc:DocumentType,omitempty"`
}

// PricingReference representa referencias de precios alternativos
type PricingReference struct {
	AlternativeConditionPrice []AlternativeConditionPrice `xml:"cac:AlternativeConditionPrice,omitempty"`
}

// AlternativeConditionPrice representa un precio alternativo
type AlternativeConditionPrice struct {
	PriceAmount   AmountType `xml:"cbc:PriceAmount"`
	PriceTypeCode string     `xml:"cbc:PriceTypeCode,omitempty"`
	PriceType     string     `xml:"cbc:PriceType,omitempty"`
}

// Item representa un item/producto
type Item struct {
	Description                  string              `xml:"cbc:Description"`
	BrandName                    string              `xml:"cbc:BrandName,omitempty"`
	ModelName                    string              `xml:"cbc:ModelName,omitempty"`
	SellersItemIdentification    *ItemIdentification `xml:"cac:SellersItemIdentification,omitempty"`
	StandardItemIdentification   *ItemIdentification `xml:"cac:StandardItemIdentification,omitempty"`
	AdditionalItemIdentification *ItemIdentification `xml:"cac:AdditionalItemIdentification,omitempty"`
}

// ItemIdentification representa la identificación de un item
type ItemIdentification struct {
	ID IDType `xml:"cbc:ID"`
}

// Price representa el precio de un item
type Price struct {
	PriceAmount  AmountType `xml:"cbc:PriceAmount"`
	BaseQuantity Quantity   `xml:"cbc:BaseQuantity,omitempty"`
}

// NewInvoice crea una nueva factura con valores por defecto
func NewInvoice(id string) *Invoice {
	now := time.Now()
	return &Invoice{
		XmlnsCbc:           "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		XmlnsExt:           "urn:oasis:names:specification:ubl:schema:xsd:CommonExtensionComponents-2",
		XmlnsCac:           "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		XmlnsSts:           "dian:gov:co:facturaelectronica:Structures-2-1",
		UBLVersionID:       "UBL 2.1",
		CustomizationID:    "05",
		ProfileID:          "DIAN 2.1",
		ProfileExecutionID: "2",
		ID:                 id,
		UUID: UUIDType{
			SchemeID:   "2",
			SchemeName: "CUFE-SHA384",
		},
		IssueDate:       now.Format("2006-01-02"),
		IssueTime:       now.Format("15:04:05-07:00"),
		InvoiceTypeCode: "01",
		DocumentCurrencyCode: DocumentCurrencyType{
			Value:          "COP",
			ListAgencyID:   "6",
			ListAgencyName: "United Nations Economic Commission for Europe",
			ListID:         "ISO 4217 Alpha",
		},
		LineCountNumeric: 0,
		TaxTotal:         []TaxTotal{},
		InvoiceLines:     []InvoiceLine{},
	}
}

// AddLine agrega una línea a la factura
func (i *Invoice) AddLine(line InvoiceLine) {
	i.InvoiceLines = append(i.InvoiceLines, line)
	i.LineCountNumeric = len(i.InvoiceLines)
}

// CalculateTotals calcula los totales de la factura
func (i *Invoice) CalculateTotals() {
	var lineExtension, taxExclusive, taxInclusive, totalTax float64

	for _, line := range i.InvoiceLines {
		lineExtension += line.LineExtensionAmount.Value

		for _, tax := range line.TaxTotal {
			totalTax += tax.TaxAmount.Value
		}
	}

	taxExclusive = lineExtension
	taxInclusive = lineExtension + totalTax

	i.LegalMonetaryTotal = LegalMonetaryTotal{
		LineExtensionAmount: AmountType{Value: lineExtension, CurrencyID: "COP"},
		TaxExclusiveAmount:  AmountType{Value: taxExclusive, CurrencyID: "COP"},
		TaxInclusiveAmount:  AmountType{Value: taxInclusive, CurrencyID: "COP"},
		PayableAmount:       AmountType{Value: taxInclusive, CurrencyID: "COP"},
	}

	// Agregar TaxTotal a nivel de factura para CUFE
	if totalTax > 0 {
		i.TaxTotal = []TaxTotal{
			{
				TaxAmount: AmountType{Value: totalTax, CurrencyID: "COP"},
			},
		}
	}
}
