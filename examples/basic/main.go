package main

import (
	"fmt"
	"log"

	"github.com/diegofxm/go-dian"
)

func main() {
	client, err := dian.NewClient(dian.Config{
		NIT:         "6382356",
		Environment: dian.EnvironmentTest,
		SoftwareID:  "74fdde3e-8dc0-4d90-8515-4f9c19634999",
		Certificate: dian.Certificate{
			Path:     "../certificate.p12",
			Password: "ScBPigJrvrKjmbqg",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	invoice := dian.NewInvoice("BEC496329154")

	invoice.AccountingSupplierParty = dian.AccountingSupplierParty{
		AdditionalAccountID: "1",
		Party: dian.Party{
			IndustryClassificationCode: "6120",
			PartyIdentification: dian.PartyIdentification{
				ID: dian.IDType{
					Value:            "830122566",
					SchemeID:         "1",
					SchemeName:       "31",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
			PartyName: dian.PartyName{
				Name: "COLOMBIA TELECOMUNICACIONES S.A. E.S.P. BIC",
			},
			PartyTaxScheme: dian.PartyTaxScheme{
				RegistrationName: "COLOMBIA TELECOMUNICACIONES S.A. E.S.P. BIC",
				CompanyID: dian.IDType{
					Value:            "830122566",
					SchemeID:         "1",
					SchemeName:       "31",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
				TaxLevelCode: "O-13;O-15;O-23",
				TaxScheme: dian.TaxScheme{
					ID:   "01",
					Name: "IVA",
				},
			},
			PartyLegalEntity: dian.PartyLegalEntity{
				RegistrationName: "COLOMBIA TELECOMUNICACIONES S.A. E.S.P. BIC",
				CompanyID: dian.IDType{
					Value:            "830122566",
					SchemeID:         "1",
					SchemeName:       "31",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
			Contact: &dian.Contact{
				Telephone:      "018000930930",
				ElectronicMail: "contacto@empresa.com",
			},
		},
	}

	invoice.AccountingCustomerParty = dian.AccountingCustomerParty{
		AdditionalAccountID: "2",
		Party: dian.Party{
			PartyIdentification: dian.PartyIdentification{
				ID: dian.IDType{
					Value:            "6382356",
					SchemeID:         "13",
					SchemeName:       "13",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
			PartyName: dian.PartyName{
				Name: "DIEGO FERNANDO MONTOYA",
			},
			PartyTaxScheme: dian.PartyTaxScheme{
				RegistrationName: "DIEGO FERNANDO MONTOYA",
				CompanyID: dian.IDType{
					Value:            "6382356",
					SchemeID:         "13",
					SchemeName:       "13",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
				TaxLevelCode: "R-99-PN",
				TaxScheme: dian.TaxScheme{
					ID:   "01",
					Name: "IVA",
				},
			},
			PartyLegalEntity: dian.PartyLegalEntity{
				RegistrationName: "DIEGO FERNANDO MONTOYA",
				CompanyID: dian.IDType{
					Value: "6382356",
				},
			},
			Contact: &dian.Contact{
				Name:           "DIEGO FERNANDO MONTOYA",
				Telephone:      "3186708084",
				ElectronicMail: "cliente@correo.com",
			},
		},
	}

	invoice.AddLine(dian.InvoiceLine{
		ID: "1",
		InvoicedQuantity: dian.Quantity{
			Value:    1,
			UnitCode: "94",
		},
		LineExtensionAmount: dian.AmountType{
			Value:      23950,
			CurrencyID: "COP",
		},
		Item: dian.Item{
			Description: "Servicio de telecomunicaciones",
			BrandName:   "N/A",
			ModelName:   "N/A",
			SellersItemIdentification: &dian.ItemIdentification{
				ID: dian.IDType{Value: "1616"},
			},
		},
		Price: dian.Price{
			PriceAmount: dian.AmountType{
				Value:      23950,
				CurrencyID: "COP",
			},
			BaseQuantity: dian.Quantity{
				Value:    1,
				UnitCode: "94",
			},
		},
		TaxTotal: []dian.TaxTotal{
			{
				TaxAmount: dian.AmountType{
					Value:      4550.50,
					CurrencyID: "COP",
				},
				TaxSubtotal: []dian.TaxSubtotal{
					{
						TaxableAmount: dian.AmountType{
							Value:      23950,
							CurrencyID: "COP",
						},
						TaxAmount: dian.AmountType{
							Value:      4550.50,
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

	invoice.CalculateTotals()

	// Agregar totales de impuestos a nivel de factura
	invoice.TaxTotal = []dian.TaxTotal{
		{
			TaxAmount: dian.AmountType{
				Value:      4550.50,
				CurrencyID: "COP",
			},
			TaxSubtotal: []dian.TaxSubtotal{
				{
					TaxableAmount: dian.AmountType{
						Value:      23950,
						CurrencyID: "COP",
					},
					TaxAmount: dian.AmountType{
						Value:      4550.50,
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
	}

	xmlData, err := client.GenerateInvoiceXML(invoice)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("XML generado:")
	fmt.Println(string(xmlData))

	response, err := client.SendInvoice(invoice)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nRespuesta DIAN: %+v\n", response)
}
