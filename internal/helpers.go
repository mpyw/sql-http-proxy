package internal

import (
	"io"
	"log/slog"
)

// CloseQuietly closes c, demoting a close failure to a debug log. For
// read-side closes where the data is already in hand and nothing can act on
// the error.
func CloseQuietly(c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Debug("Failed to close file", "error", err)
	}
}

func SliceOf[T any](items ...T) []T {
	return items
}
