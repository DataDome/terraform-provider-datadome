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

		assert.Empty(t, diags)
		assert.NotNil(t, meta)

		config, ok := meta.(*ProviderConfig)
		assert.True(t, ok, "meta should be of type *ProviderConfig")
		assert.NotNil(t, config.ClientCustomRule)
		assert.NotNil(t, config.ClientEndpoint)
		clientCustomRule := config.ClientCustomRule.(*datadome.ClientCustomRule)
		assert.Equal(t, "", clientCustomRule.Token)
		clientEndpoint := config.ClientEndpoint.(*datadome.ClientEndpoint)
		assert.Equal(t, "", clientEndpoint.Token)
	})

	t.Run("With custom host (direct)", func(t *testing.T) {
		host := "custom_host"
		rd := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
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
		host := "custom_host"
		err := os.Setenv("DATADOME_HOST", host)
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

func testAccCustomRuleResourcePreCheck(t *testing.T) {}

// TestAccCustomRuleResource_basic test the creation and the read of a new custom rule
func TestAccCustomRuleResource_basic(t *testing.T) {
	mockClient := datadome.NewMockClientCustomRule()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientCustomRule: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCustomRuleResourcePreCheck(t) },
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
		PreCheck:          func() { testAccCustomRuleResourcePreCheck(t) },
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
		PreCheck:          func() { testAccCustomRuleResourcePreCheck(t) },
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
		PreCheck:          func() { testAccCustomRuleResourcePreCheck(t) },
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
		PreCheck:          func() { testAccCustomRuleResourcePreCheck(t) },
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
		PreCheck:          func() { testAccCustomRuleResourcePreCheck(t) },
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

const testAccEndpointConfigUpdate = `
provider "datadome" {}
resource "datadome_endpoint" "simple" {
  cookie_same_site     = "Lax"
  description          = "This is a test"
  detection_enabled    = false
  name                 = "test-terraform-updated"
  protection_enabled   = false
  response_format      = "auto"
  source               = "Mobile App"
  traffic_usage        = "Account Creation"
  user_agent_inclusion = "TFTEST"
}
`

func TestAccEndpointResource_update(t *testing.T) {
	mockClient := datadome.NewMockClientEndpoint()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return &ProviderConfig{
			ClientEndpoint: mockClient,
		}, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccCustomRuleResourcePreCheck(t) },
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
			{
				Config: testAccEndpointConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists("datadome_endpoint.simple"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "name", "test-terraform-updated"),
					resource.TestCheckResourceAttr("datadome_endpoint.simple", "source", "Mobile App"),
				),
			},
		},
	})
}
