package builtin

import (
	"errors"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"unicode/utf8"
)

// blobToRawString converts a blob to a string without encoding the value in
// base64. It returns an error when the blob contains invalid UTF-8 runes.
func blobToRawString(b data.Blob) (data.String, error) {
	if !utf8.Valid([]byte(b)) {
		return data.String(""), errors.New("blob contains invalid UTF-8 runes")
	}
	return data.String(b), nil
}
