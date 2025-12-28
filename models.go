package dian

import (
	"encoding/xml"
	"fmt"
	"time"
)

// Invoice representa una factura electrónica UBL 2.1
type Invoice struct {
	XMLName xml.Name `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`

	UBLExtensions *UBLExtensions `xml:"ext:UBLExtensions,omitempty"`

	UBLVersionID         string `xml:"cbc:UBLVersionID"`
	CustomizationID      string `xml:"cbc:CustomizationID"`
	ProfileID            string `xml:"cbc:ProfileID"`
	ProfileExecutionID   string `xml:"cbc:ProfileExecutionID"`
	ID                   string `xml:"cbc:ID"`
	UUID                 string `xml:"cbc:UUID,attr"`
	IssueDate            string `xml:"cbc:IssueDate"`
	IssueTime            string `xml:"cbc:IssueTime"`
	DueDate              string `xml:"cbc:DueDate,omitempty"`
	InvoiceTypeCode      string `xml:"cbc:InvoiceTypeCode"`
	Note                 string `xml:"cbc:Note,omitempty"`
	DocumentCurrencyCode string `xml:"cbc:DocumentCurrencyCode"`
	LineCountNumeric     int    `xml:"cbc:LineCountNumeric"`

	InvoicePeriod           *InvoicePeriod          `xml:"cac:InvoicePeriod,omitempty"`
	AccountingSupplierParty AccountingSupplierParty `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty AccountingCustomerParty `xml:"cac:AccountingCustomerParty"`
	PaymentMeans            *PaymentMeans           `xml:"cac:PaymentMeans,omitempty"`
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

// InvoicePeriod representa el período de facturación
type InvoicePeriod struct {
	StartDate string `xml:"cbc:StartDate"`
	EndDate   string `xml:"cbc:EndDate"`
}

// AccountingSupplierParty representa el emisor de la factura
type AccountingSupplierParty struct {
	AdditionalAccountID string `xml:"cbc:AdditionalAccountID"`
	Party               Party  `xml:"cac:Party"`
}

// AccountingCustomerParty representa el cliente
type AccountingCustomerParty struct {
	AdditionalAccountID string `xml:"cbc:AdditionalAccountID"`
	Party               Party  `xml:"cac:Party"`
}

// Party representa una parte (emisor o cliente)
type Party struct {
	PartyIdentification        PartyIdentification `xml:"cac:PartyIdentification"`
	PartyName                  PartyName           `xml:"cac:PartyName"`
	PhysicalLocation           *Address            `xml:"cac:PhysicalLocation>cac:Address,omitempty"`
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
	RegistrationName    string    `xml:"cbc:RegistrationName"`
	CompanyID           IDType    `xml:"cbc:CompanyID"`
	TaxLevelCode        string    `xml:"cbc:TaxLevelCode"`
	RegistrationAddress *Address  `xml:"cac:RegistrationAddress,omitempty"`
	TaxScheme           TaxScheme `xml:"cac:TaxScheme"`
}

// PartyLegalEntity representa la entidad legal
type PartyLegalEntity struct {
	RegistrationName string `xml:"cbc:RegistrationName"`
	CompanyID        IDType `xml:"cbc:CompanyID"`
}

// Contact representa información de contacto
type Contact struct {
	Name           string `xml:"cbc:Name,omitempty"`
	Telephone      string `xml:"cbc:Telephone,omitempty"`
	ElectronicMail string `xml:"cbc:ElectronicMail,omitempty"`
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
	Name               string `xml:"cbc:Name"`
}

// IDType representa un ID con atributos
type IDType struct {
	Value            string `xml:",chardata"`
	SchemeID         string `xml:"schemeID,attr,omitempty"`
	SchemeName       string `xml:"schemeName,attr,omitempty"`
	SchemeAgencyID   string `xml:"schemeAgencyID,attr,omitempty"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr,omitempty"`
}

// PaymentMeans representa los medios de pago
type PaymentMeans struct {
	ID               string `xml:"cbc:ID"`
	PaymentMeansCode string `xml:"cbc:PaymentMeansCode"`
	PaymentDueDate   string `xml:"cbc:PaymentDueDate,omitempty"`
	PaymentID        string `xml:"cbc:PaymentID,omitempty"`
}

// TaxTotal representa el total de impuestos
type TaxTotal struct {
	TaxAmount      AmountType    `xml:"cbc:TaxAmount"`
	RoundingAmount AmountType    `xml:"cbc:RoundingAmount,omitempty"`
	TaxSubtotal    []TaxSubtotal `xml:"cac:TaxSubtotal"`
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

// AmountType representa un monto con moneda
type AmountType struct {
	Value      float64 `xml:",chardata"`
	CurrencyID string  `xml:"currencyID,attr"`
}

// LegalMonetaryTotal representa el total monetario legal
type LegalMonetaryTotal struct {
	LineExtensionAmount  AmountType `xml:"cbc:LineExtensionAmount"`
	TaxExclusiveAmount   AmountType `xml:"cbc:TaxExclusiveAmount"`
	TaxInclusiveAmount   AmountType `xml:"cbc:TaxInclusiveAmount"`
	AllowanceTotalAmount AmountType `xml:"cbc:AllowanceTotalAmount,omitempty"`
	ChargeTotalAmount    AmountType `xml:"cbc:ChargeTotalAmount,omitempty"`
	PayableAmount        AmountType `xml:"cbc:PayableAmount"`
}

// InvoiceLine representa una línea de factura
type InvoiceLine struct {
	ID                  string     `xml:"cbc:ID"`
	InvoicedQuantity    Quantity   `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount AmountType `xml:"cbc:LineExtensionAmount"`
	TaxTotal            []TaxTotal `xml:"cac:TaxTotal,omitempty"`
	Item                Item       `xml:"cac:Item"`
	Price               Price      `xml:"cac:Price"`
}

// Quantity representa una cantidad con unidad
type Quantity struct {
	Value    float64 `xml:",chardata"`
	UnitCode string  `xml:"unitCode,attr"`
}

// Item representa un item/producto
type Item struct {
	Description                string              `xml:"cbc:Description"`
	BrandName                  string              `xml:"cbc:BrandName,omitempty"`
	ModelName                  string              `xml:"cbc:ModelName,omitempty"`
	SellersItemIdentification  *ItemIdentification `xml:"cac:SellersItemIdentification,omitempty"`
	StandardItemIdentification *ItemIdentification `xml:"cac:StandardItemIdentification,omitempty"`
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
		UBLVersionID:         "UBL 2.1",
		CustomizationID:      "10",
		ProfileID:            "DIAN 2.1: Factura Electrónica de Venta",
		ProfileExecutionID:   "1",
		ID:                   id,
		IssueDate:            now.Format("2006-01-02"),
		IssueTime:            now.Format("15:04:05-07:00"),
		InvoiceTypeCode:      "01",
		DocumentCurrencyCode: "COP",
		LineCountNumeric:     0,
		TaxTotal:             []TaxTotal{},
		InvoiceLines:         []InvoiceLine{},
	}
}

// AddLine agrega una línea a la factura
func (i *Invoice) AddLine(line InvoiceLine) {
	i.InvoiceLines = append(i.InvoiceLines, line)
	i.LineCountNumeric = len(i.InvoiceLines)
}

// CalculateTotals calcula los totales de la factura
func (i *Invoice) CalculateTotals() {
	var lineExtension, taxExclusive, taxInclusive float64

	for _, line := range i.InvoiceLines {
		lineExtension += line.LineExtensionAmount.Value

		for _, tax := range line.TaxTotal {
			taxInclusive += tax.TaxAmount.Value
		}
	}

	taxExclusive = lineExtension
	taxInclusive = lineExtension + taxInclusive

	i.LegalMonetaryTotal = LegalMonetaryTotal{
		LineExtensionAmount: AmountType{Value: lineExtension, CurrencyID: "COP"},
		TaxExclusiveAmount:  AmountType{Value: taxExclusive, CurrencyID: "COP"},
		TaxInclusiveAmount:  AmountType{Value: taxInclusive, CurrencyID: "COP"},
		PayableAmount:       AmountType{Value: taxInclusive, CurrencyID: "COP"},
	}
}
