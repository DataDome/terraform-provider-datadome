package datadome

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/datadome/terraform-provider/common"
	dd "github.com/datadome/terraform-provider/datadome-client-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceCustomRule define the CRUD operations and the schema definition for DataDome custom rules.
func resourceEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointCreate,
		ReadContext:   resourceEndpointRead,
		UpdateContext: resourceEndpointUpdate,
		DeleteContext: resourceEndpointDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"position_before": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Computed:     true,
			},
			"traffic_usage": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if !(value == "Account Creation" ||
						value == "Cart" ||
						value == "Form" ||
						value == "Forms" ||
						value == "General" ||
						value == "Login" ||
						value == "Payment" ||
						value == "Rss") {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   fmt.Sprintf("%q is not an acceptable traffic_usage", value),
						}
						diags = append(diags, diag)
					}
					return diags
				},
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if !(value == "Api" || value == "Mobile App" || value == "Web Browser") {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   fmt.Sprintf("%q is not an acceptable source", value),
						}
						diags = append(diags, diag)
					}
					return diags
				},
			},
			"cookie_same_site": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if !(value == "Lax" || value == "Strict" || value == "None") {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   fmt.Sprintf("%q is not an acceptable cookie_same_site", value),
						}
						diags = append(diags, diag)
					}
					return diags
				},
				Default: "Lax",
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion", "query"},
			},
			"path_inclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion", "query"},
			},
			"path_exclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion", "query"},
			},
			"user_agent_inclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion", "query"},
			},
			"query": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion", "query"},
			},
			"response_format": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if !(value == "json" || value == "html" || value == "auto") {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   fmt.Sprintf("%q is not an acceptable response_format", value),
						}
						diags = append(diags, diag)
					}
					return diags
				},
				Default: "auto",
			},
			"detection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		CustomizeDiff: customizeDiffEndpoints,
	}
}

// customizeDiffEndpoints applies additional verifications regarding the fields of the endpoint
// It raises an error when:
// - the "traffic_usage" value does not fit with the "source" value
// - the "protection_enabled" is set to `true` and the "detection_enabled" is set to `false`
// - the "query" field is not empty and one of "domain", "path_inclusion", "path_exclusion", or "user_agent_inclusion" is not empty either
func customizeDiffEndpoints(ctx context.Context, data *schema.ResourceDiff, meta interface{}) error {
	source := data.Get("source").(string)
	trafficUsage := data.Get("traffic_usage").(string)

	switch source {
	case "Api":
		expectedTrafficUsage := []string{"General"}
		if trafficUsage != "General" {
			return fmt.Errorf(`expected "traffic_usage" to be one of {%s}, got %q`, strings.Join(expectedTrafficUsage, ", "), trafficUsage)
		}
	case "Mobile App":
		expectedTrafficUsage := []string{"General", "Login", "Payment", "Cart", "Forms", "Account Creation"}
		if trafficUsage != "General" && trafficUsage != "Login" && trafficUsage != "Payment" && trafficUsage != "Cart" && trafficUsage != "Forms" && trafficUsage != "Account Creation" {
			return fmt.Errorf(`expected "traffic_usage" to be one of {%s}, got %q`, strings.Join(expectedTrafficUsage, ", "), trafficUsage)
		}
	}

	protectionEnabled := data.Get("protection_enabled").(bool)
	detectionEnabled := data.Get("detection_enabled").(bool)
	if !detectionEnabled && protectionEnabled {
		return fmt.Errorf("the detection must be activated in order to activate the protection")
	}

	_, domainExists := data.GetOk("domain")
	_, pathExclusionExists := data.GetOk("path_exclusion")
	_, pathInclusionExists := data.GetOk("path_inclusion")
	_, userAgentInclusionExists := data.GetOk("user_agent_inclusion")
	_, queryExists := data.GetOk("query")
	if queryExists && (domainExists || pathExclusionExists || pathInclusionExists || userAgentInclusionExists) {
		return fmt.Errorf(`"query" must be empty if whether "domain", "path_inclusion", "path_exclusion", or "user_agent_inclusion" is filled`)
	}

	return nil
}

// resourceCustomRuleCreate is used to create new custom rule
func resourceEndpointCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientEndpoint

	description := common.GetOptionalValueWithoutZeroValue[string](data, "description")
	positionBefore := common.GetOptionalValueWithoutZeroValue[string](data, "position_before")
	domain := common.GetOptionalValueWithoutZeroValue[string](data, "domain")
	pathInclusion := common.GetOptionalValueWithoutZeroValue[string](data, "path_inclusion")
	pathExclusion := common.GetOptionalValueWithoutZeroValue[string](data, "path_exclusion")
	userAgentInclusion := common.GetOptionalValueWithoutZeroValue[string](data, "user_agent_inclusion")
	query := common.GetOptionalValueWithoutZeroValue[string](data, "query")

	newEndpoint := dd.Endpoint{
		Name:               data.Get("name").(string),
		Description:        description,
		PositionBefore:     positionBefore,
		TrafficUsage:       data.Get("traffic_usage").(string),
		Source:             data.Get("source").(string),
		CookieSameSite:     data.Get("cookie_same_site").(string),
		Domain:             domain,
		PathInclusion:      pathInclusion,
		PathExclusion:      pathExclusion,
		UserAgentInclusion: userAgentInclusion,
		Query:              query,
		ResponseFormat:     data.Get("response_format").(string),
		DetectionEnabled:   data.Get("detection_enabled").(bool),
		ProtectionEnabled:  data.Get("protection_enabled").(bool),
	}

	id, err := c.Create(ctx, newEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(*id)

	return resourceEndpointRead(ctx, data, meta)
}

// resourceCustomRuleRead is used to fetch the custom rule by its ID
func resourceEndpointRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientEndpoint

	var diags diag.Diagnostics

	endpoint, err := c.Read(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = data.Set("name", endpoint.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("description", endpoint.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("position_before", endpoint.PositionBefore); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("traffic_usage", endpoint.TrafficUsage); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("source", endpoint.Source); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("cookie_same_site", endpoint.CookieSameSite); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("domain", endpoint.Domain); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("path_inclusion", endpoint.PathInclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("path_exclusion", endpoint.PathExclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("user_agent_inclusion", endpoint.UserAgentInclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("query", endpoint.Query); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("response_format", endpoint.ResponseFormat); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("detection_enabled", endpoint.DetectionEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("protection_enabled", endpoint.ProtectionEnabled); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceCustomRuleUpdate is used to update a custom rule by its ID
func resourceEndpointUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientEndpoint

	ID := data.Id()
	description := common.GetOptionalValueWithoutZeroValue[string](data, "description")
	positionBefore := common.GetOptionalValueWithoutZeroValue[string](data, "position_before")
	domain := common.GetOptionalValueWithoutZeroValue[string](data, "domain")
	pathInclusion := common.GetOptionalValueWithoutZeroValue[string](data, "path_inclusion")
	pathExclusion := common.GetOptionalValueWithoutZeroValue[string](data, "path_exclusion")
	userAgentInclusion := common.GetOptionalValueWithoutZeroValue[string](data, "user_agent_inclusion")
	query := common.GetOptionalValueWithoutZeroValue[string](data, "query")

	newEndpoint := dd.Endpoint{
		ID:                 &ID,
		Name:               data.Get("name").(string),
		Description:        description,
		PositionBefore:     positionBefore,
		TrafficUsage:       data.Get("traffic_usage").(string),
		Source:             data.Get("source").(string),
		CookieSameSite:     data.Get("cookie_same_site").(string),
		Domain:             domain,
		PathInclusion:      pathInclusion,
		PathExclusion:      pathExclusion,
		UserAgentInclusion: userAgentInclusion,
		Query:              query,
		ResponseFormat:     data.Get("response_format").(string),
		DetectionEnabled:   data.Get("detection_enabled").(bool),
		ProtectionEnabled:  data.Get("protection_enabled").(bool),
	}

	o, err := c.Update(ctx, newEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(*o.ID)

	return resourceEndpointRead(ctx, data, meta)
}

// resourceCustomRuleDelete is used to delete a custom rule by its ID
func resourceEndpointDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientEndpoint

	var diags diag.Diagnostics

	err := c.Delete(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
