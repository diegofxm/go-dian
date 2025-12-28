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

### ✅ Implementado (MVP v1.0)

- [x] Modelos UBL 2.1 completos para facturas electrónicas
- [x] Generación de XML conforme a estándar DIAN
- [x] Cálculo de CUFE (Código Único de Factura Electrónica)
- [x] Firma digital XMLDSig con soporte P12 y PEM
- [x] Cliente SOAP para envío a DIAN
- [x] Validaciones básicas de estructura
- [x] Utilidades y helpers
- [x] Soporte para ambientes de prueba y producción
- [x] Extensiones DIAN (InvoiceControl, SoftwareProvider, etc.)
- [x] Generación de código QR
- [x] Generación de código de seguridad de software

### 🚧 En Desarrollo / Próximas Versiones

#### v1.1 - Notas Crédito y Débito
- [ ] **Modelos para Notas Crédito** (`CreditNote`)
  - Estructura UBL 2.1 para notas crédito
  - Cálculo de CUDE (Código Único de Documento Electrónico)
  - Referencia a factura original
  - Motivos de devolución/descuento
- [ ] **Modelos para Notas Débito** (`DebitNote`)
  - Estructura UBL 2.1 para notas débito
  - Cálculo de CUDE
  - Referencia a factura original
  - Motivos de cargo adicional
- [ ] **Generación de XML** para notas crédito/débito
- [ ] **Firma digital** para notas crédito/débito
- [ ] **Envío SOAP** para notas crédito/débito

#### v1.2 - Documentos Soporte
- [ ] **Modelos para Documentos Soporte** (`SupportDocument`)
  - Estructura para documentos soporte de adquisiciones
  - Validaciones específicas para no obligados a facturar
  - Cálculo de CUDS (Código Único de Documento Soporte)
- [ ] **Generación de XML** para documentos soporte
- [ ] **Firma digital** para documentos soporte
- [ ] **Envío SOAP** para documentos soporte

#### v1.3 - Eventos DIAN
- [ ] **ApplicationResponse** - Modelo para eventos
  - Acuse de recibo
  - Aceptación expresa
  - Aceptación tácita
  - Rechazo de factura
  - Reclamo de factura
- [ ] **Generación de XML** para eventos
- [ ] **Firma digital** para eventos
- [ ] **Envío SOAP** para eventos
- [ ] **Consulta de estado** de documentos en DIAN
- [ ] **Notificaciones** de eventos recibidos

#### v1.4 - Validaciones Avanzadas
- [ ] **Validaciones de negocio DIAN**
  - Validación de rangos de numeración autorizados
  - Validación de fechas de vigencia de autorización
  - Validación de resoluciones DIAN
  - Validación de códigos DANE (ciudades, departamentos)
  - Validación de códigos UNSPSC (productos)
- [ ] **Validaciones de esquema XML**
  - Validación contra XSD oficial DIAN
  - Validación de firma digital
  - Validación de CUFE/CUDE
- [ ] **Validaciones de datos**
  - Validación de NIT con dígito de verificación
  - Validación de formatos de documentos
  - Validación de montos y cálculos

#### v1.5 - Mejoras de Infraestructura
- [ ] **Retry Logic** para envío SOAP
  - Reintentos automáticos con backoff exponencial
  - Manejo de timeouts
  - Circuit breaker para protección
  - Logging detallado de reintentos
- [ ] **Caché de respuestas DIAN**
  - Caché de consultas de estado
  - Reducción de llamadas redundantes
- [ ] **Métricas y observabilidad**
  - Métricas de tiempo de respuesta
  - Contadores de éxito/error
  - Trazabilidad de documentos

#### v2.0 - Nómina Electrónica
- [ ] **Modelos UBL para nómina**
  - Estructura de nómina individual
  - Ajustes de nómina
  - Notas de ajuste
- [ ] **Cálculo de CUNE** (Código Único de Nómina Electrónica)
- [ ] **Generación de XML** para nómina
- [ ] **Firma digital** para nómina
- [ ] **Envío SOAP** para nómina

### 💡 Mejoras Opcionales Futuras
- [ ] Soporte para múltiples certificados digitales
- [ ] Integración con HSM (Hardware Security Module)
- [ ] API REST wrapper sobre la librería
- [ ] Dashboard de monitoreo
- [ ] Exportación a PDF de documentos
- [ ] Integración con proveedores de firma en la nube
- [ ] Soporte para facturación masiva (batch processing)
- [ ] Webhooks para notificaciones asíncronas

## Estado Actual

**Versión:** 1.0.0 (MVP)

**Funcionalidad Core Completa:**
- ✅ Facturación electrónica básica
- ✅ Generación de XML UBL 2.1
- ✅ Firma digital XMLDSig
- ✅ Envío a DIAN vía SOAP
- ✅ Cálculo de CUFE
- ✅ Extensiones DIAN

**Lo que NO incluye (pero está en roadmap):**
- ❌ Notas crédito/débito (v1.1)
- ❌ Documentos soporte (v1.2)
- ❌ Eventos DIAN - aceptación/rechazo (v1.3)
- ❌ Validaciones avanzadas exhaustivas (v1.4)
- ❌ Retry logic automático en SOAP (v1.5)
- ❌ Nómina electrónica (v2.0)

## Licencia

MIT

## Contribuciones

Las contribuciones son bienvenidas. Por favor abre un issue o PR.

### Áreas prioritarias para contribuir:
1. Implementación de notas crédito/débito
2. Validaciones avanzadas contra anexos técnicos DIAN
3. Retry logic y manejo de errores robusto
4. Tests de integración con ambiente de pruebas DIAN
5. Documentación y ejemplos adicionales
