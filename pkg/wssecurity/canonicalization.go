package wssecurity

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExclusiveC14N implementa Exclusive XML Canonicalization (exc-c14n)
// según el estándar: https://www.w3.org/TR/xml-exc-c14n/
type ExclusiveC14N struct {
	inclusiveNamespaces []string
}

// NewExclusiveC14N crea un nuevo canonicalizador
func NewExclusiveC14N(inclusiveNamespaces ...string) *ExclusiveC14N {
	return &ExclusiveC14N{
		inclusiveNamespaces: inclusiveNamespaces,
	}
}

// Canonicalize canonicaliza un fragmento XML
func (c *ExclusiveC14N) Canonicalize(xmlData []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	var buf bytes.Buffer

	var nsContext []map[string]string
	nsContext = append(nsContext, make(map[string]string))

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing XML: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			if err := c.canonicalizeStartElement(&buf, t, &nsContext); err != nil {
				return nil, err
			}

		case xml.EndElement:
			buf.WriteString("</")
			if t.Name.Space != "" {
				prefix := c.findPrefix(t.Name.Space, nsContext)
				if prefix != "" {
					buf.WriteString(prefix)
					buf.WriteString(":")
				}
			}
			buf.WriteString(t.Name.Local)
			buf.WriteString(">")

			if len(nsContext) > 1 {
				nsContext = nsContext[:len(nsContext)-1]
			}

		case xml.CharData:
			c.canonicalizeText(&buf, string(t))

		case xml.Comment:
			// Los comentarios se omiten en c14n

		case xml.ProcInst:
			// Las instrucciones de procesamiento se omiten
		}
	}

	return buf.Bytes(), nil
}

func (c *ExclusiveC14N) canonicalizeStartElement(buf *bytes.Buffer, elem xml.StartElement, nsContext *[]map[string]string) error {
	// Crear nuevo contexto de namespace
	currentNS := make(map[string]string)
	for k, v := range (*nsContext)[len(*nsContext)-1] {
		currentNS[k] = v
	}

	// Agregar namespaces del elemento
	for _, attr := range elem.Attr {
		if attr.Name.Space == "xmlns" || attr.Name.Local == "xmlns" {
			prefix := attr.Name.Local
			if attr.Name.Space == "xmlns" {
				prefix = attr.Name.Local
			} else if attr.Name.Local == "xmlns" {
				prefix = ""
			}
			currentNS[prefix] = attr.Value
		}
	}
	*nsContext = append(*nsContext, currentNS)

	// Escribir elemento
	buf.WriteString("<")
	if elem.Name.Space != "" {
		prefix := c.findPrefix(elem.Name.Space, *nsContext)
		if prefix != "" {
			buf.WriteString(prefix)
			buf.WriteString(":")
		}
	}
	buf.WriteString(elem.Name.Local)

	// Recolectar y ordenar atributos
	var attrs []xml.Attr
	var nsAttrs []xml.Attr

	for _, attr := range elem.Attr {
		if attr.Name.Space == "xmlns" || attr.Name.Local == "xmlns" {
			nsAttrs = append(nsAttrs, attr)
		} else {
			attrs = append(attrs, attr)
		}
	}

	// Ordenar namespaces
	sort.Slice(nsAttrs, func(i, j int) bool {
		return nsAttrs[i].Name.Local < nsAttrs[j].Name.Local
	})

	// Ordenar atributos
	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i].Name.Space != attrs[j].Name.Space {
			return attrs[i].Name.Space < attrs[j].Name.Space
		}
		return attrs[i].Name.Local < attrs[j].Name.Local
	})

	// Escribir namespaces
	for _, attr := range nsAttrs {
		buf.WriteString(" ")
		if attr.Name.Local == "xmlns" {
			buf.WriteString("xmlns")
		} else {
			buf.WriteString("xmlns:")
			buf.WriteString(attr.Name.Local)
		}
		buf.WriteString("=\"")
		c.escapeAttribute(buf, attr.Value)
		buf.WriteString("\"")
	}

	// Escribir atributos
	for _, attr := range attrs {
		buf.WriteString(" ")
		if attr.Name.Space != "" {
			prefix := c.findPrefix(attr.Name.Space, *nsContext)
			if prefix != "" {
				buf.WriteString(prefix)
				buf.WriteString(":")
			}
		}
		buf.WriteString(attr.Name.Local)
		buf.WriteString("=\"")
		c.escapeAttribute(buf, attr.Value)
		buf.WriteString("\"")
	}

	buf.WriteString(">")
	return nil
}

func (c *ExclusiveC14N) canonicalizeText(buf *bytes.Buffer, text string) {
	for _, r := range text {
		switch r {
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '&':
			buf.WriteString("&amp;")
		case '\r':
			buf.WriteString("&#xD;")
		default:
			buf.WriteRune(r)
		}
	}
}

func (c *ExclusiveC14N) escapeAttribute(buf *bytes.Buffer, s string) {
	for _, r := range s {
		switch r {
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '&':
			buf.WriteString("&amp;")
		case '"':
			buf.WriteString("&quot;")
		case '\t':
			buf.WriteString("&#x9;")
		case '\n':
			buf.WriteString("&#xA;")
		case '\r':
			buf.WriteString("&#xD;")
		default:
			buf.WriteRune(r)
		}
	}
}

func (c *ExclusiveC14N) findPrefix(namespace string, nsContext []map[string]string) string {
	for i := len(nsContext) - 1; i >= 0; i-- {
		for prefix, ns := range nsContext[i] {
			if ns == namespace {
				return prefix
			}
		}
	}
	return ""
}

// CanonicalizeSigned canonicaliza específicamente para SignedInfo
func CanonicalizeSigned(xmlData string) ([]byte, error) {
	// Limpiar espacios en blanco innecesarios
	xmlData = strings.TrimSpace(xmlData)

	c14n := NewExclusiveC14N()
	return c14n.Canonicalize([]byte(xmlData))
}
