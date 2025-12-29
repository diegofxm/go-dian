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
	cufeData := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
		invoiceNumber,
		issueDate,
		issueTime,
		amount,
		"01",
		taxAmount,
		"04",
		totalAmount,
		nit,
		customerNIT,
	)

	if environment == "2" {
		cufeData += technicalKey
	}

	return CalculateSHA384(cufeData)
}
