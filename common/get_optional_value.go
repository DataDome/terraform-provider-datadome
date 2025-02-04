package common

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// GetOptionalValue is a generic function that retrieve the expected field's value by its name.
// If the fields is set, it returns a pointer of this field.
// Otherwise, it returns a nil pointer.
func GetOptionalValue[T comparable](data *schema.ResourceData, field string) *T {
	var finalValue *T
	if value, ok := data.GetOk(field); ok {
		typedValue := value.(T)
		finalValue = &typedValue
	}

	return finalValue
}