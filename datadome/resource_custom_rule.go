package datadome

import (
	"context"
	"fmt"
	"slices"
	"strconv"
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
				ForceNew: true,
			},
			"activated_at": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if value != "" {
						parsedTime, err := time.Parse(time.DateTime, value)
						if err != nil {
							diag := diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "invalid date format",
								Detail:   fmt.Sprintf("date '%s' does not match the required format 'YYYY-MM-DD HH:MM:SS'", value),
							}
							diags = append(diags, diag)
						} else {
							if parsedTime.Before(time.Now().UTC()) {
								diag := diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "invalid date",
									Detail:   fmt.Sprintf("date '%s' must not be in the past", value),
								}
								diags = append(diags, diag)
							}
						}
					}
					return diags
				},
			},
			"expired_at": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: func(v any, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					value := v.(string)
					if value != "" {
						parsedTime, err := time.Parse(time.DateTime, value)
						if err != nil {
							diag := diag.Diagnostic{
								Severity: diag.Error,
								Summary:  "invalid date format",
								Detail:   fmt.Sprintf("date '%s' does not match the required format 'YYYY-MM-DD HH:MM:SS'", value),
							}
							diags = append(diags, diag)
						} else {
							if parsedTime.Before(time.Now().UTC()) {
								diag := diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "invalid date",
									Detail:   fmt.Sprintf("date '%s' must not be in the past", value),
								}
								diags = append(diags, diag)
							}
						}
					}
					return diags
				},
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
		CustomizeDiff: customizeDiffCustomRules,
	}
}

// customizeDiffCustomRules applies additional verifications regarding the fields of the custom rule
// It raises error when:
// - expired_at is before activated_at
func customizeDiffCustomRules(ctx context.Context, data *schema.ResourceDiff, meta interface{}) error {
	activatedAt := data.Get("activated_at").(string)
	expiredAt := data.Get("expired_at").(string)

	if activatedAt != "" && expiredAt != "" {
		activatedTime, _ := time.Parse(time.DateTime, activatedAt)
		expiredTime, _ := time.Parse(time.DateTime, expiredAt)
		if activatedTime.After(expiredTime) {
			return fmt.Errorf("expired_at date must be after activated_at date")
		}
	}

	return nil
}

// resourceCustomRuleCreate is used to create new custom rule
func resourceCustomRuleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientCustomRule

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	activatedAt := common.GetOptionalValue[string](data, "activated_at")
	enabled := common.GetOptionalValue[bool](data, "enabled")
	expiredAt := common.GetOptionalValue[string](data, "expired_at")

	newCustomRule := dd.CustomRule{
		Name:         data.Get("name").(string),
		Response:     data.Get("response").(string),
		Query:        data.Get("query").(string),
		EndpointType: data.Get("endpoint_type").(string),
		Priority:     data.Get("priority").(string),
		Enabled:      enabled,
		ActivatedAt:  activatedAt,
		ExpiredAt:    expiredAt,
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
	if err = data.Set("activated_at", customRule.ActivatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err = data.Set("expired_at", customRule.ExpiredAt); err != nil {
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

	activatedAt := common.GetOptionalValue[string](data, "activated_at")
	enabled := common.GetOptionalValue[bool](data, "enabled")
	expiredAt := common.GetOptionalValue[string](data, "expired_at")

	newCustomRule := dd.CustomRule{
		ID:           &id,
		Name:         data.Get("name").(string),
		Response:     data.Get("response").(string),
		Query:        data.Get("query").(string),
		EndpointType: data.Get("endpoint_type").(string),
		Priority:     data.Get("priority").(string),
		Enabled:      enabled,
		ActivatedAt:  activatedAt,
		ExpiredAt:    expiredAt,
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
