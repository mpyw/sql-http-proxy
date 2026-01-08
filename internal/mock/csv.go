package mock

import (
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
)

// CSVSource holds pre-parsed CSV data.
type CSVSource struct {
	rows []map[string]any
}

// ParseCSV parses inline CSV data into a CSVSource.
func ParseCSV(data string) (*CSVSource, error) {
	return parseCSVReader(strings.NewReader(data))
}

// ParseCSVFile parses a CSV file into a CSVSource.
func ParseCSVFile(path string) (*CSVSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return parseCSVReader(f)
}

func parseCSVReader(r io.Reader) (*CSVSource, error) {
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
				row[header] = parseValue(record[i])
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
// ctx and sql are ignored for static data sources.
func (s *CSVSource) Data(_ map[string]any, _ string, _ map[string]any) (any, map[string]any, error) {
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
