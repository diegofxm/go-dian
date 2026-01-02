package wssecurity

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Header representa el WS-Security Header completo
type Header struct {
	BinarySecurityToken *BinarySecurityToken
	Timestamp           *Timestamp
	Signature           *Signature
	WsaToID             string
	WsaToURL            string
}

// BinarySecurityToken representa el certificado en el header
type BinarySecurityToken struct {
	ID           string
	EncodingType string
	ValueType    string
	Value        string
}

// HeaderBuilder construye un WS-Security Header completo
type HeaderBuilder struct {
	certificate *x509.Certificate
	privateKey  *rsa.PrivateKey
	certBytes   []byte
}

// NewHeaderBuilder crea un nuevo builder
func NewHeaderBuilder(cert tls.Certificate) (*HeaderBuilder, error) {
	// Extraer certificado X.509
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("error parseando certificado: %w", err)
	}

	// Extraer llave privada RSA
	privateKey, ok := cert.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("la llave privada no es RSA")
	}

	return &HeaderBuilder{
		certificate: x509Cert,
		privateKey:  privateKey,
		certBytes:   cert.Certificate[0],
	}, nil
}

// Build construye el WS-Security Header completo con wsa:To firmado
func (hb *HeaderBuilder) Build(wsaToURL string) (*Header, error) {
	// 1. Generar IDs únicos
	securityTokenID := "SecurityToken-" + uuid.New().String()
	timestampID := "Timestamp-" + uuid.New().String()
	wsaToID := "ID-" + uuid.New().String()

	// 2. Crear BinarySecurityToken
	certBase64 := base64.StdEncoding.EncodeToString(hb.certBytes)
	binaryToken := &BinarySecurityToken{
		ID:           securityTokenID,
		EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary",
		ValueType:    "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3",
		Value:        certBase64,
	}

	// 3. Crear Timestamp (válido por 5 minutos)
	timestamp := NewTimestamp(timestampID, 5*time.Minute)

	// 4. Crear Signature (firma Timestamp Y wsa:To - requerido por DIAN)
	signer := NewSigner(hb.privateKey, hb.certificate, hb.certBytes)
	signature, err := signer.SignTimestamp(timestamp, securityTokenID, wsaToID, wsaToURL)
	if err != nil {
		return nil, fmt.Errorf("error firmando timestamp: %w", err)
	}

	return &Header{
		BinarySecurityToken: binaryToken,
		Timestamp:           timestamp,
		Signature:           signature,
		WsaToID:             wsaToID,
		WsaToURL:            wsaToURL,
	}, nil
}

// ToXML genera el XML completo del WS-Security Header con wsa:Action y wsa:To
func (h *Header) ToXML(wsaAction string) string {
	xml := `<wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">`

	// BinarySecurityToken
	xml += fmt.Sprintf(`<wsse:BinarySecurityToken wsu:Id="%s" EncodingType="%s" ValueType="%s">%s</wsse:BinarySecurityToken>`,
		h.BinarySecurityToken.ID,
		h.BinarySecurityToken.EncodingType,
		h.BinarySecurityToken.ValueType,
		h.BinarySecurityToken.Value)

	// Timestamp
	xml += h.Timestamp.ToXML()

	// Signature
	xml += h.Signature.ToXML()

	xml += `</wsse:Security>`

	// wsa:Action (fuera de wsse:Security)
	xml += fmt.Sprintf(`<wsa:Action xmlns:wsa="http://www.w3.org/2005/08/addressing">%s</wsa:Action>`, wsaAction)

	// wsa:To (fuera de wsse:Security, con wsu:Id para que pueda ser firmado)
	xml += fmt.Sprintf(`<wsa:To xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" wsu:Id="%s">%s</wsa:To>`, h.WsaToID, h.WsaToURL)

	return xml
}
