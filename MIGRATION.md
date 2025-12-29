# Migración a Arquitectura Modular v0.3.0

**Fecha:** 2025-12-29  
**Versión:** v0.2.0 → v0.3.0  
**Tipo:** Refactorización interna (sin breaking changes)

---

## Objetivo

Reorganizar go-dian en una estructura modular escalable que soporte:
- Facturas electrónicas (actual)
- Notas crédito/débito (v1.1)
- Documentos soporte (v1.2)
- Eventos DIAN (v1.3)
- Nómina electrónica (v2.0)

**Garantía:** API pública NO cambia. Código de usuario sigue funcionando sin modificaciones.

---

## Estructura Anterior (v0.2.0)

```
go-dian/
├── dian.go          # 228 líneas - Cliente + Invoice
├── models.go        # 482 líneas - Structs UBL
├── signature.go     # 450 líneas - Firma XMLDSig
├── extensions.go    # 184 líneas - Extensiones DIAN
├── soap.go          # 204 líneas - Cliente SOAP
└── examples/
    └── basic/
        └── main.go
```

**Total:** 1,548 líneas en 5 archivos

---

## Estructura Nueva (v0.3.0)

**Basada en mejores prácticas de Go (sin `pkg/`):**

```
go-dian/
├── dian.go              # Cliente principal y API pública
├── config.go            # Config, Certificate, Environment
├── errors.go            # Errores personalizados
│
├── invoice.go           # API pública para Invoice
│
├── invoice/             # Paquete Invoice
│   ├── invoice.go       # Struct Invoice y métodos
│   ├── line.go          # Struct InvoiceLine
│   ├── generator.go     # Generación XML
│   └── cufe.go          # Cálculo CUFE
│
├── common/              # Structs compartidos UBL
│   ├── party.go         # Party, Address, Contact
│   ├── tax.go           # TaxTotal, TaxSubtotal
│   ├── amount.go        # AmountType, Quantity
│   └── monetary.go      # MonetaryTotal, AllowanceCharge
│
├── signature/           # Firma digital
│   ├── signature.go     # Servicio de firma
│   ├── xmldsig.go       # XMLDSig
│   └── certificate.go   # Manejo certificados
│
├── extensions/          # Extensiones DIAN
│   ├── invoice.go       # Extensiones Invoice
│   └── types.go         # DianExtensions structs
│
├── soap/                # Cliente SOAP
│   ├── client.go        # SOAPClient
│   ├── invoice.go       # Endpoints Invoice
│   └── types.go         # SOAP structs
│
├── internal/            # Código privado
│   ├── hash/
│   │   └── hash.go      # Funciones hash SHA384
│   └── xml/
│       └── utils.go     # Utilidades XML
│
├── examples/
│   └── basic/
│       └── main.go
│
├── go.mod
├── go.sum
├── README.md
├── MIGRATION.md
├── CHANGELOG.md
└── LICENSE
```

**Nota:** Solo se crean carpetas necesarias para v0.3.0. Carpetas futuras (creditnote, debitnote, etc.) se crearán cuando se implementen.

---

## Mapeo de Archivos

### Archivos Actuales → Nueva Ubicación

| Archivo Actual | Nueva Ubicación | Acción |
|----------------|-----------------|--------|
| `dian.go` | `dian.go` + `config.go` + `invoice.go` + `invoice/` | Separar responsabilidades |
| `models.go` | `invoice/invoice.go` + `invoice/line.go` + `common/` | Separar structs Invoice de comunes |
| `signature.go` | `signature/signature.go` + `signature/xmldsig.go` + `signature/certificate.go` | Mover a paquete signature |
| `extensions.go` | `extensions/invoice.go` + `extensions/types.go` | Mover a paquete extensions |
| `soap.go` | `soap/client.go` + `soap/invoice.go` + `soap/types.go` | Mover a paquete soap |

---

## Pasos de Migración

### 1. Crear estructura de carpetas
```cmd
REM Ejecutar en: C:\Users\codev\Desktop\project_e-invoince\go-dian
mkdir invoice
mkdir common
mkdir signature
mkdir extensions
mkdir soap
mkdir internal
mkdir internal\hash
mkdir internal\xml
```

### 2. Crear archivos base
- `config.go` - Extraer Config, Certificate de `dian.go`
- `errors.go` - Errores personalizados
- `invoice.go` - API pública para Invoice

### 3. Refactorizar archivos existentes
- `dian.go` → Mantener solo Client y NewClient
- `models.go` → Separar en `invoice/`, `common/`
- `signature.go` → Mover a `signature/`
- `extensions.go` → Mover a `extensions/`
- `soap.go` → Mover a `soap/`

### 4. Mantener compatibilidad
```go
// dian.go - API pública NO cambia
package dian

import "github.com/diegofxm/go-dian/invoice"

// GenerateInvoiceXML mantiene API pública
func (c *Client) GenerateInvoiceXML(inv *invoice.Invoice) ([]byte, error) {
    return invoice.Generate(c, inv)
}
```

### 5. Imports en api-dian NO cambian
```go
// Antes y después - MISMO CÓDIGO
import "github.com/diegofxm/go-dian"

client, _ := dian.NewClient(dian.Config{...})
```

---

## Compatibilidad

### API Pública NO Cambia

**Código de usuario sigue funcionando sin modificaciones:**
```go
// v0.2.0 y v0.3.0 - MISMO CÓDIGO
import "github.com/diegofxm/go-dian"

client, _ := dian.NewClient(dian.Config{...})
invoice := dian.NewInvoice("FACT-001")
xml, _ := client.GenerateInvoiceXML(invoice)
response, _ := client.SendInvoice(invoice)
```

### Imports Internos (solo para desarrollo)

```go
// Usuarios avanzados pueden importar paquetes específicos
import (
    "github.com/diegofxm/go-dian"
    "github.com/diegofxm/go-dian/invoice"    // Opcional
    "github.com/diegofxm/go-dian/common"     // Opcional
    "github.com/diegofxm/go-dian/signature"  // Opcional
)
```

---

## Beneficios

1. ✅ **Escalabilidad** - Fácil agregar nuevos tipos de documentos
2. ✅ **Mantenibilidad** - Archivos pequeños (~200 líneas)
3. ✅ **Organización** - Separación clara de responsabilidades
4. ✅ **Testing** - Fácil testear módulos individuales
5. ✅ **Colaboración** - Múltiples desarrolladores sin conflictos
6. ✅ **Compatibilidad** - API pública sin cambios

---

## Próximos Pasos

### v1.1 - Notas Crédito/Débito
- Crear `creditnote/` y `debitnote/`
- Agregar `creditnote.go` y `debitnote.go` en raíz
- Extender `extensions/` y `soap/`

### v1.2 - Documentos Soporte
- Crear `support/`
- Agregar `support.go` en raíz

### v1.3 - Eventos DIAN
- Crear `events/`
- Implementar ApplicationResponse, aceptación, rechazo

### v2.0 - Nómina Electrónica
- Crear `payroll/`
- Agregar `payroll.go` en raíz

---

## Notas Técnicas

- **Versión Go:** 1.21+
- **Breaking Changes:** Ninguno en v0.3.0
- **Deprecations:** Ninguno en v0.3.0
- **Tests:** Todos los tests existentes deben pasar sin modificación
