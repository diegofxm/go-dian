package extensions

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/xml"
	"fmt"
)

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

type IDType struct {
	Value            string `xml:",chardata"`
	SchemeID         string `xml:"schemeID,attr,omitempty"`
	SchemeName       string `xml:"schemeName,attr,omitempty"`
	SchemeAgencyID   string `xml:"schemeAgencyID,attr,omitempty"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr,omitempty"`
}

type SoftwareSecurityCode struct {
	Value            string `xml:",chardata"`
	SchemeAgencyID   string `xml:"schemeAgencyID,attr"`
	SchemeAgencyName string `xml:"schemeAgencyName,attr"`
}

type AuthorizationProvider struct {
	AuthorizationProviderID IDType `xml:"sts:AuthorizationProviderID"`
}

func GenerateSoftwareSecurityCode(softwareID, pin string) string {
	// SoftwareSecurityCode = SHA-384(SoftwareID + PIN)
	data := softwareID + pin
	hash := sha512.Sum384([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func GenerateQRCode(nit, invoiceID, uuid string) string {
	data := fmt.Sprintf("%s%s%s", nit, invoiceID, uuid)
	hash := sha256.Sum256([]byte(data))
	documentKey := base64.URLEncoding.EncodeToString(hash[:])
	return fmt.Sprintf("https://catalogo-vpfe.dian.gov.co/document/searchqr?documentkey=%s", documentKey)
}
