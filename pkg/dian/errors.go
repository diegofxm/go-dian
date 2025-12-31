package dian

import "fmt"

// Error types
var (
	ErrInvalidNIT         = fmt.Errorf("NIT inválido")
	ErrMissingCertificate = fmt.Errorf("certificado no configurado")
	ErrInvalidInvoice     = fmt.Errorf("factura inválida")
)
