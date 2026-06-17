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
					if value != "allow" && value != "captcha" && value != "block" && value != "device_check" && value != "intent_based" && value != "monetize" {
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
					if value != "high" && value != "normal" && value != "low" {
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
					validEnpointTypes := []string{"web", "account-creation", "login", "cart", "forms", "payment-web", "rss", "submit", "api-app-mobile", "account-creation-app-mobile", "api-app-mobile-login", "cart-app-mobile", "forms-app-mobile", "payment-app-mobile", "agentic-general", "agentic-account-creation", "agentic-login", "agentic-cart", "agentic-forms", "agentic-payment", "api"}
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
			"overridden_bot": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_box": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"policy_options.0.time_box", "policy_options.0.rate_limit"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorized_hours_of_the_week": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 167),
										},
									},
									"response_outside_time_box": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"block", "captcha", "device_check"}, false),
									},
								},
							},
						},
						"rate_limit": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"policy_options.0.time_box", "policy_options.0.rate_limit"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"applies_to": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"all_traffic", "ip", "session"}, false),
									},
									"threshold": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"time_frame": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"1m", "15m", "1h", "4h", "1d"}, false),
									},
									"response_after_threshold": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"block", "captcha", "device_check"}, false),
									},
								},
							},
						},
					},
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

// expandPolicyOptions converts the Terraform schema data for policy_options into a *dd.PolicyOptions.
// Returns nil when the block is absent or empty.
func expandPolicyOptions(data *schema.ResourceData) *dd.PolicyOptions {
	raw, ok := data.GetOk("policy_options")
	if !ok {
		return nil
	}
	list := raw.([]interface{})
	if len(list) == 0 {datadome/provider_test.go
		return nil
	}
	block := list[0].(map[string]interface{})

	po := &dd.PolicyOptions{}

	if v, ok := block["time_box"].([]interface{}); ok && len(v) > 0 {
		tb := v[0].(map[string]interface{})
		hoursRaw := tb["authorized_hours_of_the_week"].([]interface{})
		hours := make([]int, len(hoursRaw))
		for i, h := range hoursRaw {
			hours[i] = h.(int)
		}
		po.TimeBox = &dd.TimeBoxOptions{
			AuthorizedHoursOfTheWeek: hours,
			ResponseOutsideTimeBox:   tb["response_outside_time_box"].(string),
		}
	}

	if v, ok := block["rate_limit"].([]interface{}); ok && len(v) > 0 {
		rl := v[0].(map[string]interface{})
		po.RateLimit = &dd.RateLimitOptions{
			AppliesTo:              rl["applies_to"].(string),
			Threshold:              rl["threshold"].(int),
			TimeFrame:              rl["time_frame"].(string),
			ResponseAfterThreshold: rl["response_after_threshold"].(string),
		}
	}

	if po.TimeBox == nil && po.RateLimit == nil {
		return nil
	}
	return po
}

// flattenPolicyOptions converts a *dd.PolicyOptions into the list-of-maps representation
// expected by the Terraform schema. Returns nil when po is nil.
func flattenPolicyOptions(po *dd.PolicyOptions) []interface{} {
	if po == nil {
		return nil
	}

	block := map[string]interface{}{
		"time_box":   []interface{}{},
		"rate_limit": []interface{}{},
	}

	if po.TimeBox != nil {
		hours := make([]interface{}, len(po.TimeBox.AuthorizedHoursOfTheWeek))
		for i, h := range po.TimeBox.AuthorizedHoursOfTheWeek {
			hours[i] = h
		}
		block["time_box"] = []interface{}{
			map[string]interface{}{
				"authorized_hours_of_the_week": hours,
				"response_outside_time_box":    po.TimeBox.ResponseOutsideTimeBox,
			},
		}
	}

	if po.RateLimit != nil {
		block["rate_limit"] = []interface{}{
			map[string]interface{}{
				"applies_to":               po.RateLimit.AppliesTo,
				"threshold":                po.RateLimit.Threshold,
				"time_frame":               po.RateLimit.TimeFrame,
				"response_after_threshold": po.RateLimit.ResponseAfterThreshold,
			},
		}
	}

	return []interface{}{block}
}

// resourceCustomRuleCreate is used to create new custom rule
func resourceCustomRuleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*ProviderConfig)
	c := config.ClientCustomRule

	activatedAt := common.GetOptionalValue[string](data, "activated_at")
	enabled := common.GetOptionalValue[bool](data, "enabled")
	expiredAt := common.GetOptionalValue[string](data, "expired_at")

	newCustomRule := dd.CustomRule{
		Name:          data.Get("name").(string),
		Response:      data.Get("response").(string),
		Query:         data.Get("query").(string),
		EndpointType:  data.Get("endpoint_type").(string),
		Priority:      data.Get("priority").(string),
		Enabled:       enabled,
		ActivatedAt:   activatedAt,
		ExpiredAt:     expiredAt,
		PolicyOptions: expandPolicyOptions(data),
	}

	if v, ok := data.GetOk("overridden_bot"); ok {
		list := v.([]interface{})
		if len(list) > 0 {
			block := list[0].(map[string]interface{})
			newCustomRule.OverriddenBot = &dd.OverriddenBot{UUID: block["uuid"].(string)}
		}
	}

	id, err := c.Create(ctx, newCustomRule)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(*id))

	return resourceCustomRuleRead(ctx, data, meta)
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

	var overriddenBot []interface{}
	if customRule.OverriddenBot != nil {
		overriddenBot = []interface{}{
			map[string]interface{}{
				"uuid": customRule.OverriddenBot.UUID,
				"name": customRule.OverriddenBot.Name,
			},
		}
	}
	if err = data.Set("overridden_bot", overriddenBot); err != nil {
		return diag.FromErr(err)
	}

	if err = data.Set("policy_options", flattenPolicyOptions(customRule.PolicyOptions)); err != nil {
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
		ID:            &id,
		Name:          data.Get("name").(string),
		Response:      data.Get("response").(string),
		Query:         data.Get("query").(string),
		EndpointType:  data.Get("endpoint_type").(string),
		Priority:      data.Get("priority").(string),
		Enabled:       enabled,
		ActivatedAt:   activatedAt,
		ExpiredAt:     expiredAt,
		PolicyOptions: expandPolicyOptions(data),
	}

	if v, ok := data.GetOk("overridden_bot"); ok {
		list := v.([]interface{})
		if len(list) > 0 {
			block := list[0].(map[string]interface{})
			newCustomRule.OverriddenBot = &dd.OverriddenBot{UUID: block["uuid"].(string)}
		}
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
