package soap

import (
	"fmt"
)

func (c *Client) SendInvoiceSync(fileName string, signedXML []byte) (*Response, error) {
	return c.SendInvoice(fileName, signedXML)
}

func (c *Client) SendInvoiceAsync(fileName string, signedXML []byte) (*Response, error) {
	return nil, fmt.Errorf("envío asíncrono no implementado aún")
}

func (c *Client) GetInvoiceStatus(trackingID string) (*Response, error) {
	return nil, fmt.Errorf("consulta de estado no implementada aún")
}

func (c *Client) SendTestSetInvoice(testSetID string, signedXML []byte) (*Response, error) {
	return c.SendTestSet(testSetID, signedXML)
}
