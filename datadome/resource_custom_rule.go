package datadome

import (
	"context"
	"strconv"

	dd "github.com/datadome/terraform-provider/datadome-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRuleCreate,
		ReadContext:   resourceCustomRuleRead,
		UpdateContext: resourceCustomRuleUpdate,
		DeleteContext: resourceCustomRuleDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"query": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"response": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCustomRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dd.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	newCustomRule := dd.CustomRule{
		Name:         d.Get("name").(string),
		Response:     d.Get("response").(string),
		Query:        d.Get("query").(string),
		EndpointType: d.Get("endpoint_type").(string),
		Priority:     d.Get("priority").(string),
		Enabled:      d.Get("enabled").(bool),
	}

	o, err := c.CreateCustomRule(newCustomRule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(o.ID))

	return diags
}

func resourceCustomRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dd.Client)

	var diags diag.Diagnostics

	o, err := c.GetCustomRules()
	if err != nil {
		return diag.FromErr(err)
	}

	customRule := dd.CustomRule{}
	for _, v := range o {
		if strconv.Itoa(v.ID) == d.Id() {
			customRule = v
		}
	}

	d.Set("name", customRule.Name)
	d.Set("response", customRule.Response)
	d.Set("query", customRule.Query)
	d.Set("endpoint_type", customRule.EndpointType)
	d.Set("priority", customRule.Priority)
	d.Set("enabled", customRule.Enabled)

	return diags
}

func resourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dd.Client)

	id, err := strconv.Atoi(d.Id())

	newCustomRule := dd.CustomRule{
		ID:           id,
		Name:         d.Get("name").(string),
		Response:     d.Get("response").(string),
		Query:        d.Get("query").(string),
		EndpointType: d.Get("endpoint_type").(string),
		Priority:     d.Get("priority").(string),
		Enabled:      d.Get("enabled").(bool),
	}

	o, err := c.UpdateCustomRule(newCustomRule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.ID))
	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dd.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())

	customRuleToDelete := dd.CustomRule{
		ID: id,
	}

	_, err = c.DeleteCustomRule(customRuleToDelete)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
