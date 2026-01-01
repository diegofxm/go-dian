#!/bin/bash
# Script para convertir certificado .p12 a formato PEM
# Uso: ./convert_p12.sh <archivo.p12> <password>

if [ "$#" -ne 2 ]; then
    echo "Uso: $0 <archivo.p12> <password>"
    exit 1
fi

P12_FILE="$1"
PASSWORD="$2"

if [ ! -f "$P12_FILE" ]; then
    echo "Error: Archivo $P12_FILE no encontrado"
    exit 1
fi

echo "Convirtiendo $P12_FILE a formato PEM..."

# Extraer certificado
openssl pkcs12 -in "$P12_FILE" -clcerts -nokeys -out certificate.pem -password "pass:$PASSWORD"
if [ $? -eq 0 ]; then
    echo "✅ Certificado extraído: certificate.pem"
else
    echo "❌ Error extrayendo certificado"
    exit 1
fi

# Extraer llave privada
openssl pkcs12 -in "$P12_FILE" -nocerts -nodes -out private_key.pem -password "pass:$PASSWORD"
if [ $? -eq 0 ]; then
    echo "✅ Llave privada extraída: private_key.pem"
else
    echo "❌ Error extrayendo llave privada"
    exit 1
fi

echo ""
echo "Conversión completada exitosamente"
echo "Archivos generados:"
echo "  - certificate.pem (certificado público)"
echo "  - private_key.pem (llave privada)"
