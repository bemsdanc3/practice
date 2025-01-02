package validation

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"os"
)

func ValidateJSONWithSchemaFile(jsonData []byte, schemaPath string) error {
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	documentLoader := gojsonschema.NewBytesLoader(jsonData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed to validate JSON: %w", err)
	}

	if !result.Valid() {
		var validationErrors string
		for _, desc := range result.Errors() {
			validationErrors += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("validation failed:\n%s", validationErrors)
	}

	return nil
}
