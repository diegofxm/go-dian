package invoice

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/diegofxm/go-dian/pkg/common"
)

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

	InvoicePeriod           *InvoicePeriod            `xml:"cac:InvoicePeriod,omitempty"`
	BillingReference        []BillingReference        `xml:"cac:BillingReference,omitempty"`
	AccountingSupplierParty AccountingSupplierParty   `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty AccountingCustomerParty   `xml:"cac:AccountingCustomerParty"`
	TaxRepresentativeParty  *TaxRepresentativeParty   `xml:"cac:TaxRepresentativeParty,omitempty"`
	Delivery                *Delivery                 `xml:"cac:Delivery,omitempty"`
	DeliveryTerms           *DeliveryTerms            `xml:"cac:DeliveryTerms,omitempty"`
	PaymentMeans            []common.PaymentMeans     `xml:"cac:PaymentMeans,omitempty"`
	PaymentTerms            []common.PaymentTerms     `xml:"cac:PaymentTerms,omitempty"`
	PrepaidPayment          []common.PrepaidPayment   `xml:"cac:PrepaidPayment,omitempty"`
	TaxTotal                []common.TaxTotal         `xml:"cac:TaxTotal"`
	LegalMonetaryTotal      common.LegalMonetaryTotal `xml:"cac:LegalMonetaryTotal"`
	InvoiceLines            []InvoiceLine             `xml:"cac:InvoiceLine"`
}

type UBLExtensions struct {
	UBLExtension []UBLExtension `xml:"ext:UBLExtension"`
}

type UBLExtension struct {
	ExtensionContent ExtensionContent `xml:"ext:ExtensionContent"`
}

type ExtensionContent struct {
	Content interface{} `xml:",innerxml"`
}

type UUIDType struct {
	Value      string `xml:",chardata"`
	SchemeID   string `xml:"schemeID,attr"`
	SchemeName string `xml:"schemeName,attr"`
}

type DocumentCurrencyType struct {
	Value          string `xml:",chardata"`
	ListAgencyID   string `xml:"listAgencyID,attr,omitempty"`
	ListAgencyName string `xml:"listAgencyName,attr,omitempty"`
	ListID         string `xml:"listID,attr,omitempty"`
}

type InvoicePeriod struct {
	StartDate string `xml:"cbc:StartDate"`
	EndDate   string `xml:"cbc:EndDate"`
}

type BillingReference struct {
	InvoiceDocumentReference InvoiceDocumentReference `xml:"cac:InvoiceDocumentReference"`
}

type InvoiceDocumentReference struct {
	ID                  string   `xml:"cbc:ID"`
	UUID                UUIDType `xml:"cbc:UUID,omitempty"`
	IssueDate           string   `xml:"cbc:IssueDate,omitempty"`
	DocumentDescription string   `xml:"cbc:DocumentDescription,omitempty"`
}

type AccountingSupplierParty struct {
	AdditionalAccountID common.AdditionalAccountIDType `xml:"cbc:AdditionalAccountID"`
	Party               common.Party                   `xml:"cac:Party"`
}

type AccountingCustomerParty struct {
	AdditionalAccountID common.AdditionalAccountIDType `xml:"cbc:AdditionalAccountID"`
	Party               common.Party                   `xml:"cac:Party"`
}

type TaxRepresentativeParty struct {
	PartyIdentification common.PartyIdentification `xml:"cac:PartyIdentification"`
}

type Delivery struct {
	DeliveryAddress *common.Address `xml:"cac:DeliveryAddress,omitempty"`
	DeliveryParty   *DeliveryParty  `xml:"cac:DeliveryParty,omitempty"`
}

type DeliveryParty struct {
	PartyName        []common.PartyName       `xml:"cac:PartyName,omitempty"`
	PhysicalLocation *common.PhysicalLocation `xml:"cac:PhysicalLocation,omitempty"`
	PartyTaxScheme   common.PartyTaxScheme    `xml:"cac:PartyTaxScheme"`
	PartyLegalEntity common.PartyLegalEntity  `xml:"cac:PartyLegalEntity"`
	Contact          *common.Contact          `xml:"cac:Contact,omitempty"`
}

type DeliveryTerms struct {
	ID                     string            `xml:"cbc:ID,omitempty"`
	SpecialTerms           string            `xml:"cbc:SpecialTerms,omitempty"`
	LossRiskResponsibility string            `xml:"cbc:LossRiskResponsibilityCode,omitempty"`
	DeliveryLocation       *DeliveryLocation `xml:"cac:DeliveryLocation,omitempty"`
}

type DeliveryLocation struct {
	ID common.IDType `xml:"cbc:ID"`
}

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
		TaxTotal:         []common.TaxTotal{},
		InvoiceLines:     []InvoiceLine{},
	}
}

func (i *Invoice) AddLine(line InvoiceLine) {
	i.InvoiceLines = append(i.InvoiceLines, line)
	i.LineCountNumeric = len(i.InvoiceLines)
}

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

	i.LegalMonetaryTotal = common.LegalMonetaryTotal{
		LineExtensionAmount: common.AmountType{Value: lineExtension, CurrencyID: "COP"},
		TaxExclusiveAmount:  common.AmountType{Value: taxExclusive, CurrencyID: "COP"},
		TaxInclusiveAmount:  common.AmountType{Value: taxInclusive, CurrencyID: "COP"},
		PayableAmount:       common.AmountType{Value: taxInclusive, CurrencyID: "COP"},
	}

	if totalTax > 0 {
		i.TaxTotal = []common.TaxTotal{
			{
				TaxAmount: common.AmountType{Value: totalTax, CurrencyID: "COP"},
			},
		}
	}
}
