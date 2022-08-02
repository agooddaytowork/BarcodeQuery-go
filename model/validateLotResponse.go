package model

type ValidateLotResponsePayload struct {
	UnexpectedSerialNumbers []string `json:"unexpected_serial_numbers"`
	ErrorDetails            string   `json:"error_details"`
}
