package server

const (
	// requestResourceNotFoundErrorCode means that the request URI was
	// correct but the requested resource was not found.
	requestResourceNotFoundErrorCode = "E0001"

	// formValidationErrorCode means that validation of request body failed.
	// When this error happens, Error.Meta should have detailed error messages
	// for each field. Each field must have a slice of strings so that clients
	// can always write error handling codes assuming that they're arrays.
	formValidationErrorCode = "E0005"

	// bqlStmtParseErrorCode is returned when a statement cannot be parsed.
	// When this error happens, Error.Meta should have parse error messages
	// in Meta["parse_errors"] as an array of strings and the statement which
	// couldn't be parsed in Meta["statement"].
	bqlStmtParseErrorCode = "E0006"

	// bqlStmtProcessingErrorCode is returned when a statement cannot be
	// processed successfully. When this error happens, Error.Meta should have
	// an error message in Meta["error"] and statement in Meta["statement"].
	bqlStmtProcessingErrorCode = "E0007"

	// nonWebSocketRequestErrorCode is returned when a requested action only
	// supports WebSocket and a request is a regular HTTP request.
	nonWebSocketRequestErrorCode = "E0008"
)
