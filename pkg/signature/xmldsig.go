package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"
)

func CreateXMLSignature(xmlData []byte, cert *x509.Certificate, privateKey *rsa.PrivateKey) (*Signature, error) {
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

	sig := &Signature{
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

	return sig, nil
}
