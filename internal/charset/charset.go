// Package charset provides character encoding conversion utilities.
package charset

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

// ToUTF8 converts data from the given charset to UTF-8.
// Returns the original data unchanged if charset is already UTF-8.
func ToUTF8(data []byte, charsetName string) ([]byte, error) {
	charsetName = strings.ToLower(charsetName)
	switch charsetName {
	case "utf-8", "utf8":
		return data, nil
	}

	enc, err := ianaindex.IANA.Encoding(charsetName)
	if err != nil {
		return nil, fmt.Errorf("unsupported charset: %s", charsetName)
	}
	if enc == nil {
		// Some charsets like "us-ascii" return nil encoding (treated as UTF-8 compatible)
		return data, nil
	}

	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())
	return io.ReadAll(reader)
}
