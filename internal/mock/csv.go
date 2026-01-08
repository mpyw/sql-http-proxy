package mock

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// CSVSource holds pre-parsed CSV data.
type CSVSource struct {
	rows []map[string]any
}

// ParseCSVOptions contains options for CSV parsing.
type ParseCSVOptions struct {
	ValueParser *ValueParser
}

// ParseCSV parses inline CSV data into a CSVSource.
func ParseCSV(data string) (*CSVSource, error) {
	return ParseCSVWithOptions(data, ParseCSVOptions{})
}

// ParseCSVWithOptions parses inline CSV data with custom options.
func ParseCSVWithOptions(data string, opts ParseCSVOptions) (*CSVSource, error) {
	return parseCSVReaderWithOptions(strings.NewReader(data), opts)
}

// ParseCSVFile parses a CSV file into a CSVSource.
func ParseCSVFile(path string) (*CSVSource, error) {
	return ParseCSVFileWithOptions(path, ParseCSVOptions{})
}

// ParseCSVFileWithOptions parses a CSV file with custom options.
func ParseCSVFileWithOptions(path string, opts ParseCSVOptions) (*CSVSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return parseCSVReaderWithOptions(f, opts)
}

func parseCSVReaderWithOptions(r io.Reader, opts ParseCSVOptions) (*CSVSource, error) {
	reader := csv.NewReader(r)
	// Allow variable number of fields per record
	reader.FieldsPerRecord = -1

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("csv: empty data, header row required")
		}
		return nil, err
	}

	if len(headers) == 0 {
		return nil, errors.New("csv: header row is empty")
	}

	expectedCols := len(headers)
	var rows []map[string]any
	lineNum := 1 // Header is line 1

	for {
		lineNum++
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		row := make(map[string]any, expectedCols)
		actualCols := len(record)

		// Warn about column count mismatch
		if actualCols != expectedCols {
			if actualCols < expectedCols {
				slog.Warn("csv: row has fewer columns than header",
					"line", lineNum,
					"expected", expectedCols,
					"actual", actualCols,
				)
			} else {
				slog.Warn("csv: row has more columns than header, extra columns ignored",
					"line", lineNum,
					"expected", expectedCols,
					"actual", actualCols,
				)
			}
		}

		for i, header := range headers {
			if i < actualCols {
				var val any
				var err error
				if opts.ValueParser != nil {
					val, err = opts.ValueParser.Parse(record[i])
					if err != nil {
						return nil, fmt.Errorf("line %d, column %q: %w", lineNum, header, err)
					}
				} else {
					val = parseValue(record[i])
				}
				row[header] = val
			} else {
				// Missing column: use empty string
				row[header] = ""
			}
		}

		rows = append(rows, row)
	}

	return &CSVSource{rows: rows}, nil
}

// Data returns the parsed CSV rows.
// ctx, sql, and tc are ignored for static data sources.
func (s *CSVSource) Data(_ map[string]any, _ string, _ map[string]any, _ *js.TransformContext) (any, map[string]any, error) {
	// Return a copy to prevent mutation
	result := make([]map[string]any, len(s.rows))
	for i, row := range s.rows {
		rowCopy := make(map[string]any, len(row))
		for k, v := range row {
			rowCopy[k] = v
		}
		result[i] = rowCopy
	}
	return result, nil, nil
}
