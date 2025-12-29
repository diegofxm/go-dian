package common

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
