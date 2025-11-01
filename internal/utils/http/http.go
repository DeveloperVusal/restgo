package http

type ContentType string

const (
	ContentType_JSON ContentType = "application/json"
	ContentType_XML  ContentType = "application/xml"
	ContentType_HTML ContentType = "text/html"
	ContentType_TEXT ContentType = "text/plain"
)
