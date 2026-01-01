package soap

import "encoding/xml"

// Envelope representa el envelope SOAP 1.2
type Envelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	SoapEnv string   `xml:"xmlns:soapenv,attr"`
	Wcf     string   `xml:"xmlns:wcf,attr"`
	Header  *Header  `xml:"soapenv:Header,omitempty"`
	Body    Body     `xml:"soapenv:Body"`
}

// Header representa el header SOAP (contendrá WS-Security)
type Header struct {
	XMLName  xml.Name `xml:"soapenv:Header"`
	Security string   `xml:",innerxml"`
}

// Body representa el body SOAP
type Body struct {
	XMLName      xml.Name     `xml:"soapenv:Body"`
	SendBillSync SendBillSync `xml:"wcf:SendBillSync"`
}

// SendBillSync representa la operación de envío
type SendBillSync struct {
	XMLName     xml.Name `xml:"wcf:SendBillSync"`
	FileName    string   `xml:"wcf:fileName"`
	ContentFile string   `xml:"wcf:contentFile"`
}

// Response representa la respuesta de DIAN
type Response struct {
	IsValid       bool
	StatusCode    string
	StatusMessage string
	ErrorMessages []string
	CUFE          string
}

// ResponseEnvelope representa el envelope de respuesta SOAP
type ResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		SendBillSyncResponse struct {
			Result string `xml:"SendBillSyncResult"`
		} `xml:"SendBillSyncResponse"`
	} `xml:"Body"`
}
