# Guía de Uso - go-dian

## Instalación

```bash
go get github.com/diegofxm/go-dian
```

Dependencias requeridas:
```bash
go get software.sslmate.com/src/go-pkcs12
```

## Configuración Inicial

### 1. Obtener certificado digital

Necesitas un certificado digital `.p12` emitido por una entidad certificadora autorizada por DIAN (GSE, Certicámara, etc.).

### 2. Configurar cliente

```go
import "github.com/diegofxm/go-dian"

client, err := dian.NewClient(dian.Config{
    NIT:         "830122566",
    Environment: dian.EnvironmentTest, // o EnvironmentProduction
    SoftwareID:  "tu-software-id-dian",
    TestSetID:   "tu-test-set-id", // Solo para habilitación
    Certificate: dian.Certificate{
        Path:     "./certificado.p12",
        Password: "password-certificado",
    },
})
```

## Crear una Factura

### Factura básica

```go
// Crear factura
invoice := dian.NewInvoice("BEC496329154")

// Configurar emisor
invoice.AccountingSupplierParty = dian.AccountingSupplierParty{
    AdditionalAccountID: dian.GetPersonTypeCode("JURIDICA"),
    Party: dian.Party{
        IndustryClassificationCode: "6120",
        PartyIdentification: dian.PartyIdentification{
            ID: dian.IDType{
                Value:            "830122566",
                SchemeID:         dian.GetIdentificationTypeCode("NIT"),
                SchemeName:       "31",
                SchemeAgencyID:   "195",
                SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
            },
        },
        PartyName: dian.PartyName{
            Name: "MI EMPRESA SAS",
        },
        PartyTaxScheme: dian.PartyTaxScheme{
            RegistrationName: "MI EMPRESA SAS",
            CompanyID: dian.IDType{
                Value:      "830122566",
                SchemeID:   "1",
                SchemeName: "31",
            },
            TaxLevelCode: "O-13",
            TaxScheme: dian.TaxScheme{
                ID:   dian.GetTaxSchemeCode("IVA"),
                Name: "IVA",
            },
        },
        PartyLegalEntity: dian.PartyLegalEntity{
            RegistrationName: "MI EMPRESA SAS",
            CompanyID: dian.IDType{Value: "830122566"},
        },
        Contact: &dian.Contact{
            Telephone:      "3001234567",
            ElectronicMail: "contacto@miempresa.com",
        },
    },
}

// Configurar cliente
invoice.AccountingCustomerParty = dian.AccountingCustomerParty{
    AdditionalAccountID: dian.GetPersonTypeCode("NATURAL"),
    Party: dian.Party{
        PartyIdentification: dian.PartyIdentification{
            ID: dian.IDType{
                Value:      "6382356",
                SchemeID:   dian.GetIdentificationTypeCode("CC"),
                SchemeName: "13",
            },
        },
        PartyName: dian.PartyName{
            Name: "CLIENTE EJEMPLO",
        },
        PartyTaxScheme: dian.PartyTaxScheme{
            RegistrationName: "CLIENTE EJEMPLO",
            CompanyID: dian.IDType{Value: "6382356"},
            TaxLevelCode: "R-99-PN",
            TaxScheme: dian.TaxScheme{
                ID:   "01",
                Name: "IVA",
            },
        },
        PartyLegalEntity: dian.PartyLegalEntity{
            RegistrationName: "CLIENTE EJEMPLO",
            CompanyID: dian.IDType{Value: "6382356"},
        },
    },
}
```

### Agregar líneas de factura

```go
// Producto con IVA 19%
invoice.AddLine(dian.InvoiceLine{
    ID: "1",
    InvoicedQuantity: dian.Quantity{
        Value:    2,
        UnitCode: "94", // Unidad
    },
    LineExtensionAmount: dian.AmountType{
        Value:      100000,
        CurrencyID: "COP",
    },
    Item: dian.Item{
        Description: "Producto de ejemplo",
        BrandName:   "Marca",
        ModelName:   "Modelo",
    },
    Price: dian.Price{
        PriceAmount: dian.AmountType{
            Value:      50000,
            CurrencyID: "COP",
        },
    },
    TaxTotal: []dian.TaxTotal{
        {
            TaxAmount: dian.AmountType{
                Value:      19000,
                CurrencyID: "COP",
            },
            TaxSubtotal: []dian.TaxSubtotal{
                {
                    TaxableAmount: dian.AmountType{
                        Value:      100000,
                        CurrencyID: "COP",
                    },
                    TaxAmount: dian.AmountType{
                        Value:      19000,
                        CurrencyID: "COP",
                    },
                    TaxCategory: dian.TaxCategory{
                        Percent: 19.0,
                        TaxScheme: dian.TaxScheme{
                            ID:   "01",
                            Name: "IVA",
                        },
                    },
                },
            },
        },
    },
})

// Calcular totales automáticamente
invoice.CalculateTotals()
```

## Enviar a DIAN

### Flujo completo

```go
// 1. Generar XML
xmlData, err := client.GenerateInvoiceXML(invoice)
if err != nil {
    log.Fatal(err)
}

// 2. Firmar XML
signedXML, err := client.SignXML(xmlData)
if err != nil {
    log.Fatal(err)
}

// 3. Enviar a DIAN
response, err := client.SendInvoice(invoice)
if err != nil {
    log.Fatal(err)
}

// 4. Verificar respuesta
if response.Success {
    fmt.Printf("✅ Factura enviada exitosamente\n")
    fmt.Printf("CUFE: %s\n", response.CUFE)
} else {
    fmt.Printf("❌ Error: %s\n", response.Message)
    for _, err := range response.Errors {
        fmt.Printf("  - %s\n", err)
    }
}
```

### Solo generar XML (sin enviar)

```go
xmlData, err := client.GenerateInvoiceXML(invoice)
if err != nil {
    log.Fatal(err)
}

// Guardar XML
os.WriteFile("factura.xml", xmlData, 0644)
```

## Utilidades

### Validaciones

```go
// Validar NIT
if err := dian.ValidateNIT("830122566"); err != nil {
    log.Fatal(err)
}

// Validar email
if err := dian.ValidateEmail("test@example.com"); err != nil {
    log.Fatal(err)
}
```

### Cálculos

```go
// Calcular IVA
iva := dian.CalculateIVA(100000, 19) // 19000

// Calcular total con impuesto
total := dian.CalculateTotalWithTax(100000, 19) // 119000

// Redondear
rounded := dian.RoundToDecimals(123.456, 2) // 123.46
```

### Formateo

```go
// Formatear fecha
dateStr := dian.FormatDate(time.Now()) // "2025-12-14"

// Formatear fecha y hora
dateTimeStr := dian.FormatDateTime(time.Now()) // "2025-12-14T08:17:34-05:00"

// Generar número de factura
invoiceNum := dian.GenerateInvoiceNumber("BEC", 123) // "BEC123"
```

## Ambientes

### Pruebas (Habilitación)

```go
client, _ := dian.NewClient(dian.Config{
    Environment: dian.EnvironmentTest,
    // ... resto de configuración
})
```

### Producción

```go
client, _ := dian.NewClient(dian.Config{
    Environment: dian.EnvironmentProduction,
    // ... resto de configuración
})
```

## Manejo de Errores

```go
response, err := client.SendInvoice(invoice)
if err != nil {
    // Error de conexión o técnico
    log.Printf("Error técnico: %v", err)
    return
}

if !response.Success {
    // Error de validación DIAN
    log.Printf("Error DIAN: %s", response.Message)
    for _, validationErr := range response.Errors {
        log.Printf("  - %s", validationErr)
    }
    return
}

// Éxito
fmt.Printf("CUFE: %s\n", response.CUFE)
```

## Mejores Prácticas

1. **Validar antes de enviar**
   ```go
   if err := invoice.Validate(); err != nil {
       log.Fatal(err)
   }
   ```

2. **Guardar XMLs generados**
   ```go
   os.WriteFile(fmt.Sprintf("facturas/%s.xml", invoice.ID), xmlData, 0644)
   ```

3. **Manejar certificados de forma segura**
   ```go
   // Usar variables de entorno
   certPath := os.Getenv("DIAN_CERT_PATH")
   certPass := os.Getenv("DIAN_CERT_PASSWORD")
   ```

4. **Logs detallados**
   ```go
   log.Printf("Enviando factura %s a DIAN...", invoice.ID)
   log.Printf("CUFE calculado: %s", invoice.UUID)
   ```

## Soporte

Para más información, consulta:
- [Documentación DIAN](https://www.dian.gov.co)
- [Anexos técnicos](https://www.dian.gov.co/impuestos/factura-electronica)
- [Issues en GitHub](https://github.com/diegofxm/go-dian/issues)
