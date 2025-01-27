package datadome

import (
	"context"
	"fmt"
	"strings"
	"time"

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
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion"},
			},
			"path_inclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion"},
			},
			"path_exclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion"},
			},
			"user_agent_inclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion"},
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
		CustomizeDiff: customizeDiffSourceTrafficUsage,
	}
}

// customizeDiffSourceTrafficUsage raise an error in case the "traffic_usage" value does not fit with the "source" value
func customizeDiffSourceTrafficUsage(ctx context.Context, data *schema.ResourceDiff, meta interface{}) error {
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
		if trafficUsage != "General" && trafficUsage != "Login" && trafficUsage != "Payment" && trafficUsage != "Cart" && trafficUsage != "Forms" && trafficUsage != "Account" {
			return fmt.Errorf(`expected "traffic_usage" to be one of {%s}, got %q`, strings.Join(expectedTrafficUsage, ", "), trafficUsage)
		}
	}

	return nil
}

// resourceCustomRuleCreate is used to create new custom rule
func resourceEndpointCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientEndpoint

	var description *string
	descriptionValue, ok := data.GetOk("description")
	if ok {
		descriptionString := descriptionValue.(string)
		description = &descriptionString
	}
	var positionBefore *string
	positionBeforeValue, ok := data.GetOk("position_before")
	if ok {
		positionBeforeString := positionBeforeValue.(string)
		positionBefore = &positionBeforeString
	}
	var domain *string
	domainValue, ok := data.GetOk("domain")
	if ok {
		domainString := domainValue.(string)
		domain = &domainString
	}
	var pathInclusion *string
	pathInclusionValue, ok := data.GetOk("path_inclusion")
	if ok {
		pathInclusionString := pathInclusionValue.(string)
		pathInclusion = &pathInclusionString
	}
	var pathExclusion *string
	pathExclusionValue, ok := data.GetOk("path_exclusion")
	if ok {
		pathExclusionString := pathExclusionValue.(string)
		pathExclusion = &pathExclusionString
	}
	var userAgentInclusion *string
	userAgentInclusionValue, ok := data.GetOk("user_agent_inclusion")
	if ok {
		userAgentInclusionString := userAgentInclusionValue.(string)
		userAgentInclusion = &userAgentInclusionString
	}

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
	var description *string
	descriptionValue, ok := data.GetOk("description")
	if ok {
		descriptionString := descriptionValue.(string)
		description = &descriptionString
	}
	var positionBefore *string
	positionBeforeValue, ok := data.GetOk("position_before")
	if ok {
		positionBeforeString := positionBeforeValue.(string)
		positionBefore = &positionBeforeString
	}
	var domain *string
	domainValue, ok := data.GetOk("domain")
	if ok {
		domainString := domainValue.(string)
		domain = &domainString
	}
	var pathInclusion *string
	pathInclusionValue, ok := data.GetOk("path_inclusion")
	if ok {
		pathInclusionString := pathInclusionValue.(string)
		pathInclusion = &pathInclusionString
	}
	var pathExclusion *string
	pathExclusionValue, ok := data.GetOk("path_exclusion")
	if ok {
		pathExclusionString := pathExclusionValue.(string)
		pathExclusion = &pathExclusionString
	}
	var userAgentInclusion *string
	userAgentInclusionValue, ok := data.GetOk("user_agent_inclusion")
	if ok {
		userAgentInclusionString := userAgentInclusionValue.(string)
		userAgentInclusion = &userAgentInclusionString
	}

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
