package datadome

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/datadome/terraform-provider/datadome-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

var (
	testAccProviders map[string]func() (*schema.Provider, error)
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"datadome": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

/*
Provider tests
*/

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProviderConfigure(t *testing.T) {
	t.Run("With apiKey (direct)", func(t *testing.T) {
		apiKey := "valid_api_key"
		rd := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
			"apikey": apiKey,
		})

		meta, diags := providerConfigure(context.Background(), rd)

		assert.Empty(t, diags)
		assert.NotNil(t, meta)

		config, ok := meta.(*ProviderConfig)
		assert.True(t, ok, "meta should be of type *ProviderConfig")
		assert.NotNil(t, config.ClientCustomRule)
		assert.NotNil(t, config.ClientEndpoint)
		clientCustomRule := config.ClientCustomRule.(*datadome.ClientCustomRule)
		assert.Equal(t, apiKey, clientCustomRule.Token)
		clientEndpoint := config.ClientEndpoint.(*datadome.ClientEndpoint)
		assert.Equal(t, apiKey, clientEndpoint.Token)
	})

	t.Run("With apiKey (env)", func(t *testing.T) {
		apiKey := "valid_api_key"
		err := os.Setenv("DATADOME_APIKEY", apiKey)
		if err != nil {
			t.Fatalf("fail to set DATADOME_APIKEY with value %q", apiKey)
			return
		}
		defer os.Unsetenv("DATADOME_APIKEY")

		rd := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})

		meta, diags := providerConfigure(context.Background(), rd)

		assert.Empty(t, diags)
		assert.NotNil(t, meta)

		config, ok := meta.(*ProviderConfig)
		assert.True(t, ok, "meta should be of type *ProviderConfig")
		assert.NotNil(t, config.ClientCustomRule)
		assert.NotNil(t, config.ClientEndpoint)
		clientCustomRule := config.ClientCustomRule.(*datadome.ClientCustomRule)
		assert.Equal(t, apiKey, clientCustomRule.Token)
		clientEndpoint := config.ClientEndpoint.(*datadome.ClientEndpoint)
		assert.Equal(t, apiKey, clientEndpoint.Token)
	})

	t.Run("Without apiKey", func(t *testing.T) {
		rd := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})

		meta, diags := providerConfigure(context.Background(), rd)

		assert.NotNil(t, diags)
		assert.Nil(t, meta)

		assert.Len(t, diags, 1, "Expected one diag error")
		assert.Equal(t, diag.Error, diags[0].Severity)
		assert.Equal(t, diags[0].Summary, "Missing required 'apikey' value")
		assert.Equal(t, diags[0].Detail, "The 'apikey' field is required but not set.")
	})

	t.Run("With custom host (direct)", func(t *testing.T) {
		// Set required value
		apiKey := "valid_api_key"
		host := "custom_host"
		rd := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
			"apikey": apiKey,
			"host": host,
		})

		meta, diags := providerConfigure(context.Background(), rd)

		assert.Empty(t, diags)
		assert.NotNil(t, meta)

		config, ok := meta.(*ProviderConfig)
		assert.True(t, ok, "meta should be of type *ProviderConfig")
		assert.NotNil(t, config.ClientCustomRule)
		assert.NotNil(t, config.ClientEndpoint)
		clientCustomRule := config.ClientCustomRule.(*datadome.ClientCustomRule)
		assert.Equal(t, host, clientCustomRule.HostURL)
		clientEndpoint := config.ClientEndpoint.(*datadome.ClientEndpoint)
		assert.Equal(t, host, clientEndpoint.HostURL)
	})

	t.Run("With custom host (env)", func(t *testing.T) {
		// Set required value
		apiKey := "valid_api_key"
		err := os.Setenv("DATADOME_APIKEY", apiKey)
		if err != nil {
			t.Fatalf("fail to set DATADOME_APIKEY with value %q", apiKey)
			return
		}
		defer os.Unsetenv("DATADOME_APIKEY")

		host := "custom_host"
		err = os.Setenv("DATADOME_HOST", host)
		if err != nil {
			t.Fatalf("fail to set DATADOME_HOST with value %q", host)
			return
		}
		defer os.Unsetenv("DATADOME_HOST")

		rd := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})

		meta, diags := providerConfigure(context.Background(), rd)

		assert.Empty(t, diags)
		assert.NotNil(t, meta)

		config, ok := meta.(*ProviderConfig)
		assert.True(t, ok, "meta should be of type *ProviderConfig")
		assert.NotNil(t, config.ClientCustomRule)
		assert.NotNil(t, config.ClientEndpoint)
		clientCustomRule := config.ClientCustomRule.(*datadome.ClientCustomRule)
		assert.Equal(t, host, clientCustomRule.HostURL)
		clientEndpoint := config.ClientEndpoint.(*datadome.ClientEndpoint)
		assert.Equal(t, host, clientEndpoint.HostURL)
	})
}

/*
Resources test helpers
*/

func testAccResourcePreCheck(t *testing.T) {}

// testAccCheckResourceExists check if the given resourceName exists
func testAccCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %q", resourceName)
		}

		return nil
	}
}

// testAccCheckCustomRuleResourceExists check if the given resourceName does not exists
func testAccCheckResourceDoesNotExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("resource still exists: %q", resourceName)
		}

		return nil
	}
}

/*
Resources CustomRules tests
*/

const testAccCustomRuleResourceConfig = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = "acc-test"
  query         = "ip: 192.168.0.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "low"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigEmptyName = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = ""
  query         = "ip: 192.168.0.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "low"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigBlankName = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = "           "
  query         = "ip: 192.168.0.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "low"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigEmptyQuery = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = "acc-test"
  query         = ""
  response      = "allow"
  endpoint_type = "web"
  priority      = "low"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigWrongResponse = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = "acc-test"
  query         = "ip: 192.168.0.1"
  response      = "wrong_response"
  endpoint_type = "web"
  priority      = "low"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigWrongEndpoint = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = "acc-test"
  query         = "ip: 192.168.0.1"
  response      = "allow"
  endpoint_type = "wrong_endpoint"
  priority      = "low"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigWrongPriority = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = "acc-test"
  query         = "ip: 192.168.0.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "wrong_priority"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigUpdate = `
provider "datadome" {}

resource "datadome_custom_rule" "accConfig" {
  name          = "acc-test-updated"
  query         = "ip: 192.168.0.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "normal"
  enabled		= true
}
`

const testAccCustomRuleResourceConfigEmpty = `
provider "datadome" {}
`

// TestAccCustomRuleResource_basic test the creation and the read of a new custom rule
func TestAccCustomRuleResource_basic(t *testing.T) {
	mockClient := datadome.NewMockClientCustomRule()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientCustomRule: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRuleResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_custom_rule.accConfig"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "name", "acc-test"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "query", "ip: 192.168.0.1"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "response", "allow"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "endpoint_type", "web"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "priority", "low"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "enabled", "true"),
				),
			},
		},
	})
}

// TestAccCustomRuleResource_update test the creation of a new custom rule and update it
func TestAccCustomRuleResource_update(t *testing.T) {
	mockClient := datadome.NewMockClientCustomRule()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientCustomRule: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRuleResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_custom_rule.accConfig"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "name", "acc-test"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "query", "ip: 192.168.0.1"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "response", "allow"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "endpoint_type", "web"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "priority", "low"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "enabled", "true"),
				),
			},
			{
				Config: testAccCustomRuleResourceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_custom_rule.accConfig"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "name", "acc-test-updated"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "priority", "normal"),
				),
			},
		},
	})
}

// TestAccCustomRuleResource_delete test the creation of a new custom rule and delete it
func TestAccCustomRuleResource_delete(t *testing.T) {
	mockClient := datadome.NewMockClientCustomRule()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientCustomRule: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRuleResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_custom_rule.accConfig"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "name", "acc-test"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "query", "ip: 192.168.0.1"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "response", "allow"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "endpoint_type", "web"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "priority", "low"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "enabled", "true"),
				),
			},
			{
				Config: testAccCustomRuleResourceConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceDoesNotExists("datadome_custom_rule.accConfig"),
				),
			},
		},
	})
}

// TestAccCustomRuleResource_wrongParameters test the creation with wrong parameters (i.e. response, endpoint, and priority)
func TestAccCustomRuleResource_wrongParameters(t *testing.T) {
	mockClient := datadome.NewMockClientCustomRule()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientCustomRule: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomRuleResourceConfigWrongResponse,
				ExpectError: regexp.MustCompile(`"wrong_response" is not an acceptable response`),
			},
			{
				Config:      testAccCustomRuleResourceConfigWrongEndpoint,
				ExpectError: regexp.MustCompile(`"wrong_endpoint" is not an acceptable endpoint`),
			},
			{
				Config:      testAccCustomRuleResourceConfigWrongPriority,
				ExpectError: regexp.MustCompile(`"wrong_priority" is not an acceptable priority`),
			},
			{
				Config:      testAccCustomRuleResourceConfigEmptyName,
				ExpectError: regexp.MustCompile(`the name value should not be blank`),
			},
			{
				Config:      testAccCustomRuleResourceConfigBlankName,
				ExpectError: regexp.MustCompile(`the name value should not be blank`),
			},
			{
				Config:      testAccCustomRuleResourceConfigEmptyQuery,
				ExpectError: regexp.MustCompile(`expected "query" to not be an empty string`),
			},
		},
	})
}

// TestAccCustomRuleResource_createAlreadyExists test the creation when a custom rule already exists with the same name
func TestAccCustomRuleResource_createAlreadyExists(t *testing.T) {
	mockClient := datadome.NewMockClientCustomRule()
	mockClient.CreateFunc = func(ctx context.Context, params datadome.CustomRule) (*int, error) {
		return nil, fmt.Errorf("The rule with name: 'acc-test' already exists")
	}

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientCustomRule: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCustomRuleResourceConfig,
				ExpectError: regexp.MustCompile(`'acc-test' already exists`),
			},
		},
	})
}

// TestAccCustomRuleResource_updateAlreadyExists test the update when a custom rule already exists with the same name
func TestAccCustomRuleResource_updateAlreadyExists(t *testing.T) {
	mockClient := datadome.NewMockClientCustomRule()
	mockClient.UpdateFunc = func(ctx context.Context, params datadome.CustomRule) (*datadome.CustomRule, error) {
		return nil, fmt.Errorf("Another rule with name: 'acc-test-updated' already exists")
	}

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientCustomRule: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRuleResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_custom_rule.accConfig"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "name", "acc-test"),
				),
			},
			{
				Config:      testAccCustomRuleResourceConfigUpdate,
				ExpectError: regexp.MustCompile(`'acc-test-updated' already exists`),
			},
		},
	})
}

/*
Resources Endpoints tests
*/

const testAccEndpointConfig = `
provider "datadome" {}

resource "datadome_endpoint" "simple" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigWithRegex = `
provider "datadome" {}

resource "datadome_endpoint" "simple" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "(?:BLOCK|CHALLENGE)UA"
}
`

const testAccEndpointConfigWithoutOptionalFields = `
provider "datadome" {}

resource "datadome_endpoint" "simple" {
  name                 = "test-terraform"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigWithPositionBefore = `
provider "datadome" {}

resource "datadome_endpoint" "first" {
  name                 = "test-terraform"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
}

resource "datadome_endpoint" "second" {
  name                 = "test-terraform"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
  position_before      = datadome_endpoint.first.id
}
`

const testAccEndpointConfigMissingFields = `
provider "datadome" {}

resource "datadome_endpoint" "missing_fields" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
}
`

const testAccEndpointConfigWrongResponseFormat = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_response_format" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "wrong_response_format"
  source               = "Web Browser"
  user_agent_inclusion = "TFTEST"
  traffic_usage        = "Account Creation"
}
`

const testAccEndpointConfigWrongSource = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_source" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "wrong_source"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigWrongTrafficUsage = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_traffic_usage" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Web Browser"
  traffic_usage        = "wrong_traffic_usage"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigWrongTrafficUsageWithSourceApi = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_traffic_usage" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Api"
  traffic_usage        = "Rss"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigWrongTrafficUsageWithSourceMobileApp = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_traffic_usage" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Mobile App"
  traffic_usage        = "Rss"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigWrongCookieSameSite = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_cookie_same_site" {
  cookie_same_site     = "wrong_cookie_same_site"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigWrongPositionBeforeFormat = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_cookie_same_site" {
  cookie_same_site     = "Lax"
  position_before      = "some_incorrect_id"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
}
`

const testAccEndpointConfigInvalidRegex = `
provider "datadome" {}

resource "datadome_endpoint" "wrong_cookie_same_site" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Web Browser"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "wrong(.*"
}
`

// TestAccEndpointResource_basic tests the creation and the read of a new endpoint
func TestAccEndpointResource_basic(t *testing.T) {
	mockClient := datadome.NewMockClientEndpoint()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientEndpoint: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_endpoint.simple"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "cookie_same_site", "Lax"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "description", "This is a test"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "detection_enabled", "false"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "name", "test-terraform"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "protection_enabled", "false"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "response_format", "auto"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "source", "Web Browser"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "traffic_usage", "Account Creation"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "user_agent_inclusion", "TFTEST"),
				),
			},
		},
	})
}

// TestAccEndpointResource_createWithoutOptionalFields tests the creation of an endpoint resource without optional fields
func TestAccEndpointResource_createWithoutOptionalFields(t *testing.T) {
	mockClient := datadome.NewMockClientEndpoint()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientEndpoint: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigWithoutOptionalFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_endpoint.simple"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "cookie_same_site", "Lax"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "detection_enabled", "true"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "name", "test-terraform"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "protection_enabled", "false"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "response_format", "auto"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "source", "Web Browser"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "traffic_usage", "Account Creation"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "user_agent_inclusion", "TFTEST"),
				),
			},
		},
	})
}

// TestAccEndpointResource_createWithPositionBefore tests the creation of two resources
// The first resource is created with the minimum required fields
// The second resource is created with the minimum required fields and specify the `position_before` to use the ID of the first resource
func TestAccEndpointResource_createWithPositionBefore(t *testing.T) {
	mockClient := datadome.NewMockClientEndpoint()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientEndpoint: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigWithPositionBefore,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_endpoint.first"),
					testAccCheckResourceExists("datadome_endpoint.second"),
					resource.TestCheckResourceAttrSet("datadome_endpoint.first", "id"),
					resource.TestCheckResourceAttrSet("datadome_endpoint.second", "id"),
					resource.TestCheckResourceAttrSet("datadome_endpoint.second", "position_before"),
					resource.TestCheckResourceAttrPair(
						"datadome_endpoint.second", "position_before",
						"datadome_endpoint.first", "id",
					),
				),
			},
		},
	})
}

// TestAccEndpointResource_createWithRegex tests the creation of a new endpoint with a Regex format for the "user_agent_inclusion" field
func TestAccEndpointResource_createWithRegex(t *testing.T) {
	mockClient := datadome.NewMockClientEndpoint()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientEndpoint: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigWithRegex,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_endpoint.simple"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "cookie_same_site", "Lax"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "description", "This is a test"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "detection_enabled", "false"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "name", "test-terraform"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "protection_enabled", "false"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "response_format", "auto"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "source", "Web Browser"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "traffic_usage", "Account Creation"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "user_agent_inclusion", "(?:BLOCK|CHALLENGE)UA"),
				),
			},
		},
	})
}

// TestAccEndpointResource_wrongParameters tests the creation of an endpoint resource by providing wrong inputs
func TestAccEndpointResource_wrongParameters(t *testing.T) {
	mockClient := datadome.NewMockClientEndpoint()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientEndpoint: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccResourcePreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEndpointConfigMissingFields,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccEndpointConfigWrongResponseFormat,
				ExpectError: regexp.MustCompile(`"wrong_response_format" is not an acceptable response_format`),
			},
			{
				Config:      testAccEndpointConfigWrongSource,
				ExpectError: regexp.MustCompile(`"wrong_source" is not an acceptable source`),
			},
			{
				Config:      testAccEndpointConfigWrongTrafficUsage,
				ExpectError: regexp.MustCompile(`"wrong_traffic_usage" is not an acceptable traffic_usage`),
			},
			{
				Config:      testAccEndpointConfigWrongTrafficUsageWithSourceApi,
				ExpectError: regexp.MustCompile(`expected "traffic_usage" to be one of {General}, got "Rss"`),
			},
			{
				Config:      testAccEndpointConfigWrongTrafficUsageWithSourceMobileApp,
				ExpectError: regexp.MustCompile(`expected "traffic_usage" to be one of {General, Login, Payment, Cart, Forms, Account Creation}, got "Rss"`),
			},
			{
				Config:      testAccEndpointConfigWrongCookieSameSite,
				ExpectError: regexp.MustCompile(`"wrong_cookie_same_site" is not an acceptable cookie_same_site`),
			},
			{
				Config:      testAccEndpointConfigWrongPositionBeforeFormat,
				ExpectError: regexp.MustCompile(`expected "position_before" to be a valid UUID, got some_incorrect_id`),
			},
			{
				Config:      testAccEndpointConfigInvalidRegex,
				ExpectError: regexp.MustCompile(`.*error parsing regexp.*`),
			},
		},
	})
}
