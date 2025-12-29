package invoice

import "github.com/diegofxm/go-dian/common"

// InvoiceLine representa una línea de factura
type InvoiceLine struct {
	ID                    string                     `xml:"cbc:ID"`
	InvoicedQuantity      common.Quantity            `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount   common.AmountType          `xml:"cbc:LineExtensionAmount"`
	FreeOfChargeIndicator *bool                      `xml:"cbc:FreeOfChargeIndicator,omitempty"`
	Delivery              *InvoiceLineDelivery       `xml:"cac:Delivery,omitempty"`
	AllowanceCharge       []common.AllowanceCharge   `xml:"cac:AllowanceCharge,omitempty"`
	TaxTotal              []common.TaxTotal          `xml:"cac:TaxTotal,omitempty"`
	DocumentReference     []common.DocumentReference `xml:"cac:DocumentReference,omitempty"`
	PricingReference      *PricingReference          `xml:"cac:PricingReference,omitempty"`
	Item                  Item                       `xml:"cac:Item"`
	Price                 Price                      `xml:"cac:Price"`
}

// InvoiceLineDelivery representa la entrega de una línea
type InvoiceLineDelivery struct {
	DeliveryLocation *DeliveryLocation `xml:"cac:DeliveryLocation,omitempty"`
}

// PricingReference representa referencias de precios alternativos
type PricingReference struct {
	AlternativeConditionPrice []AlternativeConditionPrice `xml:"cac:AlternativeConditionPrice,omitempty"`
}

// AlternativeConditionPrice representa un precio alternativo
type AlternativeConditionPrice struct {
	PriceAmount   common.AmountType `xml:"cbc:PriceAmount"`
	PriceTypeCode string            `xml:"cbc:PriceTypeCode,omitempty"`
	PriceType     string            `xml:"cbc:PriceType,omitempty"`
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
	ID common.IDType `xml:"cbc:ID"`
}

// Price representa el precio de un item
type Price struct {
	PriceAmount  common.AmountType `xml:"cbc:PriceAmount"`
	BaseQuantity common.Quantity   `xml:"cbc:BaseQuantity,omitempty"`
}
