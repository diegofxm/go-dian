package signature

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
)

type CertificateManager struct {
	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey
}

func NewCertificateManager(certPath, password string) (*CertificateManager, error) {
	cert, key, err := LoadCertificate(certPath, password)
	if err != nil {
		return nil, fmt.Errorf("error cargando certificado: %w", err)
	}

	return &CertificateManager{
		Certificate: cert,
		PrivateKey:  key,
	}, nil
}

// NewCertManagerFromPEM crea un CertificateManager desde strings PEM
func NewCertManagerFromPEM(certPEM, keyPEM string) (*CertificateManager, error) {
	cert, key, err := LoadPEMStrings(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("error cargando certificado PEM: %w", err)
	}

	return &CertificateManager{
		Certificate: cert,
		PrivateKey:  key,
	}, nil
}

func (cm *CertificateManager) GetCertificate() *x509.Certificate {
	return cm.Certificate
}

func (cm *CertificateManager) GetPrivateKey() *rsa.PrivateKey {
	return cm.PrivateKey
}

func (cm *CertificateManager) Validate() error {
	if cm.Certificate == nil {
		return fmt.Errorf("certificado no cargado")
	}
	if cm.PrivateKey == nil {
		return fmt.Errorf("clave privada no cargada")
	}
	return nil
}
