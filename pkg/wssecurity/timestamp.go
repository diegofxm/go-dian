package wssecurity

import (
	"fmt"
	"time"
)

// Timestamp representa un timestamp de WS-Security
type Timestamp struct {
	ID      string
	Created string
	Expires string
}

// NewTimestamp crea un nuevo timestamp con duración especificada
func NewTimestamp(id string, duration time.Duration) *Timestamp {
	now := time.Now().UTC()
	expires := now.Add(duration)

	return &Timestamp{
		ID:      id,
		Created: now.Format("2006-01-02T15:04:05.000Z"),
		Expires: expires.Format("2006-01-02T15:04:05.000Z"),
	}
}

// ToXML genera el XML del timestamp
func (t *Timestamp) ToXML() string {
	return fmt.Sprintf(`<wsu:Timestamp xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" wsu:Id="%s"><wsu:Created>%s</wsu:Created><wsu:Expires>%s</wsu:Expires></wsu:Timestamp>`,
		t.ID, t.Created, t.Expires)
}
