package common

// LegalMonetaryTotal representa el total monetario legal
type LegalMonetaryTotal struct {
	LineExtensionAmount AmountType  `xml:"cbc:LineExtensionAmount"`
	TaxExclusiveAmount  AmountType  `xml:"cbc:TaxExclusiveAmount"`
	TaxInclusiveAmount  AmountType  `xml:"cbc:TaxInclusiveAmount"`
	PrepaidAmount       *AmountType `xml:"cbc:PrepaidAmount,omitempty"`
	PayableAmount       AmountType  `xml:"cbc:PayableAmount"`
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
