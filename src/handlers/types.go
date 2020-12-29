package handlers

// ResJSONError is an error field that is used in JSON error responses
type ResJSONError struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}
