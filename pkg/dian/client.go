package dian

import (
	"fmt"
	"regexp"
	"strings"

	xmlutil "github.com/diegofxm/go-dian/internal/xml"
	"github.com/diegofxm/go-dian/pkg/invoice"
	"github.com/diegofxm/go-dian/pkg/signature"
)

// Client representa el cliente para interactuar con DIAN
type Client struct {
	Config      Config
	certManager *signature.CertificateManager
}

// NewClient crea una nueva instancia del cliente DIAN
func NewClient(config Config) (*Client, error) {
	if config.NIT == "" {
		return nil, ErrInvalidNIT
	}

	if config.Certificate.PEMPath == "" && config.Certificate.CertPEM == "" {
		return nil, ErrMissingCertificate
	}

	var certManager *signature.CertificateManager
	var err error

	if config.Certificate.PEMPath != "" {
		certManager, err = signature.NewCertificateManager(config.Certificate.PEMPath)
		if err != nil {
			return nil, fmt.Errorf("error cargando certificado: %w", err)
		}
	} else if config.Certificate.CertPEM != "" && config.Certificate.KeyPEM != "" {
		certManager, err = signature.NewCertManagerFromPEM(config.Certificate.CertPEM, config.Certificate.KeyPEM)
		if err != nil {
			return nil, fmt.Errorf("error cargando certificado PEM: %w", err)
		}
	}

	return &Client{
		Config:      config,
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
