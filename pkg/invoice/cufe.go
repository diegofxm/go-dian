package invoice

import (
	"fmt"

	"github.com/diegofxm/go-dian/internal/hash"
)

func CalculateCUFE(inv *Invoice, nit string, technicalKey string, environment string) (string, error) {
	if len(inv.TaxTotal) == 0 {
		return "", fmt.Errorf("la factura debe tener al menos un TaxTotal")
	}

	invoiceNumber := inv.ID
	issueDate := inv.IssueDate
	issueTime := inv.IssueTime
	amount := fmt.Sprintf("%.2f", inv.LegalMonetaryTotal.LineExtensionAmount.Value)
	taxAmount := fmt.Sprintf("%.2f", inv.TaxTotal[0].TaxAmount.Value)
	totalAmount := fmt.Sprintf("%.2f", inv.LegalMonetaryTotal.PayableAmount.Value)
	customerNIT := inv.AccountingCustomerParty.Party.PartyTaxScheme.CompanyID.Value

	cufe := hash.CalculateCUFE(
		invoiceNumber,
		issueDate,
		issueTime,
		amount,
		taxAmount,
		totalAmount,
		nit,
		customerNIT,
		technicalKey,
		environment,
	)

	return cufe, nil
}
