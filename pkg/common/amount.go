package common

import (
	"encoding/xml"
	"fmt"
)

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
