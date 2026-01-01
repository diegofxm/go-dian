# Prueba de Envío a DIAN - go-dian

Este ejemplo demuestra cómo generar, firmar y enviar una factura electrónica a DIAN usando la librería `go-dian`.

## 📋 Datos de Habilitación Configurados

```
URL: https://vpfe-hab.dian.gov.co/WcfDianCustomerServices.svc
TestSetId: e6784f41-2aba-4ed3-bcb6-d045ab217e72
SoftwareID: 23bf9eac-4dbe-4300-af06-541cc3efc7ca
Clave Técnica: fc8eac422eba16e22ffd8c6f94b3f40a6e38162c
PIN: 40125
Prefijo: SETP
Rango: 990000000 - 995000000
Cuota: 50 facturas (30 FE, 10 ND, 10 NC)
```

## 🚀 Ejecutar Prueba

### 1. Verificar Certificado

Asegúrate de que tu certificado PEM esté en:
```
go-dian-v0.3.0/examples/certificates/certificate.pem
```

El certificado debe contener:
- Certificado público (-----BEGIN CERTIFICATE-----)
- Clave privada (-----BEGIN PRIVATE KEY-----)

### 2. Compilar y Ejecutar

Desde la carpeta `examples/basic/`:

```bash
# Compilar
go build -o basic.exe main.go

# Ejecutar
./basic.exe
```

O directamente:
```bash
go run main.go
```

## 📤 Proceso de Envío

El programa ejecutará los siguientes pasos:

1. **Generar XML UBL 2.1**
   - Crea factura SETP990000001
   - Calcula CUFE (SHA-384)
   - Agrega DIAN Extensions
   - Guarda: `invoice_unsigned.xml`

2. **Firmar Digitalmente**
   - Firma con certificado PEM
   - Agrega XMLDSig en UBLExtensions
   - Guarda: `invoice_signed.xml`

3. **Validar Conexión DIAN**
   - Verifica que el endpoint esté disponible

4. **Enviar vía SOAP**
   - Codifica XML en base64
   - Envía a DIAN con TestSetId
   - Recibe ApplicationResponse

5. **Procesar Respuesta**
   - Muestra resultado (aceptada/rechazada)
   - Lista errores y advertencias
   - Guarda: `dian_response.xml`

## 📊 Salida Esperada

```
=== GENERANDO XML ===
✅ XML generado exitosamente
Tamaño: 12345 bytes
📄 XML sin firma guardado: invoice_unsigned.xml

=== FIRMANDO XML ===
✅ XML firmado exitosamente
Tamaño: 15678 bytes
📄 XML firmado guardado: invoice_signed.xml

=== ENVIANDO A DIAN ===
URL: https://vpfe-hab.dian.gov.co/WcfDianCustomerServices.svc
TestSetId: e6784f41-2aba-4ed3-bcb6-d045ab217e72

🔌 Validando conexión a DIAN...
✅ Conexión exitosa

📤 Enviando factura a DIAN...

=== RESPUESTA DIAN ===
✅ FACTURA ACEPTADA
Código: 00
Mensaje: Procesado correctamente
CUFE: 9042ac191bbf94bb8f224554f092feed...
Fecha respuesta: 2026-01-01 15:30:00

📄 Respuesta DIAN guardada: dian_response.xml

=== PROCESO COMPLETADO ===
```

## ⚠️ Posibles Errores

### Error: "certificate.pem not found"
**Solución:** Verifica que el certificado esté en `../certificates/certificate.pem`

### Error: "DIAN no disponible"
**Solución:** Verifica tu conexión a internet y que el endpoint de habilitación esté activo

### Error: "FACTURA RECHAZADA - Código 99"
**Posibles causas:**
- SoftwareID no coincide con el registrado
- CUFE mal calculado
- Certificado no válido
- Rango de numeración incorrecto

### Error: "Firma digital inválida"
**Solución:** Verifica que el certificado PEM contenga tanto el certificado como la clave privada

## 📁 Archivos Generados

Después de ejecutar, encontrarás:

- `invoice_unsigned.xml` - XML sin firma (debug)
- `invoice_signed.xml` - XML firmado listo para envío
- `dian_response.xml` - Respuesta completa de DIAN

## 🔍 Debugging

Para ver el contenido de los XMLs generados:

```bash
# Ver XML sin firma
cat invoice_unsigned.xml

# Ver XML firmado
cat invoice_signed.xml

# Ver respuesta DIAN
cat dian_response.xml
```

## 📝 Notas Importantes

1. **TestSetId:** Cada factura de prueba debe usar el TestSetId proporcionado por DIAN
2. **Numeración:** Usa el rango autorizado (990000000-995000000)
3. **Cuota:** Tienes 50 facturas de prueba disponibles
4. **Ambiente:** Este código usa el ambiente de habilitación (test), NO producción

## 🎯 Próximos Pasos

Una vez que la prueba funcione correctamente:

1. Integrar el código de envío en `api-dian`
2. Crear endpoint `/api/invoices/:id/send`
3. Actualizar estado de factura en BD según respuesta DIAN
4. Implementar manejo de errores y reintentos
