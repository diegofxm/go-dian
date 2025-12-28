package dian

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
)

// DianExtensions contiene las extensiones específicas de DIAN
type DianExtensions struct {
	XMLName               xml.Name              `xml:"dian:gov:co:facturaelectronica:Structures-2-1 DianExtensions"`
	InvoiceControl        InvoiceControl        `xml:"sts:InvoiceControl"`
	InvoiceSource         InvoiceSource         `xml:"sts:InvoiceSource"`
	SoftwareProvider      SoftwareProvider      `xml:"sts:SoftwareProvider"`
	SoftwareSecurityCode  SoftwareSecurityCode  `xml:"sts:SoftwareSecurityCode"`
	AuthorizationProvider AuthorizationProvider `xml:"sts:AuthorizationProvider"`
	QRCode                string                `xml:"sts:QRCode"`
}

type InvoiceControl struct {
	InvoiceAuthorization string              `xml:"sts:InvoiceAuthorization"`
	AuthorizationPeriod  AuthorizationPeriod `xml:"sts:AuthorizationPeriod"`
	AuthorizedInvoices   AuthorizedInvoices  `xml:"sts:AuthorizedInvoices"`
}

type AuthorizationPeriod struct {
	StartDate string `xml:"cbc:StartDate"`
	EndDate   string `xml:"cbc:EndDate"`
}

type AuthorizedInvoices struct {
	Prefix string `xml:"sts:Prefix"`
	From   string `xml:"sts:From"`
	To     string `xml:"sts:To"`
}

type InvoiceSource struct {
	IdentificationCode IdentificationCode `xml:"cbc:IdentificationCode"`
}

type IdentificationCode struct {
	Value          string `xml:",chardata"`
	ListAgencyID   string `xml:"listAgencyID,attr"`
	ListAgencyName string `xml:"listAgencyName,attr"`
	ListSchemeURI  string `xml:"listSchemeURI,attr"`
}

type SoftwareProvider struct {
	ProviderID IDType `xml:"sts:ProviderID"`
	SoftwareID IDType `xml:"sts:SoftwareID"`
}

type SoftwareSecurityCode struct {
	Value            string `xml:",chardata"`
	SchemeAgencyID   string `xml:"schemeAgencyID,attr"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr"`
}

type AuthorizationProvider struct {
	AuthorizationProviderID IDType `xml:"sts:AuthorizationProviderID"`
}

type UBLExtensions struct {
	UBLExtension []UBLExtension `xml:"ext:UBLExtension"`
}

type UBLExtension struct {
	ExtensionContent ExtensionContent `xml:"ext:ExtensionContent"`
}

type ExtensionContent struct {
	Content interface{} `xml:",innerxml"`
}

// buildInvoiceExtensions construye las extensiones UBL para la factura
func (c *Client) buildInvoiceExtensions(invoice *Invoice) *UBLExtensions {
	dianExt := DianExtensions{
		InvoiceControl: InvoiceControl{
			InvoiceAuthorization: "18764090648904",
			AuthorizationPeriod: AuthorizationPeriod{
				StartDate: "2025-03-18",
				EndDate:   "2027-03-18",
			},
			AuthorizedInvoices: AuthorizedInvoices{
				Prefix: "BEC",
				From:   "450000001",
				To:     "500000000",
			},
		},
		InvoiceSource: InvoiceSource{
			IdentificationCode: IdentificationCode{
				Value:          "CO",
				ListAgencyID:   "6",
				ListAgencyName: "United Nations Economic Commission for Europe",
				ListSchemeURI:  "urn:oasis:names:specification:ubl:codelist:gc:CountryIdentificationCode-2.1",
			},
		},
		SoftwareProvider: SoftwareProvider{
			ProviderID: IDType{
				Value:            c.NIT,
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				SchemeName:       "31",
				SchemeID:         "1",
			},
			SoftwareID: IDType{
				Value:            c.SoftwareID,
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
			},
		},
		SoftwareSecurityCode: SoftwareSecurityCode{
			Value:            c.generateSoftwareSecurityCode(invoice),
			SchemeAgencyID:   "195",
			SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
		},
		AuthorizationProvider: AuthorizationProvider{
			AuthorizationProviderID: IDType{
				Value:            "800197268",
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				SchemeName:       "31",
				SchemeID:         "4",
			},
		},
		QRCode: c.generateQRCode(invoice),
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

// signInvoiceXML firma el XML de la factura e inserta la firma en UBLExtensions
func (c *Client) signInvoiceXML(invoiceXML []byte, cert interface{}, privateKey interface{}) ([]byte, error) {
	// Generar firma XMLDSig
	signatureXML, err := SignXMLDocument(invoiceXML, cert.(*x509.Certificate), privateKey.(*rsa.PrivateKey))
	if err != nil {
		return nil, err
	}

	// Insertar firma en UBLExtensions
	invoiceStr := string(invoiceXML)

	extensionsEnd := indexOf(invoiceStr, "</ext:UBLExtensions>")
	if extensionsEnd == -1 {
		return nil, fmt.Errorf("no se encontró UBLExtensions en la factura")
	}

	signatureExtension := fmt.Sprintf(`  <ext:UBLExtension>
    <ext:ExtensionContent>
%s
    </ext:ExtensionContent>
  </ext:UBLExtension>
`, string(signatureXML))

	result := invoiceStr[:extensionsEnd] + signatureExtension + invoiceStr[extensionsEnd:]
	return []byte(result), nil
}

func (c *Client) generateSoftwareSecurityCode(invoice *Invoice) string {
	data := fmt.Sprintf("%s%s%s", c.SoftwareID, c.NIT, invoice.ID)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func (c *Client) generateQRCode(invoice *Invoice) string {
	data := fmt.Sprintf("%s%s%s", c.NIT, invoice.ID, invoice.UUID)
	hash := sha256.Sum256([]byte(data))
	documentKey := base64.URLEncoding.EncodeToString(hash[:])
	return fmt.Sprintf("https://catalogo-vpfe.dian.gov.co/document/searchqr?documentkey=%s", documentKey)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
