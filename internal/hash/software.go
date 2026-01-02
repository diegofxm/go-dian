package hash

// CalculateSoftwareSecurityCode calcula el código de seguridad del software
// SoftwareSecurityCode = SHA-384(SoftwareID + SoftwarePIN)
func CalculateSoftwareSecurityCode(softwareID, pin string) string {
	data := softwareID + pin
	return CalculateSHA384(data)
}
