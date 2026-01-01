package wssecurity

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// Signature representa una firma digital XMLDSig
type Signature struct {
	SignedInfo     *SignedInfo
	SignatureValue string
	KeyInfo        *KeyInfo
}

// SignedInfo contiene la información firmada
type SignedInfo struct {
	CanonicalizationMethod string
	SignatureMethod        string
	References             []*Reference
}

// Reference representa una referencia a un elemento firmado
type Reference struct {
	URI          string
	DigestMethod string
	DigestValue  string
	Transforms   []string
}

// KeyInfo contiene información sobre la clave pública
type KeyInfo struct {
	SecurityTokenReference *SecurityTokenReference
}

// SecurityTokenReference referencia al token de seguridad
type SecurityTokenReference struct {
	Reference string
}

// Signer maneja la firma de documentos
type Signer struct {
	privateKey  *rsa.PrivateKey
	certificate *x509.Certificate
	certBytes   []byte
}

// NewSigner crea un nuevo firmador con certificado y llave privada
func NewSigner(privateKey *rsa.PrivateKey, certificate *x509.Certificate, certBytes []byte) *Signer {
	return &Signer{
		privateKey:  privateKey,
		certificate: certificate,
		certBytes:   certBytes,
	}
}

// SignTimestamp firma un timestamp y retorna una Signature completa
func (s *Signer) SignTimestamp(timestamp *Timestamp, securityTokenID string) (*Signature, error) {
	// 1. Generar DigestValue del Timestamp
	timestampXML := timestamp.ToXML()
	timestampCanonical, err := CanonicalizeSigned(timestampXML)
	if err != nil {
		return nil, fmt.Errorf("error canonicalizando timestamp: %w", err)
	}

	timestampDigest := sha256.Sum256(timestampCanonical)
	timestampDigestB64 := base64.StdEncoding.EncodeToString(timestampDigest[:])

	// 2. Crear Reference para el Timestamp
	ref := &Reference{
		URI:          "#" + timestamp.ID,
		DigestMethod: "http://www.w3.org/2001/04/xmlenc#sha256",
		DigestValue:  timestampDigestB64,
		Transforms: []string{
			"http://www.w3.org/2001/10/xml-exc-c14n#",
		},
	}

	// 3. Crear SignedInfo
	signedInfo := &SignedInfo{
		CanonicalizationMethod: "http://www.w3.org/2001/10/xml-exc-c14n#",
		SignatureMethod:        "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		References:             []*Reference{ref},
	}

	// 4. Canonicalizar SignedInfo
	signedInfoXML := signedInfo.ToXML()
	signedInfoCanonical, err := CanonicalizeSigned(signedInfoXML)
	if err != nil {
		return nil, fmt.Errorf("error canonicalizando SignedInfo: %w", err)
	}

	// 5. Firmar SignedInfo
	signedInfoHash := sha256.Sum256(signedInfoCanonical)
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, signedInfoHash[:])
	if err != nil {
		return nil, fmt.Errorf("error firmando: %w", err)
	}

	signatureValue := base64.StdEncoding.EncodeToString(signatureBytes)

	// 6. Crear KeyInfo
	keyInfo := &KeyInfo{
		SecurityTokenReference: &SecurityTokenReference{
			Reference: "#" + securityTokenID,
		},
	}

	return &Signature{
		SignedInfo:     signedInfo,
		SignatureValue: signatureValue,
		KeyInfo:        keyInfo,
	}, nil
}

// ToXML genera el XML de SignedInfo
func (si *SignedInfo) ToXML() string {
	xml := `<ds:SignedInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">`
	xml += fmt.Sprintf(`<ds:CanonicalizationMethod Algorithm="%s"></ds:CanonicalizationMethod>`, si.CanonicalizationMethod)
	xml += fmt.Sprintf(`<ds:SignatureMethod Algorithm="%s"></ds:SignatureMethod>`, si.SignatureMethod)

	for _, ref := range si.References {
		xml += fmt.Sprintf(`<ds:Reference URI="%s">`, ref.URI)
		xml += `<ds:Transforms>`
		for _, transform := range ref.Transforms {
			xml += fmt.Sprintf(`<ds:Transform Algorithm="%s"></ds:Transform>`, transform)
		}
		xml += `</ds:Transforms>`
		xml += fmt.Sprintf(`<ds:DigestMethod Algorithm="%s"></ds:DigestMethod>`, ref.DigestMethod)
		xml += fmt.Sprintf(`<ds:DigestValue>%s</ds:DigestValue>`, ref.DigestValue)
		xml += `</ds:Reference>`
	}

	xml += `</ds:SignedInfo>`
	return xml
}

// ToXML genera el XML completo de la Signature
func (sig *Signature) ToXML() string {
	xml := `<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">`
	xml += sig.SignedInfo.ToXML()
	xml += fmt.Sprintf(`<ds:SignatureValue>%s</ds:SignatureValue>`, sig.SignatureValue)
	xml += `<ds:KeyInfo>`
	xml += `<wsse:SecurityTokenReference xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">`
	xml += fmt.Sprintf(`<wsse:Reference URI="%s" ValueType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3"></wsse:Reference>`, sig.KeyInfo.SecurityTokenReference.Reference)
	xml += `</wsse:SecurityTokenReference>`
	xml += `</ds:KeyInfo>`
	xml += `</ds:Signature>`
	return xml
}
