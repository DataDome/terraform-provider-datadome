package datadome

import (
	"context"
	"fmt"
	"strconv"
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"traffic_usage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
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
							Detail:   fmt.Sprintf("%q is not an acceptable cookieSameSite", value),
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
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion"},
			},
			"path_inclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion"},
			},
			"path_exclusion": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"domain", "path_inclusion", "path_exclusion", "user_agent_inclusion"},
			},
			"user_agent_inclusion": {
				Type:         schema.TypeString,
				Optional:     true,
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
							Detail:   fmt.Sprintf("%q is not an acceptable responseFormat", value),
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
	}
}

// resourceCustomRuleCreate is used to create new custom rule
func resourceEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	var diags diag.Diagnostics

	description := d.Get("description").(string)
	positionBefore := d.Get("positionBefore").(string)
	domain := d.Get("domain").(string)
	pathInclusion := d.Get("pathInclusion").(string)
	pathExclusion := d.Get("pathExclusion").(string)
	userAgentInclusion := d.Get("userAgentInclusion").(string)

	newEndpoint := dd.Endpoint{
		Name:               d.Get("name").(string),
		Description:        &description,
		PositionBefore:     &positionBefore,
		TrafficUsage:       d.Get("trafficUsage").(string),
		Source:             d.Get("source").(string),
		CookieSameSite:     d.Get("cookieSameSite").(string),
		Domain:             &domain,
		PathInclusion:      &pathInclusion,
		PathExclusion:      &pathExclusion,
		UserAgentInclusion: &userAgentInclusion,
		ResponseFormat:     d.Get("responseFormat").(string),
		DetectionEnabled:   d.Get("detectionEnabled").(bool),
		ProtectionEnabled:  d.Get("protectionEnabled").(bool),
	}

	id, err := c.Create(ctx, newEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*id))

	return diags
}

// resourceCustomRuleRead is used to fetch the custom rule by its ID
func resourceEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint, err := c.Read(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", endpoint.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", endpoint.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("positionBefore", endpoint.PositionBefore); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("trafficUsage", endpoint.TrafficUsage); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("source", endpoint.Source); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("cookieSameSite", endpoint.CookieSameSite); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("domain", endpoint.Domain); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("pathInclusion", endpoint.PathInclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("pathExclusion", endpoint.PathExclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("userAgentInclusion", endpoint.UserAgentInclusion); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("responseFormat", endpoint.ResponseFormat); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("detectionEnabled", endpoint.DetectionEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("protectionEnabled", endpoint.ProtectionEnabled); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceCustomRuleUpdate is used to update a custom rule by its ID
func resourceEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	description := d.Get("description").(string)
	positionBefore := d.Get("positionBefore").(string)
	domain := d.Get("domain").(string)
	pathInclusion := d.Get("pathInclusion").(string)
	pathExclusion := d.Get("pathExclusion").(string)
	userAgentInclusion := d.Get("userAgentInclusion").(string)

	newEndpoint := dd.Endpoint{
		ID:                 id,
		Name:               d.Get("name").(string),
		Description:        &description,
		PositionBefore:     &positionBefore,
		TrafficUsage:       d.Get("trafficUsage").(string),
		Source:             d.Get("source").(string),
		CookieSameSite:     d.Get("cookieSameSite").(string),
		Domain:             &domain,
		PathInclusion:      &pathInclusion,
		PathExclusion:      &pathExclusion,
		UserAgentInclusion: &userAgentInclusion,
		ResponseFormat:     d.Get("responseFormat").(string),
		DetectionEnabled:   d.Get("detectionEnabled").(bool),
		ProtectionEnabled:  d.Get("protectionEnabled").(bool),
	}

	o, err := c.Update(ctx, newEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.ID))

	return resourceEndpointRead(ctx, d, m)
}

// resourceCustomRuleDelete is used to delete a custom rule by its ID
func resourceEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*ProviderConfig)
	c := config.ClientEndpoint

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.Delete(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
