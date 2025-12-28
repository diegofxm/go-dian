package dian

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

// Signature representa la estructura de firma XMLDSig
type Signature struct {
	XMLName        xml.Name       `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
	ID             string         `xml:"Id,attr"`
	SignedInfo     SignedInfo     `xml:"SignedInfo"`
	SignatureValue SignatureValue `xml:"SignatureValue"`
	KeyInfo        KeyInfo        `xml:"KeyInfo"`
	Object         *Object        `xml:"Object,omitempty"`
}

type SignedInfo struct {
	XMLName                xml.Name               `xml:"SignedInfo"`
	CanonicalizationMethod CanonicalizationMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethod        `xml:"SignatureMethod"`
	Reference              []Reference            `xml:"Reference"`
}

type CanonicalizationMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type Reference struct {
	ID           string       `xml:"Id,attr,omitempty"`
	URI          string       `xml:"URI,attr"`
	Type         string       `xml:"Type,attr,omitempty"`
	Transforms   *Transforms  `xml:"Transforms,omitempty"`
	DigestMethod DigestMethod `xml:"DigestMethod"`
	DigestValue  string       `xml:"DigestValue"`
}

type Transforms struct {
	Transform []Transform `xml:"Transform"`
}

type Transform struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type DigestMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type SignatureValue struct {
	ID    string `xml:"Id,attr,omitempty"`
	Value string `xml:",chardata"`
}

type KeyInfo struct {
	ID       string   `xml:"Id,attr,omitempty"`
	X509Data X509Data `xml:"X509Data"`
}

type X509Data struct {
	X509Certificate string `xml:"X509Certificate"`
}

type Object struct {
	QualifyingProperties QualifyingProperties `xml:"http://uri.etsi.org/01903/v1.3.2# QualifyingProperties"`
}

type QualifyingProperties struct {
	Target           string           `xml:"Target,attr"`
	SignedProperties SignedProperties `xml:"SignedProperties"`
}

type SignedProperties struct {
	ID                        string                    `xml:"Id,attr"`
	SignedSignatureProperties SignedSignatureProperties `xml:"SignedSignatureProperties"`
}

type SignedSignatureProperties struct {
	SigningTime               string                    `xml:"SigningTime"`
	SigningCertificate        SigningCertificate        `xml:"SigningCertificate"`
	SignaturePolicyIdentifier SignaturePolicyIdentifier `xml:"SignaturePolicyIdentifier"`
}

type SigningCertificate struct {
	Cert []Cert `xml:"Cert"`
}

type Cert struct {
	CertDigest   CertDigest   `xml:"CertDigest"`
	IssuerSerial IssuerSerial `xml:"IssuerSerial"`
}

type CertDigest struct {
	DigestMethod DigestMethod `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue  string       `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

type IssuerSerial struct {
	X509IssuerName   string `xml:"http://www.w3.org/2000/09/xmldsig# X509IssuerName"`
	X509SerialNumber string `xml:"http://www.w3.org/2000/09/xmldsig# X509SerialNumber"`
}

type SignaturePolicyIdentifier struct {
	SignaturePolicyId SignaturePolicyId `xml:"SignaturePolicyId"`
}

type SignaturePolicyId struct {
	SigPolicyId   SigPolicyId   `xml:"SigPolicyId"`
	SigPolicyHash SigPolicyHash `xml:"SigPolicyHash"`
}

type SigPolicyId struct {
	Identifier  string `xml:"Identifier"`
	Description string `xml:"Description,omitempty"`
}

type SigPolicyHash struct {
	DigestMethod DigestMethod `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue  string       `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

// LoadCertificate carga un certificado (P12 o PEM) automáticamente
func LoadCertificate(path, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Detectar tipo de archivo
	if strings.HasSuffix(strings.ToLower(path), ".pem") {
		return loadFromPEM(path)
	}

	// Si es P12, intentar decodificar directamente
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error leyendo certificado: %w", err)
	}

	privateKey, certificate, err := pkcs12.Decode(data, password)
	if err != nil {
		// Si falla, convertir P12 a PEM usando OpenSSL y reintentar
		return convertP12ToPEMAndLoad(path, password)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("la clave privada no es RSA")
	}

	return certificate, rsaKey, nil
}

// loadFromPEM carga certificado y clave desde archivo PEM
func loadFromPEM(pemPath string) (*x509.Certificate, *rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error leyendo archivo PEM: %w", err)
	}

	var cert *x509.Certificate
	var key *rsa.PrivateKey

	// Decodificar todos los bloques PEM
	rest := pemData
	for {
		block, remaining := pem.Decode(rest)
		if block == nil {
			break
		}
		rest = remaining

		switch block.Type {
		case "CERTIFICATE":
			if cert == nil {
				cert, err = x509.ParseCertificate(block.Bytes)
				if err != nil {
					return nil, nil, fmt.Errorf("error parseando certificado: %w", err)
				}
			}

		case "PRIVATE KEY":
			parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("error parseando clave PKCS8: %w", err)
			}
			var ok bool
			key, ok = parsedKey.(*rsa.PrivateKey)
			if !ok {
				return nil, nil, fmt.Errorf("la clave privada no es RSA")
			}

		case "RSA PRIVATE KEY":
			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, fmt.Errorf("error parseando clave RSA: %w", err)
			}
		}
	}

	if cert == nil {
		return nil, nil, fmt.Errorf("certificado no encontrado en archivo PEM")
	}
	if key == nil {
		return nil, nil, fmt.Errorf("clave privada no encontrada en archivo PEM")
	}

	return cert, key, nil
}

// LoadCertificateFromPEMStrings carga certificado desde strings PEM (para BD)
func LoadCertificateFromPEMStrings(certPEM, keyPEM string) (*x509.Certificate, *rsa.PrivateKey, error) {
	var cert *x509.Certificate
	var key *rsa.PrivateKey

	// Parsear certificado
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, nil, fmt.Errorf("error decodificando certificado PEM")
	}

	var err error
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parseando certificado: %w", err)
	}

	// Parsear clave privada
	keyBlock, _ := pem.Decode([]byte(keyPEM))
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("error decodificando clave privada PEM")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			return nil, nil, fmt.Errorf("error parseando clave privada: %w", err)
		}
	}

	var ok bool
	key, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("la clave privada no es RSA")
	}

	return cert, key, nil
}

// convertP12ToPEMAndLoad convierte P12 a PEM usando OpenSSL y lo carga
func convertP12ToPEMAndLoad(p12Path, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Generar nombre del archivo PEM (mismo nombre que P12 pero con extensión .pem)
	pemPath := p12Path
	pemPath = strings.TrimSuffix(pemPath, ".p12")
	pemPath = strings.TrimSuffix(pemPath, ".pfx")
	pemPath = pemPath + ".pem"

	// Si ya existe un PEM, verificar si es del mismo P12 (comparar fecha de modificación)
	p12Info, err := os.Stat(p12Path)
	if err != nil {
		return nil, nil, fmt.Errorf("error obteniendo info del P12: %w", err)
	}

	pemInfo, err := os.Stat(pemPath)
	if err == nil {
		// PEM existe, verificar si es más antiguo que el P12
		if pemInfo.ModTime().Before(p12Info.ModTime()) {
			// P12 es más nuevo, eliminar PEM antiguo
			os.Remove(pemPath)
		} else {
			// PEM es actual, usarlo directamente
			return loadFromPEM(pemPath)
		}
	}

	// Ejecutar OpenSSL para convertir P12 a PEM (con -legacy para soportar RC2-40-CBC)
	cmd := exec.Command("openssl", "pkcs12", "-in", p12Path, "-out", pemPath, "-nodes", "-legacy", "-passin", "pass:"+password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil, fmt.Errorf("error convirtiendo P12 a PEM con OpenSSL: %w\nOutput: %s", err, string(output))
	}

	// Cargar el archivo PEM convertido
	cert, key, err := loadFromPEM(pemPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error cargando PEM convertido: %w", err)
	}

	return cert, key, nil
}

// SignXMLDocument firma un documento XML con XMLDSig
func SignXMLDocument(xmlData []byte, cert *x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	signatureID := generateID()

	// Calcular digest del documento
	documentHash := sha256.Sum256(xmlData)
	documentDigest := base64.StdEncoding.EncodeToString(documentHash[:])

	// Calcular digest del certificado
	certHash := sha256.Sum256(cert.Raw)
	certDigest := base64.StdEncoding.EncodeToString(certHash[:])

	// Crear SignedInfo
	signedInfo := SignedInfo{
		CanonicalizationMethod: CanonicalizationMethod{
			Algorithm: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315",
		},
		SignatureMethod: SignatureMethod{
			Algorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		},
		Reference: []Reference{
			{
				ID:  signatureID + "-ref0",
				URI: "",
				Transforms: &Transforms{
					Transform: []Transform{
						{Algorithm: "http://www.w3.org/2000/09/xmldsig#enveloped-signature"},
					},
				},
				DigestMethod: DigestMethod{
					Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
				},
				DigestValue: documentDigest,
			},
		},
	}

	// Serializar SignedInfo
	signedInfoXML, err := xml.Marshal(signedInfo)
	if err != nil {
		return nil, fmt.Errorf("error serializando SignedInfo: %w", err)
	}

	// Firmar SignedInfo
	signedInfoHash := sha256.Sum256(signedInfoXML)
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, signedInfoHash[:])
	if err != nil {
		return nil, fmt.Errorf("error firmando documento: %w", err)
	}

	signatureValue := base64.StdEncoding.EncodeToString(signatureBytes)

	// Codificar certificado en base64
	certPEM := base64.StdEncoding.EncodeToString(cert.Raw)

	// Crear estructura de firma completa
	signature := Signature{
		ID:         signatureID,
		SignedInfo: signedInfo,
		SignatureValue: SignatureValue{
			ID:    signatureID + "-sigvalue",
			Value: signatureValue,
		},
		KeyInfo: KeyInfo{
			ID: signatureID + "-keyinfo",
			X509Data: X509Data{
				X509Certificate: certPEM,
			},
		},
		Object: &Object{
			QualifyingProperties: QualifyingProperties{
				Target: "#" + signatureID,
				SignedProperties: SignedProperties{
					ID: signatureID + "-signedprops",
					SignedSignatureProperties: SignedSignatureProperties{
						SigningTime: time.Now().Format("2006-01-02T15:04:05-07:00"),
						SigningCertificate: SigningCertificate{
							Cert: []Cert{
								{
									CertDigest: CertDigest{
										DigestMethod: DigestMethod{
											Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
										},
										DigestValue: certDigest,
									},
									IssuerSerial: IssuerSerial{
										X509IssuerName:   cert.Issuer.String(),
										X509SerialNumber: "0",
									},
								},
							},
						},
						SignaturePolicyIdentifier: SignaturePolicyIdentifier{
							SignaturePolicyId: SignaturePolicyId{
								SigPolicyId: SigPolicyId{
									Identifier:  "https://facturaelectronica.dian.gov.co/politicadefirma/v2/politicadefirmav2.pdf",
									Description: "Política de firma para facturas electrónicas de la República de Colombia.",
								},
								SigPolicyHash: SigPolicyHash{
									DigestMethod: DigestMethod{
										Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
									},
									DigestValue: "dMoMvtcG5aIzgYo0tIsSQeVJBDnUnfSOfBpxXrmor0Y=",
								},
							},
						},
					},
				},
			},
		},
	}

	// Serializar firma
	signatureXML, err := xml.MarshalIndent(signature, "  ", "  ")
	if err != nil {
		return nil, fmt.Errorf("error serializando firma: %w", err)
	}

	// Retornar SOLO la firma XMLDSig (no el documento completo)
	return signatureXML, nil
}

// generateID genera un ID único para la firma
func generateID() string {
	return fmt.Sprintf("xmldsig-%d", time.Now().UnixNano())
}

// VerifySignature verifica una firma XMLDSig (para testing)
func VerifySignature(signedXML []byte, cert *x509.Certificate) error {
	// TODO: Implementar verificación de firma
	// Por ahora retorna nil (asume válido)
	return nil
}

// GetCertificateInfo obtiene información del certificado
func GetCertificateInfo(cert *x509.Certificate) map[string]string {
	return map[string]string{
		"Subject":      cert.Subject.String(),
		"Issuer":       cert.Issuer.String(),
		"NotBefore":    cert.NotBefore.String(),
		"NotAfter":     cert.NotAfter.String(),
		"SerialNumber": cert.SerialNumber.String(),
	}
}

// EncodeCertificatePEM codifica un certificado en formato PEM
func EncodeCertificatePEM(cert *x509.Certificate) string {
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	return string(certPEM)
}
