package datadome

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	dd "github.com/datadome/terraform-provider/datadome-client-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceCustomRule define the CRUD operations and the schema definition for DataDome custom rules.
func resourceCustomRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRuleCreate,
		ReadContext:   resourceCustomRuleRead,
		UpdateContext: resourceCustomRuleUpdate,
		DeleteContext: resourceCustomRuleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					trimedValue := strings.TrimSpace(value)
					if trimedValue == "" {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   "the name value should not be blank",
						}
						diags = append(diags, diag)
					}
					return diags
				},
			},
			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"response": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if !(value == "allow" || value == "captcha" || value == "block") {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   fmt.Sprintf("%q is not an acceptable response", value),
						}
						diags = append(diags, diag)
					}
					return diags
				},
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if !(value == "high" || value == "normal" || value == "low") {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   fmt.Sprintf("%q is not an acceptable priority", value),
						}
						diags = append(diags, diag)
					}
					return diags
				},
				Default: "high",
			},
			"endpoint_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					validEnpointTypes := []string{"account-creation", "account-creation-app-mobile", "api", "api-app-mobile", "api-app-mobile-login", "cart", "cart-app-mobile", "forms", "forms-app-mobile", "login", "payment-app-mobile", "payment-web", "rss", "submit", "web"}
					var diags diag.Diagnostics
					value := v.(string)

					if !slices.Contains(validEnpointTypes, value) {
						diag := diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "wrong value",
							Detail:   fmt.Sprintf("%q is not an acceptable endpoint_type", value),
						}
						diags = append(diags, diag)
					}
					return diags
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
func resourceCustomRuleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientCustomRule

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	newCustomRule := dd.CustomRule{
		Name:         data.Get("name").(string),
		Response:     data.Get("response").(string),
		Query:        data.Get("query").(string),
		EndpointType: data.Get("endpoint_type").(string),
		Priority:     data.Get("priority").(string),
		Enabled:      data.Get("enabled").(bool),
	}

	id, err := c.Create(ctx, newCustomRule)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(*id))

	return diags
}

// resourceCustomRuleRead is used to fetch the custom rule by its ID
func resourceCustomRuleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientCustomRule

	var diags diag.Diagnostics

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	customRule, err := c.Read(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = data.Set("name", customRule.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("response", customRule.Response); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("query", customRule.Query); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("endpoint_type", customRule.EndpointType); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("priority", customRule.Priority); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("enabled", customRule.Enabled); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceCustomRuleUpdate is used to update a custom rule by its ID
func resourceCustomRuleUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientCustomRule

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	newCustomRule := dd.CustomRule{
		ID:           &id,
		Name:         data.Get("name").(string),
		Response:     data.Get("response").(string),
		Query:        data.Get("query").(string),
		EndpointType: data.Get("endpoint_type").(string),
		Priority:     data.Get("priority").(string),
		Enabled:      data.Get("enabled").(bool),
	}

	o, err := c.Update(ctx, newCustomRule)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.Itoa(*o.ID))
	return resourceCustomRuleRead(ctx, data, meta)
}

// resourceCustomRuleDelete is used to delete a custom rule by its ID
func resourceCustomRuleDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientCustomRule

	var diags diag.Diagnostics

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.Delete(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
