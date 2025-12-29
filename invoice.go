package dian

import "github.com/diegofxm/go-dian/invoice"

// NewInvoice crea una nueva factura con valores por defecto
func NewInvoice(id string) *invoice.Invoice {
	return invoice.NewInvoice(id)
}
