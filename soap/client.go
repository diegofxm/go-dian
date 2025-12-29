package soap

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	URL        string
	Timeout    time.Duration
	HTTPClient *http.Client
}

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	Content interface{} `xml:",innerxml"`
}

type SendBillSync struct {
	XMLName     xml.Name `xml:"http://wcf.dian.colombia SendBillSync"`
	FileName    string   `xml:"fileName"`
	ContentFile string   `xml:"contentFile"`
}

type SendBillSyncResponse struct {
	XMLName xml.Name `xml:"SendBillSyncResponse"`
	Result  string   `xml:"SendBillSyncResult"`
}

type Environment string

const (
	EnvironmentTest       Environment = "test"
	EnvironmentProduction Environment = "production"
)

var Endpoints = map[Environment]string{
	EnvironmentTest:       "https://vpfe-hab.dian.gov.co/WcfDianCustomerServices.svc",
	EnvironmentProduction: "https://vpfe.dian.gov.co/WcfDianCustomerServices.svc",
}

func NewClient(environment Environment) *Client {
	return &Client{
		URL:     Endpoints[environment],
		Timeout: 30 * time.Second,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SendInvoice(fileName string, signedXML []byte) (*Response, error) {
	contentFile := base64.StdEncoding.EncodeToString(signedXML)

	request := SendBillSync{
		FileName:    fileName,
		ContentFile: contentFile,
	}

	envelope := Envelope{
		Body: Body{
			Content: request,
		},
	}

	soapXML, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando SOAP: %w", err)
	}

	soapMessage := []byte(xml.Header + string(soapXML))

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(soapMessage))
	if err != nil {
		return nil, fmt.Errorf("error creando petición HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://wcf.dian.colombia/IWcfDianCustomerServices/SendBillSync")
	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(soapMessage)))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando petición a DIAN: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DIAN retornó código %d: %s", resp.StatusCode, string(body))
	}

	var responseEnvelope Envelope
	if err := xml.Unmarshal(body, &responseEnvelope); err != nil {
		return nil, fmt.Errorf("error parseando respuesta SOAP: %w", err)
	}

	dianResponse, err := parseResponse([]byte(responseEnvelope.Body.Content.(string)))
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta DIAN: %w", err)
	}

	return dianResponse, nil
}

func (c *Client) SendTestSet(testSetID string, signedXML []byte) (*Response, error) {
	fileName := fmt.Sprintf("test_%s.xml", testSetID)
	return c.SendInvoice(fileName, signedXML)
}

func (c *Client) ValidateConnection() error {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return fmt.Errorf("error creando petición: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error conectando a DIAN: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("DIAN no disponible (código %d)", resp.StatusCode)
	}

	return nil
}
