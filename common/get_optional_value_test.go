package common

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestGetOptionalValue_String(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"optional_field": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{
		"optional_field": "hello",
	})

	result := GetOptionalValue[string](data, "optional_field")
	assert.NotNil(t, result)
	assert.Equal(t, "hello", *result)
}

func TestGetOptionalValue_Int(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"optional_field": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{
		"optional_field": 1234,
	})

	result := GetOptionalValue[int](data, "optional_field")
	assert.NotNil(t, result)
	assert.Equal(t, 1234, *result)
}

func TestGetOptionalValue_Bool(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"optional_field": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{
		"optional_field": true,
	})

	result := GetOptionalValue[bool](data, "optional_field")
	assert.NotNil(t, result)
	assert.Equal(t, true, *result)
}

func TestGetOptionalValue_MissingField(t *testing.T) {
	tests := []struct {
		name       string
		schemaType schema.ValueType
	}{
		{
			name:       "missing string field",
			schemaType: schema.TypeString,
		},
		{
			name:       "missing int field",
			schemaType: schema.TypeInt,
		},
		{
			name:       "missing bool field",
			schemaType: schema.TypeBool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSchema := map[string]*schema.Schema{
				"optional_field": {
					Type:     tt.schemaType,
					Optional: true,
				},
			}

			data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{})

			switch tt.schemaType {
			case schema.TypeString:
				result := GetOptionalValue[string](data, "optional_field")
				assert.Nil(t, result)
			case schema.TypeInt:
				result := GetOptionalValue[int](data, "optional_field")
				assert.Nil(t, result)
			case schema.TypeBool:
				result := GetOptionalValue[bool](data, "optional_field")
				assert.Nil(t, result)
			}
		})
	}
}

func TestGetOptionalValue_WrongTypeConversion(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"optional_field": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{
		"optional_field": 1234,
	})

	result := GetOptionalValue[string](data, "optional_field")
	assert.Nil(t, result)
}

func TestGetOptionalValue_ZeroValues(t *testing.T) {
	tests := []struct {
		name       string
		schemaType schema.ValueType
		value      interface{}
		expected   interface{}
	}{
		{
			name:       "zero string",
			schemaType: schema.TypeString,
		},
		{
			name:       "zero int",
			schemaType: schema.TypeInt,
			value:      0,
			expected:   0,
		},
		{
			name:       "zero bool",
			schemaType: schema.TypeBool,
			value:      false,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSchema := map[string]*schema.Schema{
				"optional_field": {
					Type:     tt.schemaType,
					Optional: true,
				},
			}

			data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{
				"optional_field": tt.value,
			})

			switch tt.schemaType {
			case schema.TypeString:
				result := GetOptionalValue[string](data, "optional_field")
				assert.Nil(t, result)
				t.Log("String type is the only one to return nil for zero value")
			case schema.TypeInt:
				result := GetOptionalValue[int](data, "optional_field")
				assert.NotNil(t, result)
				if result == nil {
					t.Fatalf("expected non-nil result for int type")
				}
				assert.Equal(t, tt.expected, *result)
			case schema.TypeBool:
				result := GetOptionalValue[bool](data, "optional_field")
				assert.NotNil(t, result)
				if result == nil {
					t.Fatalf("expected non-nil result for boolean type")
				}
				assert.Equal(t, tt.expected, *result)
			}
		})
	}
}

func TestGetOptionalValue_NullValues(t *testing.T) {
	tests := []struct {
		name       string
		schemaType schema.ValueType
	}{
		{
			name:       "null string",
			schemaType: schema.TypeString,
		},
		{
			name:       "null int",
			schemaType: schema.TypeInt,
		},
		{
			name:       "null bool",
			schemaType: schema.TypeBool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSchema := map[string]*schema.Schema{
				"optional_field": {
					Type:     tt.schemaType,
					Optional: true,
				},
			}

			data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{
				"optional_field": nil,
			})

			switch tt.schemaType {
			case schema.TypeString:
				result := GetOptionalValue[string](data, "optional_field")
				assert.Nil(t, result)
			case schema.TypeInt:
				result := GetOptionalValue[int](data, "optional_field")
				assert.Nil(t, result)
			case schema.TypeBool:
				result := GetOptionalValue[bool](data, "optional_field")
				assert.Nil(t, result)
			}
		})
	}
}
