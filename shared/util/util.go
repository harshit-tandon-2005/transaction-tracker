package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/coin-tracker/transaction-tracker/shared/constants"
)

/*
Convert a string to an int64
*/
func StringToInt(inp string) (int64, error) {

	res, err := strconv.ParseInt(inp, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp string format '%s': %w", inp, err)
	}
	return res, nil

}

/*
Convert a unix timestamp string to a formatted string
*/
func FormatUnixTimestampString(unixTimestampStr string) (string, error) {
	// 1. Parse the string timestamp to int64
	unixTimestampInt, err := StringToInt(unixTimestampStr)
	if err != nil {
		return "", err
	}

	// Handle potential zero timestamp (if needed)
	if unixTimestampInt == 0 {
		// Decide how to handle zero: return error, empty string, or a default?
		// return "", fmt.Errorf("zero timestamp provided")
		return "00-00-0000 00:00:00", nil // Or return empty string ""
	}

	// 2. Convert int64 Unix timestamp to time.Time
	// time.Unix takes seconds and nanoseconds. Assume the input is in seconds.
	t := time.Unix(unixTimestampInt, 0).UTC()

	// 3. Format the time.Time object
	// Go uses a specific reference date (Mon Jan 2 15:04:05 MST 2006)
	// to define the desired output format.
	// DD -> 02
	// MM -> 01
	// YYYY -> 2006
	// HH -> 15 (for 24-hour clock)
	// MM -> 04 (for minutes)
	// SS -> 05 (for seconds)
	formattedString := t.Format(constants.DATE_FORMAT_YYYY_MM_DD_HH_MM_SS)

	return formattedString, nil
}

/*
Get the current working directory
*/
func GetCurrentWorkingDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %w", err)
	}
	return cwd, nil
}

/*
A function called WriteCSV is generic and works with a type parameter that is T."
The type T can be absolutely any type (any constraint).

"The function takes a filePath and a slice of type T."
"The function will write the slice of type T to a CSV file."
"The function will use the field tags to determine the headers of the CSV file."
"The function will use the field tags to determine the order of the fields in the CSV file."
"The function will use the field tags to determine the type of the fields in the CSV file."
"The function will use the field tags to determine the format of the fields in the CSV file."
*/
func WriteCSV[T any](filePath string, data []T) error {
	if len(data) == 0 {
		fmt.Printf("Info: No data provided to WriteCSV for file: %s. Creating empty file with headers (if any).\n", filePath)
		// Allow creating an empty file with only headers if needed, or return an error:
		// return fmt.Errorf("no data provided to write to CSV")
	}

	// --- 1. Reflect to get headers and field indices ---
	var headers []string
	var fieldIndices []int

	// Need at least one element to reflect type, or handle empty data case separately
	if len(data) > 0 {
		dataType := reflect.TypeOf(data[0])
		if dataType.Kind() != reflect.Struct {
			return fmt.Errorf("input data must be a slice of structs, got %s", dataType.Kind())
		}

		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			tag := field.Tag.Get("csv") // Get the value associated with the "csv" key

			// Only include fields that have the 'csv' tag and it's not "-"
			if tag != "" && tag != "-" {
				headers = append(headers, tag)         // Add header name from tag
				fieldIndices = append(fieldIndices, i) // Add index of the field
			}
		}
	} else {
		// Attempt to get type from slice definition even if empty, requires Go 1.18+ generics knowledge
		sliceType := reflect.TypeOf(data)
		if sliceType.Kind() == reflect.Slice {
			elemType := sliceType.Elem()
			if elemType.Kind() == reflect.Struct {
				for i := 0; i < elemType.NumField(); i++ {
					field := elemType.Field(i)
					tag := field.Tag.Get("csv")
					if tag != "" && tag != "-" {
						headers = append(headers, tag)
						fieldIndices = append(fieldIndices, i)
					}
				}
			} else {
				return fmt.Errorf("input data slice element is not a struct, got %s", elemType.Kind())
			}
		} else {
			return fmt.Errorf("input data is not a slice")
		}

		if len(headers) == 0 {
			fmt.Printf("Warning: No fields with 'csv' tags found in struct type for file %s. CSV will be empty.\n", filePath)
			// Optionally return an error here if headers are mandatory
		}
	}

	// --- 2. Ensure directory exists ---
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil { // os.ModePerm (0777) might be too permissive, adjust if needed
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// --- 3. Create and open file ---
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close() // Ensure file is closed

	// --- 4. Create CSV writer ---
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all buffered data is written to the file

	// --- 5. Write Header Row ---
	if len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write CSV header to %s: %w", filePath, err)
		}
	} else if len(data) > 0 {
		// We have data but no headers, which is weird.
		return fmt.Errorf("data provided but no fields found with 'csv' tag")
	}

	// --- 6. Write Data Rows ---
	for _, item := range data {
		itemValue := reflect.ValueOf(item) // Get the value representation of the struct
		if itemValue.Kind() != reflect.Struct {
			// This check might be redundant with generics but good for safety
			return fmt.Errorf("encountered non-struct item in data slice")
		}

		var record []string
		for _, index := range fieldIndices {
			fieldValue := itemValue.Field(index) // Get the value of the field by index
			// Convert field value to string
			record = append(record, valueToString(fieldValue))
		}

		// Write the record to CSV
		if err := writer.Write(record); err != nil {
			// Log the problematic record?
			return fmt.Errorf("failed to write record %+v to CSV file %s: %w", record, filePath, err)
		}
	}

	// --- 7. Check for writer errors after flush ---
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error occurred during CSV writing/flushing for %s: %w", filePath, err)
	}

	fmt.Printf("Successfully wrote %d data rows to CSV file: %s\n", len(data), filePath)
	return nil // Success
}

// valueToString converts a reflect.Value to its string representation for CSV.
// Add more type cases as needed (e.g., pointers, custom types).
func valueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		// Use 'f', -1 precision (automatic), 64 bit representation
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Struct:
		// Handle time.Time specifically for better formatting
		if t, ok := v.Interface().(time.Time); ok {
			// Use a consistent format, maybe the one from constants or FormatUnixTimestampString logic
			return t.Format(constants.DATE_FORMAT_YYYY_MM_DD_HH_MM_SS) // Example: Use constant
			// Or reuse your existing function if appropriate (though it expects Unix timestamp)
			// return FormatUnixTimestampString(strconv.FormatInt(t.Unix(), 10)) ?? careful here
		}
		// Handle other structs? Return default string or error?
		return fmt.Sprintf("%v", v.Interface()) // Default representation
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "" // Represent nil pointers/interfaces as empty string
		}
		// Try converting the element it points to / holds
		return valueToString(v.Elem())
	default:
		// Fallback for other types
		if v.IsValid() {
			return fmt.Sprintf("%v", v.Interface())
		}
		return "" // Return empty for invalid/zero values
	}
}
