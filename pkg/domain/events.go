package domain

type Event map[string]any

type EventUpdatePDFField struct {
	Headers   map[string]string `mapstructure:"HEADERS"`
	Operation string            `mapstructure:"OPERATION"`

	PDFID   string `mapstructure:"pdf_id"`
	PageNum int    `mapstructure:"page_num"`
}
