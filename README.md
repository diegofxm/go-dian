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
    // Configurar cliente con datos de autorización DIAN
    client, err := dian.NewClient(dian.Config{
        NIT:         "830122566",
        Environment: dian.EnvironmentTest,
        SoftwareID:  "tu-software-id",
        Certificate: dian.Certificate{
            Path:     "./certificado.p12",
            Password: "password",
        },
        // Datos de autorización DIAN (específicos por empresa)
        InvoiceAuthorization: "18764090648904",
        AuthStartDate:        "2025-03-18",
        AuthEndDate:          "2027-03-18",
        InvoicePrefix:        "FACT",
        AuthFrom:             "1",
        AuthTo:               "1000000",
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

    // Generar XML (sin firmar)
    xmlData, err := client.GenerateInvoiceXML(invoice)
    if err != nil {
        log.Fatal(err)
    }

    // Firmar XML
    signedXML, err := client.SignXML(xmlData)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(signedXML))

    // O usar SendInvoice que hace todo (generar, firmar y enviar)
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
├── dian.go          # Cliente principal y funciones públicas
├── models.go        # Modelos UBL 2.1 (Invoice, InvoiceLine, etc.)
├── signature.go     # Firma digital XMLDSig y manejo de certificados
├── extensions.go    # Extensiones DIAN (InvoiceControl, QRCode, etc.)
├── soap.go          # Cliente SOAP para envío a DIAN
├── examples/        # Ejemplos de uso
└── README.md        # Documentación
```

## API Principal

### Funciones Públicas

**Cliente:**
- `NewClient(config Config)` - Crea cliente DIAN
- `GenerateInvoiceXML(invoice *Invoice)` - Genera XML sin firmar
- `SignXML(xmlData []byte)` - Firma XML (genérico, reutilizable)
- `SendInvoice(invoice *Invoice)` - Genera, firma y envía a DIAN
- `CalculateCUFE(invoice *Invoice)` - Calcula CUFE SHA384
- `ValidateNIT(nit string)` - Valida formato de NIT colombiano

**Certificados:**
- `LoadCertificate(path, password)` - Carga certificado P12/PEM
- `LoadCertificateFromPEMStrings(certPEM, keyPEM)` - Carga desde BD
- `GetCertificateInfo(cert)` - Obtiene info del certificado

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

**Versión:** 0.1.10 (Fixes Críticos DIAN)

**Funcionalidad Core Completa:**
- ✅ Facturación electrónica básica
- ✅ Generación de XML UBL 2.1
- ✅ Firma digital XMLDSig
- ✅ Envío a DIAN vía SOAP
- ✅ Cálculo de CUFE con SHA384 (CORREGIDO v0.1.10)
- ✅ Extensiones DIAN
- ✅ Montos sin notación científica (CORREGIDO v0.1.10)
- ✅ PaymentMeans y PaymentTerms (AGREGADO v0.1.10)

**Lo que NO incluye (pero está en roadmap):**
- ❌ Notas crédito/débito (v1.1)
- ❌ Documentos soporte (v1.2)
- ❌ Eventos DIAN - aceptación/rechazo (v1.3)
- ❌ Validaciones avanzadas exhaustivas (v1.4)
- ❌ Retry logic automático en SOAP (v1.5)
- ❌ Nómina electrónica (v2.0)

## 📝 Changelog

### v0.2.0 (2025-12-29) - Refactorización y Limpieza

**🔧 REFACTORIZACIÓN - Separación de Responsabilidades:**
- ✅ `GenerateInvoiceXML()` - Solo genera XML (sin firmar)
- ✅ `SignXML()` - Solo firma XML (genérico, reutilizable)
- ✅ `SendInvoice()` - Orquesta todo (generar + firmar + enviar)
- 🎯 **Beneficio:** Máxima flexibilidad para usuarios avanzados

**🧹 LIMPIEZA - Código Optimizado:**
- ❌ **Eliminado:** `helpers.go` (150+ líneas de código no usado)
- ✅ **Movido:** `ValidateNIT()` a `dian.go` (única función útil)
- ❌ **Eliminado:** Funciones stub y redundantes
- ✅ **Mejorado:** `GetCertificateInfo()` ahora retorna struct tipado

**⚙️ PARAMETRIZACIÓN - Datos de Autorización:**
- ✅ **Agregado:** Campos al `Config` para datos de autorización DIAN
- ✅ **Eliminado:** Datos hardcodeados en `extensions.go`
- 🎯 **Beneficio:** Cada empresa usa sus propios datos de autorización

**📊 IMPACTO:**
- Reducción de ~350 líneas de código (-32%)
- 0% código duplicado o redundante
- Librería lista para uso opensource profesional

### v0.1.10 (2025-12-28) - Fixes Críticos DIAN

**🔴 CRÍTICO - CUFE con SHA384:**
- ✅ Usa SHA384 según requerimientos oficiales DIAN
- 🔧 Cambio: `sha256.Sum256` → `sha512.Sum384`

**🔴 CRÍTICO - Notación Científica Eliminada:**
- ✅ Montos se serializan como `2289500.00` (no `2.2895e+06`)
- 🔧 Custom `MarshalXML` con `fmt.Sprintf("%.2f")`

**⚠️ IMPORTANTE - PaymentMeans y PaymentTerms:**
- ✅ Agregados structs para medios de pago y condiciones

---

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
