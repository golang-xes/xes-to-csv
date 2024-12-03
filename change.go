package xes_to_csv

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// XES represents the structure of the XES file.
type XES struct {
	XMLName xml.Name `xml:"log"`
	Traces  []Trace  `xml:"trace"` // Changed variable name to plural form for consistency.
}

// Trace represents a single trace in the XES file.
type Trace struct {
	Events           []Event           `xml:"event"` // Changed variable name to plural form for consistency.
	StringAttributes []StringAttribute `xml:"string"`
}

// Event represents a single event within a trace in the XES file.
type Event struct {
	StringAttributes []StringAttribute `xml:"string"`
	DateAttributes   []DateAttribute   `xml:"date"`
}

// StringAttribute represents a string attribute in an event or trace.
type StringAttribute struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

// DateAttribute represents a date attribute in an event.
type DateAttribute struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}

// ConvertXESToCSV reads an XES file and writes its contents to a CSV file.
func ConvertXESToCSV(XESFilePath, CSVFilePath string) error {
	// Validate and clean the input file path.
	inputPath := filepath.Clean(XESFilePath)
	if !isValidXESFile(inputPath) {
		return fmt.Errorf("input file must be an XES file: %s", inputPath)
	}

	// Open the XES file.
	xesFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open XES file: %w", err)
	}
	defer func(xesFile *os.File) {
		_ = xesFile.Close()
	}(xesFile) // Simplified defer statement.

	// Parse the XML data.
	xes := XES{}
	decoder := xml.NewDecoder(xesFile)
	if err := decoder.Decode(&xes); err != nil {
		return fmt.Errorf("failed to decode XES file: %w", err)
	}

	// Collect all unique attribute keys.
	keyMap := collectAttributeKeys(xes)

	// Prepare the CSV file.
	csvFilePath := filepath.Clean(CSVFilePath)
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer func(csvFile *os.File) {
		_ = csvFile.Close()
	}(csvFile) // Simplified defer statement.

	// Write UTF-8 BOM to ensure correct encoding.
	if _, err := csvFile.WriteString("\xEF\xBB\xBF"); err != nil {
		return fmt.Errorf("failed to write UTF-8 BOM: %w", err)
	}

	// Initialize CSV writer and write the header.
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	if err := writer.Write(collectHeader(keyMap)); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write XES data to CSV file.
	return writeXESToCSV(writer, xes, keyMap)
}

// isValidXESFile checks if the given file path has a .xes extension.
func isValidXESFile(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".xes")
}

// collectAttributeKeys gathers all unique attribute keys from the XES structure.
func collectAttributeKeys(xes XES) map[string]struct{} {
	keyMap := make(map[string]struct{})
	for _, trace := range xes.Traces {
		for _, attr := range trace.StringAttributes {
			if attr.Key == "concept:name" {
				keyMap["case:concept:name"] = struct{}{}
			} else {
				keyMap[attr.Key] = struct{}{}
			}
		}
		for _, event := range trace.Events {
			for _, attr := range event.StringAttributes {
				keyMap[attr.Key] = struct{}{}
			}
			for _, attr := range event.DateAttributes {
				keyMap[attr.Key] = struct{}{}
			}
		}
	}
	return keyMap
}

// collectHeader creates the CSV header from the attribute keys.
func collectHeader(keyMap map[string]struct{}) []string {
	var keys []string
	for key := range keyMap {
		keys = append(keys, strings.TrimSpace(key))
	}
	return keys
}

// writeXESToCSV writes the content of the XES file to the CSV file.
func writeXESToCSV(writer *csv.Writer, xes XES, keyMap map[string]struct{}) error {
	keys := collectHeader(keyMap)
	for _, trace := range xes.Traces {
		for _, event := range trace.Events {
			record := make([]string, len(keys))
			for _, attr := range event.StringAttributes {
				setAttributeValue(record, keys, attr)
			}
			for _, attr := range event.DateAttributes {
				setAttributeValue(record, keys, attr)
			}
			for _, attr := range trace.StringAttributes {
				if attr.Key == "concept:name" {
					attr.Key = "case:concept:name"
				}
				setAttributeValue(record, keys, attr)
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write CSV record: %w", err)
			}
		}
	}
	return nil
}

// setAttributeValue sets the value of an attribute in the CSV record.
func setAttributeValue(record []string, keys []string, attr interface{}) {
	switch v := attr.(type) {
	case StringAttribute:
		index := findIndex(keys, v.Key)
		if index != -1 {
			record[index] = strings.TrimSpace(v.Value)
		}
	case DateAttribute:
		index := findIndex(keys, v.Key)
		if index != -1 {
			record[index] = strings.TrimSpace(v.Value)
		}
	}
}

// findIndex finds the index of a key in the keys slice.
func findIndex(keys []string, key string) int {
	for i, k := range keys {
		if k == key {
			return i
		}
	}
	return -1
}
