# go-dian v0.3.0

Librería Go profesional para Facturación Electrónica DIAN Colombia.

## Características

- ✅ Generación de facturas electrónicas UBL 2.1
- ✅ Extensiones DIAN (InvoiceControl, SoftwareProvider, QRCode)
- ✅ Cálculo CUFE SHA384
- ✅ Firma XMLDSig con certificados PEM
- ✅ Envío a DIAN vía SOAP
- ✅ Estructura modular y escalable

## Instalación

```bash
go get github.com/diegofxm/go-dian
```

## Uso Básico

```go
import (
    "github.com/diegofxm/go-dian/pkg/dian"
    "github.com/diegofxm/go-dian/pkg/invoice"
    "github.com/diegofxm/go-dian/pkg/common"
)

// Crear cliente DIAN
client, err := dian.NewClient(dian.Config{
    NIT:         "830122566",
    Environment: dian.EnvironmentTest,
    Certificate: dian.Certificate{
        PEMPath: "certificate.pem",
    },
})

// Crear factura
inv := invoice.NewInvoice("SETP990000001")
// ... configurar factura

// Generar y enviar
xml, err := client.GenerateInvoiceXML(inv)
response, err := client.SendInvoice(xml)
```

## Estructura

```
pkg/
├── dian/          Cliente principal DIAN
├── invoice/       Factura electrónica
├── common/        Tipos compartidos UBL
├── extensions/    Extensiones DIAN
├── signature/     Firma digital (solo PEM)
├── transmission/  Cliente SOAP
└── validation/    Validaciones DIAN
```

## Licencia

MIT
