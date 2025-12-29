package xml

import (
	"encoding/xml"
	"fmt"
	"strings"
)

func InsertSignature(invoiceXML, signatureXML []byte) ([]byte, error) {
	invoiceStr := string(invoiceXML)

	extensionsEnd := strings.Index(invoiceStr, "</ext:UBLExtensions>")
	if extensionsEnd == -1 {
		return nil, fmt.Errorf("no se encontró UBLExtensions en la factura")
	}

	signatureExtension := fmt.Sprintf(`  <ext:UBLExtension>
    <ext:ExtensionContent>
%s
    </ext:ExtensionContent>
  </ext:UBLExtension>
`, string(signatureXML))

	result := invoiceStr[:extensionsEnd] + signatureExtension + invoiceStr[extensionsEnd:]
	return []byte(result), nil
}

func PrettyPrint(data []byte) ([]byte, error) {
	var v interface{}
	if err := xml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return xml.MarshalIndent(v, "", "  ")
}
