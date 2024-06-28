package datadome

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/datadome/terraform-provider/datadome-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func testAccPreCheck(t *testing.T) {}

// testAccCheckCustomRuleResourceExists check if the given resourceName exists
func testAccCheckCustomRuleResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %q", resourceName)
		}

		return nil
	}
}

// testAccCheckCustomRuleResourceExists check if the given resourceName does not exists
func testAccCheckCustomRuleResourceDoesNotExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if ok {
			return fmt.Errorf("resource still exists: %q", resourceName)
		}

		return nil
	}
}

// TestAccCustomRuleResource_basic test the creation and the read of a new custom rule
func TestAccCustomRuleResource_basic(t *testing.T) {
	mockClient := datadome.NewMockClient()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return mockClient, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRuleResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomRuleResourceExists("datadome_custom_rule.accConfig"),
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
	mockClient := datadome.NewMockClient()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return mockClient, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRuleResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomRuleResourceExists("datadome_custom_rule.accConfig"),
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
					testAccCheckCustomRuleResourceExists("datadome_custom_rule.accConfig"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "name", "acc-test-updated"),
					resource.TestCheckResourceAttr("datadome_custom_rule.accConfig", "priority", "normal"),
				),
			},
		},
	})
}

// TestAccCustomRuleResource_delete test the creation of a new custom rule and delete it
func TestAccCustomRuleResource_delete(t *testing.T) {
	mockClient := datadome.NewMockClient()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return mockClient, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomRuleResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomRuleResourceExists("datadome_custom_rule.accConfig"),
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
					testAccCheckCustomRuleResourceDoesNotExists("datadome_custom_rule.accConfig"),
				),
			},
		},
	})
}

// TestAccCustomRuleResource_wrongParameters test the creation with wrong parameters (i.e. response, endpoint, and priority)
func TestAccCustomRuleResource_wrongParameters(t *testing.T) {
	mockClient := datadome.NewMockClient()

	testAccProvider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return mockClient, nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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
		},
	})
}
