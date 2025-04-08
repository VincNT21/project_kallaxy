package gui

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/widget"
	"github.com/VincNT21/kallaxy/client/context"
)

func createMetadataForm(appCtxt *context.AppContext, mediaType string) (*widget.Form, map[string]*widget.Entry) {

	// Create var
	form := &widget.Form{}
	entryMap := make(map[string]*widget.Entry)

	// Get metadata field for this media type
	fields := appCtxt.MetadataFieldsMap[mediaType]

	for _, field := range fields {
		var entry *widget.Entry
		if field == "overview" || field == "description" {
			entry = widget.NewMultiLineEntry()
		} else {
			entry = widget.NewEntry()
		}
		form.Append(field, entry)
		entryMap[field] = entry
	}

	return form, entryMap
}

func extractMetadataValues(appCtxt *context.AppContext, entryMap map[string]*widget.Entry) map[string]interface{} {

	result := make(map[string]interface{})

	for field, value := range entryMap {
		inputText := value.Text
		result[field] = processMetadataFieldValue(appCtxt, field, inputText)
	}

	return result
}

func processMetadataFieldValue(appCtxt *context.AppContext, fieldName, textValue string) interface{} {
	// Define type specifications for each metadata field
	spec, exists := appCtxt.MetadataFieldsSpecs[fieldName]
	if !exists {
		// Default to string if not specified
		return textValue
	}

	switch spec.FieldType {
	case "list":
		// For list, split and trim the text input and add it to a list
		items := []string{}
		if textValue != "" {
			for _, item := range strings.Split(textValue, ",") {
				trimmed := strings.TrimSpace(item)
				if trimmed != "" {
					items = append(items, trimmed)
				}
			}
		}
		return items
	case "int":
		// For int, convert trimmed string to int
		if val, err := strconv.Atoi(strings.TrimSpace(textValue)); err != nil {
			return val
		}
		return 0 // Default value for int
	case "string":
		return textValue // No special processing for strings
	default:
		return textValue
	}
}

func formatMetadataValueForEntry(appCtxt *context.AppContext, fieldName string, value interface{}) string {
	// Define type specifications for each metadata field
	spec, exists := appCtxt.MetadataFieldsSpecs[fieldName]
	if !exists {
		// Default to string if not specified
		if str, ok := value.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", value)
	}

	switch spec.FieldType {
	case "list":
		// Chech if it's already the right type
		if strList, ok := value.([]string); ok {
			return strings.Join(strList, ", ")
		} else if list, ok := value.([]interface{}); ok {
			// Handle case where JSON Unmarshalling gives []interface{}
			strList := make([]string, len(list))
			for i, item := range list {
				strList[i] = fmt.Sprintf("%v", item)
			}
			return strings.Join(strList, ", ")
		}
		return "" // Default for invalid list

	case "int":
		switch v := value.(type) {
		case int:
			return strconv.Itoa(v)
		case float64:
			return strconv.Itoa(int(v))
		case string:
			return v // Already a string
		default:
			return fmt.Sprintf("%v", value)
		}

	case "string":
		if str, ok := value.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", value)

	default:
		return fmt.Sprintf("%v", value)
	}
}
