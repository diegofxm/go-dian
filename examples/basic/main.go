package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/diegofxm/go-dian/pkg/common"
	"github.com/diegofxm/go-dian/pkg/dian"
	"github.com/diegofxm/go-dian/pkg/invoice"
	"github.com/diegofxm/go-dian/pkg/soap"
)

func main() {
	// ========================================
	// DATOS REALES DE HABILITACIÓN DIAN
	// ========================================
	// TestSetId: e6784f41-2aba-4ed3-bcb6-d045ab217e72
	// SoftwareID: 23bf9eac-4dbe-4300-af06-541cc3efc7ca
	// Clave Técnica: fc8eac422eba16e22ffd8c6f94b3f40a6e38162c
	// PIN: 40125
	// Cuota: 50 facturas (30 FE, 10 ND, 10 NC)
	// ========================================

	// Configurar cliente DIAN con datos reales
	client, err := dian.NewClient(dian.Config{
		NIT:          "6382356", // NIT real de DIEGO FERNANDO MONTOYA VALLEJO
		Environment:  dian.EnvironmentTest,
		SoftwareID:   "23bf9eac-4dbe-4300-af06-541cc3efc7ca", // SoftwareID real
		TechnicalKey: "fc8eac422eba16e22ffd8c6f94b3f40a6e38162c",
		PIN:          "40125",
		Certificate: dian.Certificate{
			PEMPath: "../certificates/certificate.pem", // Ruta a tu certificado
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

	// Crear factura con número del rango autorizado
	inv := invoice.NewInvoice("SETP990000001") // Usar rango 990000000-995000000
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
				{Name: "DIEGO FERNANDO MONTOYA VALLEJO"},
			},
			PhysicalLocation: &common.PhysicalLocation{
				Address: common.Address{
					ID:                   "76520",
					CityName:             "Palmira",
					CountrySubentity:     "Valle del Cauca",
					CountrySubentityCode: "76",
					AddressLine: &common.AddressLine{
						Line: "CL 30 40 69",
					},
					Country: common.Country{
						IdentificationCode: "CO",
						Name:               "Colombia",
					},
				},
			},
			PartyTaxScheme: common.PartyTaxScheme{
				RegistrationName: "DIEGO FERNANDO MONTOYA VALLEJO",
				CompanyID: common.IDType{
					Value:            "6382356",
					SchemeID:         "13",
					SchemeName:       "13",
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
				RegistrationName: "DIEGO FERNANDO MONTOYA VALLEJO",
				CompanyID: common.IDType{
					Value:            "6382356",
					SchemeID:         "13",
					SchemeName:       "13",
					SchemeAgencyID:   "195",
					SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				},
			},
			PartyIdentification: common.PartyIdentification{
				ID: common.IDType{
					Value:            "6382356",
					SchemeID:         "13",
					SchemeName:       "13",
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
					AddressLine: &common.AddressLine{
						Line: "Calle 100 # 10-20",
					},
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
					SchemeID:         "3",
					SchemeName:       "31",
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
					SchemeID:         "3",
					SchemeName:       "31",
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

	// Agregar medio de pago (requerido por DIAN)
	inv.PaymentMeans = []common.PaymentMeans{
		{
			ID:               "1",
			PaymentMeansCode: "10", // 10 = Efectivo
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
	fmt.Println("\n=== GENERANDO XML ===")
	xmlData, err := client.GenerateInvoiceXML(inv)
	if err != nil {
		log.Fatalf("error generando XML: %v", err)
	}

	fmt.Println("✅ XML generado exitosamente")
	fmt.Printf("Tamaño: %d bytes\n", len(xmlData))

	// Guardar XML sin firma (opcional, para debug)
	if err := os.WriteFile("invoice_unsigned.xml", xmlData, 0644); err != nil {
		log.Printf("⚠️  No se pudo guardar XML sin firma: %v", err)
	} else {
		fmt.Println("📄 XML sin firma guardado: invoice_unsigned.xml")
	}

	// Firmar XML
	fmt.Println("\n=== FIRMANDO XML ===")
	signedXML, err := client.SignXML(xmlData)
	if err != nil {
		log.Fatalf("error firmando XML: %v", err)
	}

	fmt.Println("✅ XML firmado exitosamente")
	fmt.Printf("Tamaño: %d bytes\n", len(signedXML))

	// Guardar XML firmado
	if err := os.WriteFile("invoice_signed.xml", signedXML, 0644); err != nil {
		log.Printf("⚠️  No se pudo guardar XML firmado: %v", err)
	} else {
		fmt.Println("📄 XML firmado guardado: invoice_signed.xml")
	}

	// ========================================
	// ENVIAR A DIAN (SOAP)
	// ========================================
	fmt.Println("\n=== ENVIANDO A DIAN ===")
	fmt.Println("URL: https://vpfe-hab.dian.gov.co/WcfDianCustomerServices.svc")
	fmt.Println("TestSetId: e6784f41-2aba-4ed3-bcb6-d045ab217e72")

	// Crear cliente SOAP con certificado para autenticación TLS
	// DIAN requiere autenticación mediante certificado digital (mTLS)
	// El archivo PEM contiene tanto el certificado como la llave privada
	pemData, err := os.ReadFile("../certificates/certificate.pem")
	if err != nil {
		log.Fatalf("❌ Error leyendo certificado: %v", err)
	}

	// El mismo archivo PEM contiene certificado y llave privada
	soapClient, err := soap.NewClient(soap.Test, pemData, pemData)
	if err != nil {
		log.Fatalf("❌ Error creando cliente SOAP: %v", err)
	}

	// Crear ZIP con la factura firmada
	fmt.Println("\n📦 Creando archivo ZIP...")
	zipData, err := createZIP("SETP990000001.xml", signedXML)
	if err != nil {
		log.Fatalf("❌ Error creando ZIP: %v", err)
	}
	fmt.Printf("✅ ZIP creado (%d bytes)\n", len(zipData))

	// Guardar ZIP (opcional, para debug)
	if err := os.WriteFile("SETP990000001.zip", zipData, 0644); err != nil {
		log.Printf("⚠️  No se pudo guardar ZIP: %v", err)
	}

	// Enviar factura
	fmt.Println("\n📤 Enviando factura a DIAN...")
	fileName := "SETP990000001.zip"
	response, err := soapClient.SendInvoice(fileName, zipData)
	if err != nil {
		log.Fatalf("❌ Error enviando factura: %v", err)
	}

	// Mostrar respuesta
	fmt.Println("\n=== RESPUESTA DIAN ===")
	if response.IsValid {
		fmt.Println("✅ FACTURA ACEPTADA")
	} else {
		fmt.Println("❌ FACTURA RECHAZADA")
	}
	fmt.Printf("Código: %s\n", response.StatusCode)
	fmt.Printf("Mensaje: %s\n", response.StatusMessage)
	fmt.Printf("CUFE: %s\n", response.CUFE)

	if len(response.ErrorMessages) > 0 {
		fmt.Println("\n⚠️  ERRORES:")
		for i, errMsg := range response.ErrorMessages {
			fmt.Printf("  %d. %s\n", i+1, errMsg)
		}
	}

	fmt.Println("\n=== PROCESO COMPLETADO ===")
}

// createZIP crea un archivo ZIP con el contenido especificado
func createZIP(filename string, content []byte) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("error creando archivo en ZIP: %w", err)
	}

	if _, err := fileWriter.Write(content); err != nil {
		return nil, fmt.Errorf("error escribiendo contenido: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error cerrando ZIP: %w", err)
	}

	return buf.Bytes(), nil
}
