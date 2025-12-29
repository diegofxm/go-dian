package dian

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

// Environment define el ambiente de DIAN
type Environment string

const (
	EnvironmentProduction Environment = "production"
	EnvironmentTest       Environment = "test"
)
