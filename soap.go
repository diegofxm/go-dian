package dian

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SOAPClient representa un cliente SOAP para DIAN
type SOAPClient struct {
	URL        string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// SOAPEnvelope representa el sobre SOAP
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SOAPBody `xml:"Body"`
}

// SOAPBody representa el cuerpo del mensaje SOAP
type SOAPBody struct {
	Content interface{} `xml:",innerxml"`
}

// SendBillSync representa la petición de envío de factura
type SendBillSync struct {
	XMLName     xml.Name `xml:"http://wcf.dian.colombia SendBillSync"`
	FileName    string   `xml:"fileName"`
	ContentFile string   `xml:"contentFile"`
}

// SendBillSyncResponse representa la respuesta de DIAN
type SendBillSyncResponse struct {
	XMLName xml.Name `xml:"SendBillSyncResponse"`
	Result  string   `xml:"SendBillSyncResult"`
}

// DianEndpoints contiene los endpoints de DIAN
var DianEndpoints = map[Environment]string{
	EnvironmentTest:       "https://vpfe-hab.dian.gov.co/WcfDianCustomerServices.svc",
	EnvironmentProduction: "https://vpfe.dian.gov.co/WcfDianCustomerServices.svc",
}

// NewSOAPClient crea un nuevo cliente SOAP
func NewSOAPClient(environment Environment) *SOAPClient {
	return &SOAPClient{
		URL:     DianEndpoints[environment],
		Timeout: 30 * time.Second,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendInvoiceToDAIN envía una factura firmada a DIAN
func (s *SOAPClient) SendInvoiceToDAIN(fileName string, signedXML []byte) (*DianResponse, error) {
	// Codificar XML en base64
	contentFile := encodeBase64(signedXML)

	// Crear petición SOAP
	request := SendBillSync{
		FileName:    fileName,
		ContentFile: contentFile,
	}

	envelope := SOAPEnvelope{
		Body: SOAPBody{
			Content: request,
		},
	}

	// Serializar SOAP
	soapXML, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando SOAP: %w", err)
	}

	// Agregar declaración XML
	soapMessage := []byte(xml.Header + string(soapXML))

	// Crear petición HTTP
	req, err := http.NewRequest("POST", s.URL, bytes.NewBuffer(soapMessage))
	if err != nil {
		return nil, fmt.Errorf("error creando petición HTTP: %w", err)
	}

	// Headers correctos para SOAP 1.1
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://wcf.dian.colombia/IWcfDianCustomerServices/SendBillSync")
	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(soapMessage)))

	// Enviar petición
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando petición a DIAN: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Verificar código de estado HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DIAN retornó código %d: %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta SOAP
	var responseEnvelope SOAPEnvelope
	if err := xml.Unmarshal(body, &responseEnvelope); err != nil {
		return nil, fmt.Errorf("error parseando respuesta SOAP: %w", err)
	}

	// Parsear resultado
	dianResponse, err := parseDianResponse([]byte(responseEnvelope.Body.Content.(string)))
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta DIAN: %w", err)
	}

	return dianResponse, nil
}

// DianResponse representa la respuesta de DIAN
type DianResponse struct {
	Success       bool
	StatusCode    string
	StatusMessage string
	CUFE          string
	Errors        []string
	Warnings      []string
	ResponseDate  time.Time
	RawResponse   string
}

// parseDianResponse parsea la respuesta XML de DIAN
func parseDianResponse(responseXML []byte) (*DianResponse, error) {
	// Estructura simplificada de respuesta DIAN
	type ApplicationResponse struct {
		StatusCode       string `xml:"StatusCode"`
		StatusMessage    string `xml:"StatusMessage"`
		DocumentResponse struct {
			Response struct {
				ResponseCode string `xml:"ResponseCode"`
			} `xml:"Response"`
		} `xml:"DocumentResponse"`
	}

	var appResponse ApplicationResponse
	if err := xml.Unmarshal(responseXML, &appResponse); err != nil {
		// Si no se puede parsear, retornar respuesta genérica
		return &DianResponse{
			Success:       false,
			StatusMessage: "Error parseando respuesta de DIAN",
			RawResponse:   string(responseXML),
			ResponseDate:  time.Now(),
		}, nil
	}

	success := appResponse.StatusCode == "00" || appResponse.DocumentResponse.Response.ResponseCode == "00"

	return &DianResponse{
		Success:       success,
		StatusCode:    appResponse.StatusCode,
		StatusMessage: appResponse.StatusMessage,
		RawResponse:   string(responseXML),
		ResponseDate:  time.Now(),
	}, nil
}

// GetStatus consulta el estado de un documento en DIAN
func (s *SOAPClient) GetStatus(trackingID string) (*DianResponse, error) {
	// TODO: Implementar consulta de estado
	return &DianResponse{
		Success:       true,
		StatusMessage: "Consulta de estado no implementada",
		ResponseDate:  time.Now(),
	}, nil
}

// SendTestSet envía un set de pruebas a DIAN (habilitación)
func (s *SOAPClient) SendTestSet(testSetID string, signedXML []byte) (*DianResponse, error) {
	// Similar a SendInvoiceToDAIN pero con endpoint de habilitación
	fileName := fmt.Sprintf("test_%s.xml", testSetID)
	return s.SendInvoiceToDAIN(fileName, signedXML)
}

// encodeBase64 codifica datos en base64
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ValidateConnection valida la conexión con DIAN
func (s *SOAPClient) ValidateConnection() error {
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		return fmt.Errorf("error creando petición: %w", err)
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error conectando a DIAN: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("DIAN no disponible (código %d)", resp.StatusCode)
	}

	return nil
}
