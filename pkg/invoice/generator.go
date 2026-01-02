package invoice

import (
	"encoding/xml"
	"fmt"

	"github.com/diegofxm/go-dian/pkg/extensions"
)

func GenerateXML(inv *Invoice, config GeneratorConfig) ([]byte, error) {
	if err := inv.Validate(); err != nil {
		return nil, fmt.Errorf("factura inválida: %w", err)
	}

	inv.UBLExtensions = buildExtensions(inv, config)

	invoiceXML, err := xml.MarshalIndent(inv, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generando XML: %w", err)
	}

	return []byte(xml.Header + string(invoiceXML)), nil
}

type GeneratorConfig struct {
	NIT                  string
	SoftwareID           string
	PIN                  string
	InvoiceAuthorization string
	AuthStartDate        string
	AuthEndDate          string
	InvoicePrefix        string
	AuthFrom             string
	AuthTo               string
}

func buildExtensions(inv *Invoice, config GeneratorConfig) *UBLExtensions {
	dianExt := extensions.DianExtensions{
		InvoiceControl: extensions.InvoiceControl{
			InvoiceAuthorization: config.InvoiceAuthorization,
			AuthorizationPeriod: extensions.AuthorizationPeriod{
				StartDate: config.AuthStartDate,
				EndDate:   config.AuthEndDate,
			},
			AuthorizedInvoices: extensions.AuthorizedInvoices{
				Prefix: config.InvoicePrefix,
				From:   config.AuthFrom,
				To:     config.AuthTo,
			},
		},
		InvoiceSource: extensions.InvoiceSource{
			IdentificationCode: extensions.IdentificationCode{
				Value:          "CO",
				ListAgencyID:   "6",
				ListAgencyName: "United Nations Economic Commission for Europe",
				ListSchemeURI:  "urn:oasis:names:specification:ubl:codelist:gc:CountryIdentificationCode-2.1",
			},
		},
		SoftwareProvider: extensions.SoftwareProvider{
			ProviderID: extensions.IDType{
				Value:            config.NIT,
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				SchemeName:       "31",
				SchemeID:         "4",
			},
			SoftwareID: extensions.IDType{
				Value:            config.SoftwareID,
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
			},
		},
		SoftwareSecurityCode: extensions.SoftwareSecurityCode{
			Value:            extensions.GenerateSoftwareSecurityCode(config.SoftwareID, config.PIN),
			SchemeAgencyID:   "195",
			SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
		},
		AuthorizationProvider: extensions.AuthorizationProvider{
			AuthorizationProviderID: extensions.IDType{
				Value:            "800197268",
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				SchemeName:       "31",
				SchemeID:         "4",
			},
		},
		QRCode: extensions.GenerateQRCode(config.NIT, inv.ID, inv.UUID.Value),
	}

	dianExtXML, _ := xml.Marshal(dianExt)

	return &UBLExtensions{
		UBLExtension: []UBLExtension{
			{
				ExtensionContent: ExtensionContent{
					Content: string(dianExtXML),
				},
			},
		},
	}
}
