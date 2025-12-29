package signature

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

func LoadCertificate(path, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	if strings.HasSuffix(strings.ToLower(path), ".pem") {
		return loadFromPEM(path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error leyendo certificado: %w", err)
	}

	privateKey, certificate, err := pkcs12.Decode(data, password)
	if err != nil {
		return convertP12ToPEMAndLoad(path, password)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("la clave privada no es RSA")
	}

	return certificate, rsaKey, nil
}

func loadFromPEM(pemPath string) (*x509.Certificate, *rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error leyendo archivo PEM: %w", err)
	}

	return parsePEMData(pemData)
}

// LoadPEMStrings carga certificado y clave privada desde strings PEM
func LoadPEMStrings(certPEM, keyPEM string) (*x509.Certificate, *rsa.PrivateKey, error) {
	if certPEM == "" || keyPEM == "" {
		return nil, nil, fmt.Errorf("certificado y clave PEM son requeridos")
	}

	// Combinar certificado y clave en un solo string para parsear
	pemData := []byte(certPEM + "\n" + keyPEM)
	return parsePEMData(pemData)
}

func parsePEMData(pemData []byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	var cert *x509.Certificate
	var key *rsa.PrivateKey
	var err error

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

	if cert == nil || key == nil {
		return nil, nil, fmt.Errorf("certificado o clave privada no encontrados en PEM")
	}

	return cert, key, nil
}

func convertP12ToPEMAndLoad(p12Path, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	pemPath := strings.TrimSuffix(p12Path, ".p12") + ".pem"

	// Exportar certificado y clave privada en un solo archivo PEM
	// -clcerts: solo certificados de cliente
	// -nokeys: no exportar claves (para el primer comando)
	// Usamos -nodes para no encriptar la clave privada
	cmd := exec.Command("openssl", "pkcs12", "-in", p12Path, "-out", pemPath, "-nodes", "-passin", "pass:"+password, "-legacy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Intentar sin -legacy para versiones antiguas de OpenSSL
		cmd = exec.Command("openssl", "pkcs12", "-in", p12Path, "-out", pemPath, "-nodes", "-passin", "pass:"+password)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return nil, nil, fmt.Errorf("error convirtiendo P12 a PEM: %w (output: %s)", err, string(output))
		}
	}

	return loadFromPEM(pemPath)
}

func SignXMLDocument(xmlData []byte, cert *x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256(xmlData)
	digestValue := base64.StdEncoding.EncodeToString(hash[:])

	signedInfo := SignedInfo{
		CanonicalizationMethod: CanonicalizationMethod{
			Algorithm: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315",
		},
		SignatureMethod: SignatureMethod{
			Algorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		},
		Reference: []Reference{
			{
				URI: "",
				Transforms: &Transforms{
					Transform: []Transform{
						{Algorithm: "http://www.w3.org/2000/09/xmldsig#enveloped-signature"},
					},
				},
				DigestMethod: DigestMethod{
					Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
				},
				DigestValue: digestValue,
			},
		},
	}

	signedInfoXML, err := xml.Marshal(signedInfo)
	if err != nil {
		return nil, fmt.Errorf("error serializando SignedInfo: %w", err)
	}

	signedInfoHash := sha256.Sum256(signedInfoXML)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, signedInfoHash[:])
	if err != nil {
		return nil, fmt.Errorf("error firmando: %w", err)
	}

	certDER := base64.StdEncoding.EncodeToString(cert.Raw)

	sig := Signature{
		ID:         "xmldsig-signature",
		SignedInfo: signedInfo,
		SignatureValue: SignatureValue{
			Value: base64.StdEncoding.EncodeToString(signature),
		},
		KeyInfo: KeyInfo{
			X509Data: X509Data{
				X509Certificate: certDER,
			},
		},
		Object: &Object{
			QualifyingProperties: QualifyingProperties{
				Target: "#xmldsig-signature",
				SignedProperties: SignedProperties{
					ID: "xmldsig-signedproperties",
					SignedSignatureProperties: SignedSignatureProperties{
						SigningTime: time.Now().Format(time.RFC3339),
						SigningCertificate: SigningCertificate{
							Cert: []Cert{
								{
									CertDigest: CertDigest{
										DigestMethod: DigestMethod{
											Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256",
										},
										DigestValue: base64.StdEncoding.EncodeToString(hash[:]),
									},
									IssuerSerial: IssuerSerial{
										X509IssuerName:   cert.Issuer.String(),
										X509SerialNumber: cert.SerialNumber.String(),
									},
								},
							},
						},
						SignaturePolicyIdentifier: SignaturePolicyIdentifier{
							SignaturePolicyId: SignaturePolicyId{
								SigPolicyId: SigPolicyId{
									Identifier: "https://facturaelectronica.dian.gov.co/politicadefirma/v2/politicadefirmav2.pdf",
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

	return xml.MarshalIndent(sig, "      ", "  ")
}
