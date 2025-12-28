# go-dian

Paquete Go para integración con DIAN (Facturación Electrónica Colombia).

## Características

- ✅ Generación de XML UBL 2.1
- ✅ Cálculo de CUFE/CUDE
- ✅ Firma digital XMLDSig
- ✅ Envío a DIAN vía SOAP
- ✅ Validación de facturas
- ✅ Soporte para ambiente de pruebas y producción

## Instalación

```bash
go get github.com/diegofxm/go-dian
```

## Uso básico

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/diegofxm/go-dian"
)

func main() {
    // Configurar cliente
    client, err := dian.NewClient(dian.Config{
        NIT:         "830122566",
        Environment: dian.EnvironmentTest,
        SoftwareID:  "tu-software-id",
        Certificate: dian.Certificate{
            Path:     "./certificado.p12",
            Password: "password",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Crear factura
    invoice := dian.NewInvoice("BEC496329154")
    
    // Configurar emisor
    invoice.AccountingSupplierParty = dian.AccountingSupplierParty{
        AdditionalAccountID: "1",
        Party: dian.Party{
            PartyTaxScheme: dian.PartyTaxScheme{
                RegistrationName: "MI EMPRESA SAS",
                CompanyID: dian.IDType{
                    Value:      "830122566",
                    SchemeID:   "1",
                    SchemeName: "31",
                },
                TaxLevelCode: "O-13",
                TaxScheme: dian.TaxScheme{
                    ID:   "01",
                    Name: "IVA",
                },
            },
            PartyLegalEntity: dian.PartyLegalEntity{
                RegistrationName: "MI EMPRESA SAS",
                CompanyID: dian.IDType{
                    Value: "830122566",
                },
            },
        },
    }

    // Configurar cliente
    invoice.AccountingCustomerParty = dian.AccountingCustomerParty{
        AdditionalAccountID: "2",
        Party: dian.Party{
            PartyTaxScheme: dian.PartyTaxScheme{
                RegistrationName: "CLIENTE EJEMPLO",
                CompanyID: dian.IDType{
                    Value:      "900123456",
                    SchemeID:   "1",
                    SchemeName: "31",
                },
                TaxLevelCode: "O-13",
                TaxScheme: dian.TaxScheme{
                    ID:   "01",
                    Name: "IVA",
                },
            },
            PartyLegalEntity: dian.PartyLegalEntity{
                RegistrationName: "CLIENTE EJEMPLO",
                CompanyID: dian.IDType{
                    Value: "900123456",
                },
            },
        },
    }

    // Agregar línea de factura
    invoice.AddLine(dian.InvoiceLine{
        ID: "1",
        InvoicedQuantity: dian.Quantity{
            Value:    1,
            UnitCode: "94",
        },
        LineExtensionAmount: dian.AmountType{
            Value:      100000,
            CurrencyID: "COP",
        },
        Item: dian.Item{
            Description: "Producto de ejemplo",
        },
        Price: dian.Price{
            PriceAmount: dian.AmountType{
                Value:      100000,
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

    // Calcular totales
    invoice.CalculateTotals()

    // Generar XML
    xmlData, err := client.GenerateInvoiceXML(invoice)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(xmlData))

    // Enviar a DIAN
    response, err := client.SendInvoice(invoice)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Respuesta DIAN: %+v\n", response)
}
```

## Estructura del paquete

```
go-dian/
├── dian.go          # Cliente principal y lógica de negocio
├── models.go        # Modelos de datos UBL 2.1
├── signature.go     # Firma digital XMLDSig
├── soap.go          # Cliente SOAP para envío a DIAN
├── helpers.go       # Utilidades y funciones auxiliares
├── *_test.go        # Tests unitarios
├── examples/        # Ejemplos de uso
└── README.md        # Documentación
```

## Roadmap

- [x] Modelos UBL 2.1
- [x] Generación de XML
- [x] Cálculo de CUFE
- [x] Firma digital XMLDSig
- [x] Cliente SOAP
- [x] Validaciones básicas
- [x] Utilidades y helpers
- [ ] Notas crédito/débito
- [ ] Documento soporte
- [ ] Nómina electrónica
- [ ] Validaciones avanzadas DIAN

## Licencia

MIT

## Contribuciones

Las contribuciones son bienvenidas. Por favor abre un issue o PR.
