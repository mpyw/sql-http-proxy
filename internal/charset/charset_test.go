package charset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToUTF8_UTF8Passthrough(t *testing.T) {
	data := []byte("hello world")

	// utf-8
	result, err := ToUTF8(data, "utf-8")
	require.NoError(t, err)
	assert.Equal(t, data, result)

	// UTF-8 (uppercase)
	result, err = ToUTF8(data, "UTF-8")
	require.NoError(t, err)
	assert.Equal(t, data, result)

	// utf8 (no hyphen)
	result, err = ToUTF8(data, "utf8")
	require.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestToUTF8_ShiftJIS(t *testing.T) {
	// "テスト" in Shift_JIS
	shiftJISData := []byte{0x83, 0x65, 0x83, 0x58, 0x83, 0x67}
	expected := "テスト"

	result, err := ToUTF8(shiftJISData, "shift_jis")
	require.NoError(t, err)
	assert.Equal(t, expected, string(result))
}

func TestToUTF8_ISO8859_1(t *testing.T) {
	// "café" in ISO-8859-1 (Latin-1)
	latin1Data := []byte{0x63, 0x61, 0x66, 0xe9}
	expected := "café"

	result, err := ToUTF8(latin1Data, "iso-8859-1")
	require.NoError(t, err)
	assert.Equal(t, expected, string(result))
}

func TestToUTF8_UnsupportedCharset(t *testing.T) {
	data := []byte("hello")

	_, err := ToUTF8(data, "unknown-charset-xyz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported charset")
}

func TestToUTF8_NilEncodingPassthrough(t *testing.T) {
	// Some charsets like "csUnicode" return nil encoding
	// They should be treated as UTF-8 compatible
	data := []byte("hello")

	result, err := ToUTF8(data, "csUnicode")
	require.NoError(t, err)
	assert.Equal(t, data, result)
}
