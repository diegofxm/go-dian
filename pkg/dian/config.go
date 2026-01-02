package dian

// Config contiene la configuración del cliente
type Config struct {
	NIT          string
	Certificate  Certificate
	Environment  Environment
	SoftwareID   string
	TestSetID    string
	TechnicalKey string // Clave técnica del software (para CUFE)
	PIN          string // PIN del software (para SoftwareSecurityCode)

	// Datos de autorización DIAN (específicos por empresa)
	InvoiceAuthorization string // Número de autorización DIAN
	AuthStartDate        string // Fecha inicio autorización (YYYY-MM-DD)
	AuthEndDate          string // Fecha fin autorización (YYYY-MM-DD)
	InvoicePrefix        string // Prefijo de facturación
	AuthFrom             string // Consecutivo desde
	AuthTo               string // Consecutivo hasta
}

// Certificate representa el certificado digital (solo PEM)
type Certificate struct {
	PEMPath string // Ruta a certificado PEM
	CertPEM string // Certificado PEM como string (para BD)
	KeyPEM  string // Clave privada PEM como string (para BD)
}

// Environment define el ambiente de DIAN
type Environment string

const (
	EnvironmentProduction Environment = "production"
	EnvironmentTest       Environment = "test"
)
