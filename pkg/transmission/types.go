package soap

import (
	"encoding/xml"
	"time"
)

type Response struct {
	Success       bool
	StatusCode    string
	StatusMessage string
	CUFE          string
	Errors        []string
	Warnings      []string
	ResponseDate  time.Time
	RawResponse   string
}

func parseResponse(responseXML []byte) (*Response, error) {
	type ApplicationResponse struct {
		StatusCode       string `xml:"StatusCode"`
		StatusMessage    string `xml:"StatusMessage"`
		DocumentResponse struct {
			Response struct {
				ResponseCode string `xml:"ResponseCode"`
			} `xml:"Response"`
		} `xml:"DocumentResponse"`
	}

	var appResponse ApplicationResponse
	if err := xml.Unmarshal(responseXML, &appResponse); err != nil {
		return &Response{
			Success:       false,
			StatusMessage: "Error parseando respuesta de DIAN",
			RawResponse:   string(responseXML),
			ResponseDate:  time.Now(),
		}, nil
	}

	success := appResponse.StatusCode == "00" || appResponse.DocumentResponse.Response.ResponseCode == "00"

	return &Response{
		Success:       success,
		StatusCode:    appResponse.StatusCode,
		StatusMessage: appResponse.StatusMessage,
		RawResponse:   string(responseXML),
		ResponseDate:  time.Now(),
	}, nil
}
