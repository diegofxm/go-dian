package main

import (
	"fmt"
	"log"

	"github.com/diegofxm/go-dian/pkg/common"
	"github.com/diegofxm/go-dian/pkg/dian"
	"github.com/diegofxm/go-dian/pkg/invoice"
)

func main() {
	// Configurar cliente DIAN
	client, err := dian.NewClient(dian.Config{
		NIT:         "830122566",
		Environment: dian.EnvironmentTest,
		SoftwareID:  "74fdde3e-8dc0-4d90-8515-4f9c19634999",
		Certificate: dian.Certificate{
			PEMPath: "../certificates/certificate.pem",
		},
		InvoiceAuthorization: "18760000001",
		AuthStartDate:        "2019-01-19",
		AuthEndDate:          "2030-01-19",
		InvoicePrefix:        "SETP",
		AuthFrom:             "990000000",
		AuthTo:               "995000000",
	})
	if err != nil {
		log.Fatalf("error creando cliente: %v", err)
	}

	// Crear factura
	inv := invoice.NewInvoice("SETP990000001")
	inv.InvoiceTypeCode = "01"
	inv.DocumentCurrencyCode = invoice.DocumentCurrencyType{
		Value:          "COP",
		ListAgencyID:   "6",
		ListAgencyName: "United Nations Economic Commission for Europe",
	}

	// Configurar emisor
	inv.AccountingSupplierParty = invoice.AccountingSupplierParty{
		AdditionalAccountID: common.AdditionalAccountIDType{
			Value:      "1",
			SchemeName: "tipos de personas",
		},
		Party: common.Party{
			PartyName: []common.PartyName{
				{Name: "EMPRESA DE PRUEBA S.A.S"},
			},
			PhysicalLocation: &common.PhysicalLocation{
				Address: common.Address{
					ID:                   "11001",
					CityName:             "Bogotá",
					CountrySubentity:     "Bogotá",
					CountrySubentityCode: "11",
					Country: common.Country{
						IdentificationCode: "CO",
						Name:               "Colombia",
					},
				},
			},
			PartyTaxScheme: common.PartyTaxScheme{
				RegistrationName: "EMPRESA DE PRUEBA S.A.S",
				CompanyID: common.IDType{
					Value:            "830122566",
					SchemeID:         "31",
					SchemeName:       "NIT",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
				TaxLevelCode: common.TaxLevelCodeType{
					Value:    "O-13",
					ListName: "Responsabilidades",
				},
				TaxScheme: common.TaxScheme{
					ID:   "01",
					Name: "IVA",
				},
			},
			PartyLegalEntity: common.PartyLegalEntity{
				RegistrationName: "EMPRESA DE PRUEBA S.A.S",
				CompanyID: common.IDType{
					Value:            "830122566",
					SchemeID:         "31",
					SchemeName:       "NIT",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
			PartyIdentification: common.PartyIdentification{
				ID: common.IDType{
					Value:            "830122566",
					SchemeID:         "31",
					SchemeName:       "NIT",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
		},
	}

	// Configurar cliente
	inv.AccountingCustomerParty = invoice.AccountingCustomerParty{
		AdditionalAccountID: common.AdditionalAccountIDType{
			Value:      "1",
			SchemeName: "tipos de personas",
		},
		Party: common.Party{
			PartyName: []common.PartyName{
				{Name: "CLIENTE DE PRUEBA"},
			},
			PhysicalLocation: &common.PhysicalLocation{
				Address: common.Address{
					ID:                   "11001",
					CityName:             "Bogotá",
					CountrySubentity:     "Bogotá",
					CountrySubentityCode: "11",
					Country: common.Country{
						IdentificationCode: "CO",
						Name:               "Colombia",
					},
				},
			},
			PartyTaxScheme: common.PartyTaxScheme{
				RegistrationName: "CLIENTE DE PRUEBA",
				CompanyID: common.IDType{
					Value:            "900123456",
					SchemeID:         "31",
					SchemeName:       "NIT",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
				TaxLevelCode: common.TaxLevelCodeType{
					Value:    "O-13",
					ListName: "Responsabilidades",
				},
				TaxScheme: common.TaxScheme{
					ID:   "01",
					Name: "IVA",
				},
			},
			PartyLegalEntity: common.PartyLegalEntity{
				RegistrationName: "CLIENTE DE PRUEBA",
				CompanyID: common.IDType{
					Value:            "900123456",
					SchemeID:         "31",
					SchemeName:       "NIT",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
			PartyIdentification: common.PartyIdentification{
				ID: common.IDType{
					Value:            "900123456",
					SchemeID:         "31",
					SchemeName:       "NIT",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
		},
	}

	// Agregar línea de factura
	line := invoice.InvoiceLine{
		ID: "1",
		InvoicedQuantity: common.Quantity{
			Value:    1.0,
			UnitCode: "94",
		},
		LineExtensionAmount: common.AmountType{
			Value:      100000.00,
			CurrencyID: "COP",
		},
		Item: invoice.Item{
			Description: "Servicio de consultoría",
		},
		Price: invoice.Price{
			PriceAmount: common.AmountType{
				Value:      100000.00,
				CurrencyID: "COP",
			},
			BaseQuantity: common.Quantity{
				Value:    1.0,
				UnitCode: "94",
			},
		},
		TaxTotal: []common.TaxTotal{
			{
				TaxAmount: common.AmountType{
					Value:      19000.00,
					CurrencyID: "COP",
				},
				TaxSubtotal: []common.TaxSubtotal{
					{
						TaxableAmount: common.AmountType{
							Value:      100000.00,
							CurrencyID: "COP",
						},
						TaxAmount: common.AmountType{
							Value:      19000.00,
							CurrencyID: "COP",
						},
						TaxCategory: common.TaxCategory{
							Percent: 19.0,
							TaxScheme: common.TaxScheme{
								ID:   "01",
								Name: "IVA",
							},
						},
					},
				},
			},
		},
	}

	inv.AddLine(line)
	inv.CalculateTotals()

	// Generar XML
	xmlData, err := client.GenerateInvoiceXML(inv)
	if err != nil {
		log.Fatalf("error generando XML: %v", err)
	}

	fmt.Println("XML generado:")
	fmt.Println(string(xmlData))

	// Firmar XML
	signedXML, err := client.SignXML(xmlData)
	if err != nil {
		log.Fatalf("error firmando XML: %v", err)
	}

	fmt.Println("\nXML firmado exitosamente")
	fmt.Printf("Tamaño: %d bytes\n", len(signedXML))
}
