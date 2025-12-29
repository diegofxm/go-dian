module github.com/diegofxm/go-dian

go 1.23

// v0.1.10 - Fixes críticos DIAN
// - CUFE con SHA384 (antes SHA256)
// - Eliminada notación científica en montos
// - Agregados PaymentMeans y PaymentTerms5.1

require software.sslmate.com/src/go-pkcs12 v0.7.0

require golang.org/x/crypto v0.11.0 // indirect
