package soap

import (
	"encoding/xml"
	"fmt"
)

// EnvelopeBuilder construye envelopes SOAP 1.2
type EnvelopeBuilder struct {
	namespace string
	wcfNS     string
}

// NewEnvelopeBuilder crea un nuevo builder
func NewEnvelopeBuilder() *EnvelopeBuilder {
	return &EnvelopeBuilder{
		namespace: "http://www.w3.org/2003/05/soap-envelope",
		wcfNS:     "http://wcf.dian.colombia",
	}
}

// BuildSendBillSync construye un envelope para SendBillSync
func (eb *EnvelopeBuilder) BuildSendBillSync(fileName, contentFile string, wsSecurityHeader string) ([]byte, error) {
	envelope := Envelope{
		SoapEnv: eb.namespace,
		Wcf:     eb.wcfNS,
		Body: Body{
			SendBillSync: SendBillSync{
				FileName:    fileName,
				ContentFile: contentFile,
			},
		},
	}

	// Agregar WS-Security Header si está presente
	if wsSecurityHeader != "" {
		envelope.Header = &Header{
			Security: wsSecurityHeader,
		}
	}

	// Serializar a XML
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando SOAP: %w", err)
	}

	// Agregar declaración XML
	result := []byte(xml.Header)
	result = append(result, xmlData...)

	return result, nil
}
