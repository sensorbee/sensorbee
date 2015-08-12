package server

import (
	"encoding/json"
	"net/http"
)

func parseRequestJSON(js []byte) (map[string]interface{}, *Error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal(js, &parsed); err != nil {
		return nil, NewError(requestBodyParseErrorCode, "Parsing the requested JSON failed. Please check if the JSON is valid.",
			http.StatusBadRequest, err)
	}
	return parsed, nil
}

func parseJSONFromRequestBody(c *Context) (map[string]interface{}, *Error) {
	body, err := c.Body()
	if err != nil {
		return nil, NewError(requestBodyReadErrorCode, "Cannot read the request body.",
			http.StatusInternalServerError, err)
	}
	return parseRequestJSON(body)
}
