package common

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// GetOptionalValue is a generic function that retrieve the expected field's value by its name.
// If the fields is set, it returns a pointer of this field.
// Otherwise, it returns a nil pointer.
func GetOptionalValue[T comparable](data *schema.ResourceData, field string) *T {
	//lint:ignore SA1019 GetOkExists is required for zero-value detection
	value, ok := data.GetOkExists(field) //nolint:staticcheck
	if !ok {
		return nil
	}

	typedValue, ok := value.(T)
	if !ok {
		return nil
	}

	return &typedValue
}
