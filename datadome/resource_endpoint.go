package datadome

import (
	"context"
	"fmt"
	"slices"
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
func customizeDiffSourceTrafficUsage(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	source := d.Get("source").(string)
	trafficUsage := d.Get("traffic_usage").(string)

	switch source {
	case "Api":
		expectedTrafficUsage := []string{"General"}
		if !slices.Contains(expectedTrafficUsage, trafficUsage) {
			return fmt.Errorf(`expected "traffic_usage" to be one of {%s}, got %q`, strings.Join(expectedTrafficUsage, ", "), trafficUsage)
		}
	case "Mobile App":
		expectedTrafficUsage := []string{"General", "Login", "Payment", "Cart", "Forms", "Account Creation"}
		if !slices.Contains(expectedTrafficUsage, trafficUsage) {
			return fmt.Errorf(`expected "traffic_usage" to be one of {%s}, got %q`, strings.Join(expectedTrafficUsage, ", "), trafficUsage)
		}
	}

	return nil
}

// resourceEndpointCreate is used to create new custom rule
func resourceEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	var description *string
	descriptionValue, ok := d.GetOk("description")
	if ok {
		descriptionString := descriptionValue.(string)
		description = &descriptionString
	}
	var positionBefore *string
	positionBeforeValue, ok := d.GetOk("position_before")
	if ok {
		positionBeforeString := positionBeforeValue.(string)
		positionBefore = &positionBeforeString
	}
	var domain *string
	domainValue, ok := d.GetOk("domain")
	if ok {
		domainString := domainValue.(string)
		domain = &domainString
	}
	var pathInclusion *string
	pathInclusionValue, ok := d.GetOk("path_inclusion")
	if ok {
		pathInclusionString := pathInclusionValue.(string)
		pathInclusion = &pathInclusionString
	}
	var pathExclusion *string
	pathExclusionValue, ok := d.GetOk("path_exclusion")
	if ok {
		pathExclusionString := pathExclusionValue.(string)
		pathExclusion = &pathExclusionString
	}
	var userAgentInclusion *string
	userAgentInclusionValue, ok := d.GetOk("user_agent_inclusion")
	if ok {
		userAgentInclusionString := userAgentInclusionValue.(string)
		userAgentInclusion = &userAgentInclusionString
	}

	newEndpoint := dd.Endpoint{
		Name:               d.Get("name").(string),
		Description:        description,
		PositionBefore:     positionBefore,
		TrafficUsage:       d.Get("traffic_usage").(string),
		Source:             d.Get("source").(string),
		CookieSameSite:     d.Get("cookie_same_site").(string),
		Domain:             domain,
		PathInclusion:      pathInclusion,
		PathExclusion:      pathExclusion,
		UserAgentInclusion: userAgentInclusion,
		ResponseFormat:     d.Get("response_format").(string),
		DetectionEnabled:   d.Get("detection_enabled").(bool),
		ProtectionEnabled:  d.Get("protection_enabled").(bool),
	}

	id, err := c.Create(ctx, newEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*id)

	return resourceEndpointRead(ctx, d, m)
}

// resourceEndpointRead is used to fetch the custom rule by its ID
func resourceEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	var diags diag.Diagnostics

	endpoint, err := c.Read(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", endpoint.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", endpoint.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("position_before", endpoint.PositionBefore); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("traffic_usage", endpoint.TrafficUsage); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("source", endpoint.Source); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("cookie_same_site", endpoint.CookieSameSite); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("domain", endpoint.Domain); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("path_inclusion", endpoint.PathInclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("path_exclusion", endpoint.PathExclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("user_agent_inclusion", endpoint.UserAgentInclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("response_format", endpoint.ResponseFormat); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("detection_enabled", endpoint.DetectionEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("protection_enabled", endpoint.ProtectionEnabled); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceEndpointUpdate is used to update a custom rule by its ID
func resourceEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	ID := d.Id()
	var description *string
	descriptionValue, ok := d.GetOk("description")
	if ok {
		descriptionString := descriptionValue.(string)
		description = &descriptionString
	}
	var positionBefore *string
	positionBeforeValue, ok := d.GetOk("position_before")
	if ok {
		positionBeforeString := positionBeforeValue.(string)
		positionBefore = &positionBeforeString
	}
	var domain *string
	domainValue, ok := d.GetOk("domain")
	if ok {
		domainString := domainValue.(string)
		domain = &domainString
	}
	var pathInclusion *string
	pathInclusionValue, ok := d.GetOk("path_inclusion")
	if ok {
		pathInclusionString := pathInclusionValue.(string)
		pathInclusion = &pathInclusionString
	}
	var pathExclusion *string
	pathExclusionValue, ok := d.GetOk("path_exclusion")
	if ok {
		pathExclusionString := pathExclusionValue.(string)
		pathExclusion = &pathExclusionString
	}
	var userAgentInclusion *string
	userAgentInclusionValue, ok := d.GetOk("user_agent_inclusion")
	if ok {
		userAgentInclusionString := userAgentInclusionValue.(string)
		userAgentInclusion = &userAgentInclusionString
	}

	newEndpoint := dd.Endpoint{
		ID:                 &ID,
		Name:               d.Get("name").(string),
		Description:        description,
		PositionBefore:     positionBefore,
		TrafficUsage:       d.Get("traffic_usage").(string),
		Source:             d.Get("source").(string),
		CookieSameSite:     d.Get("cookie_same_site").(string),
		Domain:             domain,
		PathInclusion:      pathInclusion,
		PathExclusion:      pathExclusion,
		UserAgentInclusion: userAgentInclusion,
		ResponseFormat:     d.Get("response_format").(string),
		DetectionEnabled:   d.Get("detection_enabled").(bool),
		ProtectionEnabled:  d.Get("protection_enabled").(bool),
	}

	o, err := c.Update(ctx, newEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*o.ID)

	return resourceEndpointRead(ctx, d, m)
}

// resourceEndpointDelete is used to delete a custom rule by its ID
func resourceEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	var diags diag.Diagnostics

	err := c.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
