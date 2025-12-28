package dian

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ValidateNIT valida el formato de un NIT colombiano
func ValidateNIT(nit string) error {
	nit = strings.ReplaceAll(nit, ".", "")
	nit = strings.ReplaceAll(nit, "-", "")

	if len(nit) < 9 || len(nit) > 10 {
		return fmt.Errorf("NIT debe tener entre 9 y 10 dígitos")
	}

	matched, _ := regexp.MatchString(`^\d+$`, nit)
	if !matched {
		return fmt.Errorf("NIT debe contener solo números")
	}

	return nil
}

// FormatCurrency formatea un valor monetario
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// FormatDate formatea una fecha para DIAN (YYYY-MM-DD)
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime formatea fecha y hora para DIAN (YYYY-MM-DDTHH:MM:SS-05:00)
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05-07:00")
}

// ParseDate parsea una fecha en formato DIAN
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// EncodeBase64 codifica datos en base64
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodifica datos en base64
func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// CalculateIVA calcula el IVA de un monto
func CalculateIVA(amount float64, percentage float64) float64 {
	return amount * (percentage / 100)
}

// CalculateTotalWithTax calcula el total incluyendo impuestos
func CalculateTotalWithTax(amount float64, taxPercentage float64) float64 {
	return amount + CalculateIVA(amount, taxPercentage)
}

// RoundToDecimals redondea un número a N decimales
func RoundToDecimals(value float64, decimals int) float64 {
	multiplier := float64(1)
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	return float64(int(value*multiplier+0.5)) / multiplier
}

// GenerateInvoiceNumber genera un número de factura con prefijo
func GenerateInvoiceNumber(prefix string, number int) string {
	return fmt.Sprintf("%s%d", prefix, number)
}

// ValidateEmail valida un email
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email inválido: %s", email)
	}
	return nil
}

// SanitizeString limpia caracteres especiales de una cadena
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\t", " ")

	// Reemplazar múltiples espacios por uno solo
	spaceRegex := regexp.MustCompile(`\s+`)
	s = spaceRegex.ReplaceAllString(s, " ")

	return s
}

// GetTaxSchemeCode retorna el código del esquema tributario
func GetTaxSchemeCode(taxType string) string {
	codes := map[string]string{
		"IVA": "01",
		"INC": "04",
		"ICA": "03",
	}

	if code, ok := codes[taxType]; ok {
		return code
	}
	return "01" // Default IVA
}

// GetDocumentTypeCode retorna el código del tipo de documento
func GetDocumentTypeCode(docType string) string {
	codes := map[string]string{
		"FACTURA":      "01",
		"NOTA_CREDITO": "91",
		"NOTA_DEBITO":  "92",
		"DOC_SOPORTE":  "05",
	}

	if code, ok := codes[docType]; ok {
		return code
	}
	return "01" // Default factura
}

// GetPersonTypeCode retorna el código del tipo de persona
func GetPersonTypeCode(personType string) string {
	codes := map[string]string{
		"JURIDICA": "1",
		"NATURAL":  "2",
	}

	if code, ok := codes[personType]; ok {
		return code
	}
	return "1" // Default jurídica
}

// GetIdentificationTypeCode retorna el código del tipo de identificación
func GetIdentificationTypeCode(idType string) string {
	codes := map[string]string{
		"NIT":       "31",
		"CC":        "13",
		"CE":        "22",
		"PASAPORTE": "41",
		"TI":        "12",
	}

	if code, ok := codes[idType]; ok {
		return code
	}
	return "31" // Default NIT
}
