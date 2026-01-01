package soap

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/diegofxm/go-dian/pkg/wssecurity"
)

// Environment representa el ambiente de DIAN
type Environment string

const (
	Test       Environment = "test"
	Production Environment = "production"
)

// Endpoints de DIAN
var Endpoints = map[Environment]string{
	Test:       "https://vpfe-hab.dian.gov.co/WcfDianCustomerServices.svc",
	Production: "https://vpfe.dian.gov.co/WcfDianCustomerServices.svc",
}

// Client es el cliente SOAP para DIAN
type Client struct {
	URL             string
	Environment     Environment
	Certificate     tls.Certificate
	HTTPClient      *http.Client
	EnvelopeBuilder *EnvelopeBuilder
	HeaderBuilder   *wssecurity.HeaderBuilder
}

// NewClient crea un nuevo cliente SOAP con certificado para mTLS
func NewClient(environment Environment, certPEMBlock, keyPEMBlock []byte) (*Client, error) {
	// Cargar certificado TLS
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, fmt.Errorf("error cargando certificado: %w", err)
	}

	// Configurar TLS con mTLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Crear cliente HTTP con TLS
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Crear HeaderBuilder para WS-Security
	headerBuilder, err := wssecurity.NewHeaderBuilder(cert)
	if err != nil {
		return nil, fmt.Errorf("error creando header builder: %w", err)
	}

	return &Client{
		URL:             Endpoints[environment],
		Environment:     environment,
		Certificate:     cert,
		HTTPClient:      httpClient,
		EnvelopeBuilder: NewEnvelopeBuilder(),
		HeaderBuilder:   headerBuilder,
	}, nil
}

// SendInvoice envía una factura a DIAN
func (c *Client) SendInvoice(fileName string, signedXML []byte) (*Response, error) {
	// 1. Codificar XML en base64
	contentFile := base64.StdEncoding.EncodeToString(signedXML)

	// 2. Generar WS-Security Header
	wsSecurityHeader, err := c.HeaderBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("error generando WS-Security header: %w", err)
	}

	// 3. Construir envelope SOAP con WS-Security
	soapMessage, err := c.EnvelopeBuilder.BuildSendBillSync(
		fileName,
		contentFile,
		wsSecurityHeader.ToXML(),
	)
	if err != nil {
		return nil, fmt.Errorf("error construyendo envelope SOAP: %w", err)
	}

	// Debug: imprimir mensaje SOAP
	fmt.Println("\n=== MENSAJE SOAP ENVIADO ===")
	fmt.Println(string(soapMessage))
	fmt.Println("============================\n")

	// 4. Crear petición HTTP
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(soapMessage))
	if err != nil {
		return nil, fmt.Errorf("error creando petición HTTP: %w", err)
	}

	// Headers SOAP 1.2
	req.Header.Set("Content-Type", "application/soap+xml;charset=UTF-8")
	req.Header.Set("SOAPAction", "http://wcf.dian.colombia/IWcfDianCustomerServices/SendBillSync")
	req.Header.Set("Accept", "application/soap+xml")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(soapMessage)))

	// 5. Enviar petición
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando petición a DIAN: %w", err)
	}
	defer resp.Body.Close()

	// 6. Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Debug: imprimir respuesta
	fmt.Println("\n=== RESPUESTA DIAN ===")
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Println(string(body))
	fmt.Println("======================\n")

	// 7. Verificar código de estado
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DIAN retornó código %d: %s", resp.StatusCode, string(body))
	}

	// 8. Parsear respuesta SOAP
	var responseEnvelope ResponseEnvelope
	if err := xml.Unmarshal(body, &responseEnvelope); err != nil {
		return nil, fmt.Errorf("error parseando respuesta SOAP: %w", err)
	}

	// 9. Decodificar respuesta de DIAN (viene en base64)
	responseData, err := base64.StdEncoding.DecodeString(responseEnvelope.Body.SendBillSyncResponse.Result)
	if err != nil {
		return nil, fmt.Errorf("error decodificando respuesta DIAN: %w", err)
	}

	// 10. Parsear respuesta DIAN
	dianResponse, err := parseResponse(responseData)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta DIAN: %w", err)
	}

	return dianResponse, nil
}

// parseResponse parsea la respuesta XML de DIAN
func parseResponse(data []byte) (*Response, error) {
	type DianResponse struct {
		XMLName       xml.Name `xml:"ApplicationResponse"`
		IsValid       string   `xml:"DocumentResponse>Response>ResponseCode"`
		StatusCode    string   `xml:"DocumentResponse>Response>Status>StatusReasonCode"`
		StatusMessage string   `xml:"DocumentResponse>Response>Status>StatusReason"`
		CUFE          string   `xml:"DocumentResponse>DocumentReference>UUID"`
	}

	var dianResp DianResponse
	if err := xml.Unmarshal(data, &dianResp); err != nil {
		return nil, err
	}

	return &Response{
		IsValid:       dianResp.IsValid == "00",
		StatusCode:    dianResp.StatusCode,
		StatusMessage: dianResp.StatusMessage,
		CUFE:          dianResp.CUFE,
	}, nil
}
