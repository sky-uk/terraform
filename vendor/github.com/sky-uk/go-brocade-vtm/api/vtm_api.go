package api

// VTMApi object.
type VTMApi interface {
	Method() string
	Endpoint() string
	RequestObject() interface{}
	ResponseObject() interface{}
	StatusCode() int
	RawResponse() []byte
	Error() error

	SetResponseObject(interface{})
	SetStatusCode(int)
	SetRawResponse([]byte)
}
