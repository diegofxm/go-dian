package extensions

type ExtensionBuilder struct {
	NIT                  string
	SoftwareID           string
	PIN                  string
	InvoiceAuthorization string
	AuthStartDate        string
	AuthEndDate          string
	InvoicePrefix        string
	AuthFrom             string
	AuthTo               string
}

func NewExtensionBuilder(nit, softwareID string) *ExtensionBuilder {
	return &ExtensionBuilder{
		NIT:        nit,
		SoftwareID: softwareID,
	}
}

func (eb *ExtensionBuilder) WithAuthorization(auth, startDate, endDate, prefix, from, to string) *ExtensionBuilder {
	eb.InvoiceAuthorization = auth
	eb.AuthStartDate = startDate
	eb.AuthEndDate = endDate
	eb.InvoicePrefix = prefix
	eb.AuthFrom = from
	eb.AuthTo = to
	return eb
}

func (eb *ExtensionBuilder) Build(invoiceID, uuid string) *DianExtensions {
	return &DianExtensions{
		InvoiceControl: InvoiceControl{
			InvoiceAuthorization: eb.InvoiceAuthorization,
			AuthorizationPeriod: AuthorizationPeriod{
				StartDate: eb.AuthStartDate,
				EndDate:   eb.AuthEndDate,
			},
			AuthorizedInvoices: AuthorizedInvoices{
				Prefix: eb.InvoicePrefix,
				From:   eb.AuthFrom,
				To:     eb.AuthTo,
			},
		},
		InvoiceSource: InvoiceSource{
			IdentificationCode: IdentificationCode{
				Value:          "CO",
				ListAgencyID:   "6",
				ListAgencyName: "United Nations Economic Commission for Europe",
				ListSchemeURI:  "urn:oasis:names:specification:ubl:codelist:gc:CountryIdentificationCode-2.1",
			},
		},
		SoftwareProvider: SoftwareProvider{
			ProviderID: IDType{
				Value:            eb.NIT,
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				SchemeName:       "31",
				SchemeID:         "1",
			},
			SoftwareID: IDType{
				Value:            eb.SoftwareID,
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
			},
		},
		SoftwareSecurityCode: SoftwareSecurityCode{
			Value:            GenerateSoftwareSecurityCode(eb.SoftwareID, eb.PIN),
			SchemeAgencyID:   "195",
			SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
		},
		AuthorizationProvider: AuthorizationProvider{
			AuthorizationProviderID: IDType{
				Value:            "800197268",
				SchemeAgencyID:   "195",
				SchemeAgencyName: "CO, DIAN (Dirección de Impuestos y Aduanas Nacionales)",
				SchemeName:       "31",
				SchemeID:         "4",
			},
		},
		QRCode: GenerateQRCode(eb.NIT, invoiceID, uuid),
	}
}
