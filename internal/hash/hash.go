package hash

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func CalculateSHA384(data string) string {
	hash := sha512.New384()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func CalculateCUFE(invoiceNumber, issueDate, issueTime, amount, taxAmount, totalAmount, nit, customerNIT, technicalKey, environment string) string {
	// Formatear fecha: remover guiones (YYYY-MM-DD -> YYYYMMDD)
	formattedDate := ""
	for _, c := range issueDate {
		if c != '-' {
			formattedDate += string(c)
		}
	}

	// Formatear hora: remover : y zona horaria (HH:MM:SS-05:00 -> HHMMSS)
	formattedTime := ""
	for _, c := range issueTime {
		if c >= '0' && c <= '9' {
			formattedTime += string(c)
		}
		if c == '-' || c == '+' {
			break
		}
	}
	if len(formattedTime) > 6 {
		formattedTime = formattedTime[:6]
	}

	// Construir cadena CUFE según especificación DIAN
	// NumFac + FecFac + HorFac + ValFac + CodImp1 + ValImp1 + CodImp2 + ValImp2 + CodImp3 + ValImp3 + ValTot + NitOFE + NumAdq + ClTec + TipoAmbiente
	cufeData := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s",
		invoiceNumber,
		formattedDate,
		formattedTime,
		amount,
		"01", // CodImp1 (IVA)
		taxAmount,
		"01", // CodImp2 (IVA adicional o 0)
		"0.00",
		"01", // CodImp3
		"0.00",
		totalAmount,
		nit,
		customerNIT,
		technicalKey,
		environment,
	)

	return CalculateSHA384(cufeData)
}
