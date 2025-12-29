package dian

import (
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Client representa el cliente para interactuar con DIAN
type Client struct {
	NIT         string
	Certificate Certificate
	Environment Environment
	SoftwareID  string
	TestSetID   string

	// Datos de autorización DIAN (específicos por empresa)
	InvoiceAuthorization string
	AuthStartDate        string
	AuthEndDate          string
	InvoicePrefix        string
	AuthFrom             string
	AuthTo               string
}

// Environment define el ambiente de DIAN
type Environment string

const (
	EnvironmentProduction Environment = "production"
	EnvironmentTest       Environment = "test"
)

// NewClient crea una nueva instancia del cliente DIAN
func NewClient(config Config) (*Client, error) {
	if config.NIT == "" {
		return nil, fmt.Errorf("NIT es requerido")
	}

	return &Client{
		NIT:                  config.NIT,
		Certificate:          config.Certificate,
		Environment:          config.Environment,
		SoftwareID:           config.SoftwareID,
		TestSetID:            config.TestSetID,
		InvoiceAuthorization: config.InvoiceAuthorization,
		AuthStartDate:        config.AuthStartDate,
		AuthEndDate:          config.AuthEndDate,
		InvoicePrefix:        config.InvoicePrefix,
		AuthFrom:             config.AuthFrom,
		AuthTo:               config.AuthTo,
	}, nil
}

// Config contiene la configuración del cliente
type Config struct {
	NIT         string
	Certificate Certificate
	Environment Environment
	SoftwareID  string
	TestSetID   string

	// Datos de autorización DIAN (específicos por empresa)
	InvoiceAuthorization string // Número de autorización DIAN
	AuthStartDate        string // Fecha inicio autorización (YYYY-MM-DD)
	AuthEndDate          string // Fecha fin autorización (YYYY-MM-DD)
	InvoicePrefix        string // Prefijo de facturación
	AuthFrom             string // Consecutivo desde
	AuthTo               string // Consecutivo hasta
}

// Certificate representa el certificado digital
type Certificate struct {
	Path     string // Ruta a certificado P12 o PEM
	Password string // Contraseña del P12
	CertPEM  string // Certificado PEM como string (para BD)
	KeyPEM   string // Clave privada PEM como string (para BD)
}

// GenerateInvoiceXML genera el XML de Invoice con DianExtensions (sin firmar)
func (c *Client) GenerateInvoiceXML(invoice *Invoice) ([]byte, error) {
	if err := invoice.Validate(); err != nil {
		return nil, fmt.Errorf("factura inválida: %w", err)
	}

	// Calcular CUFE antes de generar XML
	cufe, err := c.CalculateCUFE(invoice)
	if err != nil {
		return nil, fmt.Errorf("error calculando CUFE: %w", err)
	}
	invoice.UUID.Value = cufe

	// Agregar DianExtensions a la factura
	invoice.UBLExtensions = c.buildInvoiceExtensions(invoice)

	// Generar XML de la factura
	invoiceXML, err := xml.MarshalIndent(invoice, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generando XML de factura: %w", err)
	}

	// Agregar declaración XML
	result := []byte(xml.Header + string(invoiceXML))
	return result, nil
}

// CalculateCUFE calcula el Código Único de Factura Electrónica
func (c *Client) CalculateCUFE(invoice *Invoice) (string, error) {
	// CUFE = SHA384(NumFac + FecFac + HorFac + ValFac + CodImp1 + ValImp1 + ... + NitOFE + NumAdq + ClTec + TipoAmbiente)
	data := fmt.Sprintf("%s%s%s%.2f%s%.2f%s%s%s%s",
		invoice.ID,
		invoice.IssueDate,
		invoice.IssueTime,
		invoice.LegalMonetaryTotal.PayableAmount.Value,
		"01", // IVA
		invoice.TaxTotal[0].TaxAmount.Value,
		c.NIT,
		invoice.AccountingCustomerParty.Party.PartyTaxScheme.CompanyID.Value,
		c.SoftwareID,
		string(c.Environment),
	)

	hash := sha512.Sum384([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}

// SendInvoice envía la factura a DIAN
func (c *Client) SendInvoice(invoice *Invoice) (*InvoiceResponse, error) {
	// 1. Generar XML sin firmar
	xmlData, err := c.GenerateInvoiceXML(invoice)
	if err != nil {
		return nil, err
	}

	// 2. Firmar XML
	signedXML, err := c.SignXML(xmlData)
	if err != nil {
		return nil, fmt.Errorf("error firmando XML: %w", err)
	}

	// 3. Enviar a DIAN vía SOAP
	response, err := c.sendSOAP(signedXML)
	if err != nil {
		return nil, fmt.Errorf("error enviando a DIAN: %w", err)
	}

	return response, nil
}

// SignXML firma cualquier XML con el certificado digital (método genérico reutilizable)
func (c *Client) SignXML(xmlData []byte) ([]byte, error) {
	// Cargar certificado (soporta PEM strings o archivo P12/PEM)
	var cert *x509.Certificate
	var privateKey *rsa.PrivateKey
	var err error

	if c.Certificate.CertPEM != "" && c.Certificate.KeyPEM != "" {
		cert, privateKey, err = LoadCertificateFromPEMStrings(c.Certificate.CertPEM, c.Certificate.KeyPEM)
	} else if c.Certificate.Path != "" {
		cert, privateKey, err = LoadCertificate(c.Certificate.Path, c.Certificate.Password)
	} else {
		return nil, fmt.Errorf("certificado no configurado")
	}

	if err != nil {
		return nil, fmt.Errorf("error cargando certificado: %w", err)
	}

	// Firmar XML e insertar firma en UBLExtensions
	signedXML, err := c.signInvoiceXML(xmlData, cert, privateKey)
	if err != nil {
		return nil, fmt.Errorf("error firmando XML: %w", err)
	}

	return signedXML, nil
}

// sendSOAP envía el XML firmado a DIAN vía SOAP
func (c *Client) sendSOAP(signedXML []byte) (*InvoiceResponse, error) {
	soapClient := NewSOAPClient(c.Environment)

	fileName := fmt.Sprintf("invoice_%d.xml", time.Now().Unix())

	dianResponse, err := soapClient.SendInvoiceToDAIN(fileName, signedXML)
	if err != nil {
		return nil, fmt.Errorf("error enviando a DIAN: %w", err)
	}

	return &InvoiceResponse{
		Success:      dianResponse.Success,
		Message:      dianResponse.StatusMessage,
		CUFE:         dianResponse.CUFE,
		Errors:       dianResponse.Errors,
		ResponseDate: dianResponse.ResponseDate,
	}, nil
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
