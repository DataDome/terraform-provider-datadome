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
	testSchema := map[string]*schema.Schema{
		"optional_field": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	data := schema.TestResourceDataRaw(t, testSchema, map[string]interface{}{})

	result := GetOptionalValue[string](data, "optional_field")
	assert.Nil(t, result)
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
