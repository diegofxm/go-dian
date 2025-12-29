package dian

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	xmlutil "github.com/diegofxm/go-dian/internal/xml"
	"github.com/diegofxm/go-dian/invoice"
	"github.com/diegofxm/go-dian/signature"
	"github.com/diegofxm/go-dian/soap"
)

// Client representa el cliente para interactuar con DIAN
type Client struct {
	Config      Config
	soapClient  *soap.Client
	certManager *signature.CertificateManager
}

// NewClient crea una nueva instancia del cliente DIAN
func NewClient(config Config) (*Client, error) {
	if config.NIT == "" {
		return nil, ErrInvalidNIT
	}

	if config.Certificate.Path == "" && config.Certificate.CertPEM == "" {
		return nil, ErrMissingCertificate
	}

	var certManager *signature.CertificateManager
	var err error

	if config.Certificate.Path != "" {
		certManager, err = signature.NewCertificateManager(config.Certificate.Path, config.Certificate.Password)
		if err != nil {
			return nil, fmt.Errorf("error cargando certificado: %w", err)
		}
	}

	var soapEnv soap.Environment
	if config.Environment == EnvironmentProduction {
		soapEnv = soap.EnvironmentProduction
	} else {
		soapEnv = soap.EnvironmentTest
	}

	return &Client{
		Config:      config,
		soapClient:  soap.NewClient(soapEnv),
		certManager: certManager,
	}, nil
}

// GenerateInvoiceXML genera el XML de Invoice con DianExtensions (sin firmar)
func (c *Client) GenerateInvoiceXML(inv *invoice.Invoice) ([]byte, error) {
	if err := inv.Validate(); err != nil {
		return nil, fmt.Errorf("factura inválida: %w", err)
	}

	// Calcular CUFE
	cufe, err := invoice.CalculateCUFE(inv, c.Config.NIT, string(c.Config.Environment))
	if err != nil {
		return nil, fmt.Errorf("error calculando CUFE: %w", err)
	}
	inv.UUID.Value = cufe

	// Generar XML usando el generador modular
	genConfig := invoice.GeneratorConfig{
		NIT:                  c.Config.NIT,
		SoftwareID:           c.Config.SoftwareID,
		InvoiceAuthorization: c.Config.InvoiceAuthorization,
		AuthStartDate:        c.Config.AuthStartDate,
		AuthEndDate:          c.Config.AuthEndDate,
		InvoicePrefix:        c.Config.InvoicePrefix,
		AuthFrom:             c.Config.AuthFrom,
		AuthTo:               c.Config.AuthTo,
	}

	return invoice.GenerateXML(inv, genConfig)
}

// CalculateCUFE calcula el Código Único de Factura Electrónica
func (c *Client) CalculateCUFE(inv *invoice.Invoice) (string, error) {
	return invoice.CalculateCUFE(inv, c.Config.NIT, string(c.Config.Environment))
}

// SendInvoice envía la factura a DIAN
func (c *Client) SendInvoice(inv *invoice.Invoice) (*InvoiceResponse, error) {
	// 1. Generar XML sin firmar
	xmlData, err := c.GenerateInvoiceXML(inv)
	if err != nil {
		return nil, err
	}

	// 2. Firmar XML
	signedXML, err := c.SignXML(xmlData)
	if err != nil {
		return nil, fmt.Errorf("error firmando XML: %w", err)
	}

	// 3. Enviar a DIAN vía SOAP
	fileName := fmt.Sprintf("%s%s.xml", c.Config.InvoicePrefix, inv.ID)
	response, err := c.soapClient.SendInvoice(fileName, signedXML)
	if err != nil {
		return nil, fmt.Errorf("error enviando a DIAN: %w", err)
	}

	return &InvoiceResponse{
		Success:      response.Success,
		Message:      response.StatusMessage,
		CUFE:         response.CUFE,
		Errors:       response.Errors,
		ResponseDate: response.ResponseDate,
	}, nil
}

// SignXML firma cualquier XML con el certificado digital
func (c *Client) SignXML(xmlData []byte) ([]byte, error) {
	if c.certManager == nil {
		return nil, ErrMissingCertificate
	}

	// Crear firma XMLDSig
	signatureXML, err := signature.SignXMLDocument(xmlData, c.certManager.GetCertificate(), c.certManager.GetPrivateKey())
	if err != nil {
		return nil, fmt.Errorf("error generando firma: %w", err)
	}

	// Insertar firma en UBLExtensions
	signedXML, err := xmlutil.InsertSignature(xmlData, signatureXML)
	if err != nil {
		return nil, fmt.Errorf("error insertando firma: %w", err)
	}

	return signedXML, nil
}

// InvoiceResponse representa la respuesta de DIAN
type InvoiceResponse struct {
	Success      bool
	Message      string
	CUFE         string
	Errors       []string
	ResponseDate time.Time
}

// ValidateNIT valida el formato de un NIT colombiano
func ValidateNIT(nit string) error {
	nit = strings.ReplaceAll(nit, ".", "")
	nit = strings.ReplaceAll(nit, "-", "")

	if len(nit) < 9 || len(nit) > 10 {
		return fmt.Errorf("NIT debe tener entre 9 y 10 dígitos")
	}

	matched, _ := regexp.MatchString(`^\d+$`, nit)
	if !matched {
		return fmt.Errorf("NIT debe contener solo números")
	}

	return nil
}
