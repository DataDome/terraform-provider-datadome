package datadome

import (
	"context"

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
				Required: true,
			},
			"endpoint_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
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
	}

	o, err := c.CreateCustomRule(newCustomRule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(o.Name)

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
		if v.Name == d.Id() {
			customRule = v
		}
	}

	d.Set("name", customRule.Name)
	d.Set("response", customRule.Response)
	d.Set("query", customRule.Query)
	d.Set("endpoint_type", customRule.EndpointType)
	d.Set("priority", customRule.Priority)

	return diags
}

func resourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dd.Client)

	newCustomRule := dd.CustomRule{
		Name:         d.Get("name").(string),
		Response:     d.Get("response").(string),
		Query:        d.Get("query").(string),
		EndpointType: d.Get("endpoint_type").(string),
		Priority:     d.Get("priority").(string),
	}

	o, err := c.UpdateCustomRule(newCustomRule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(o.Name)
	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dd.Client)

	var diags diag.Diagnostics

	customRuleToDelete := dd.CustomRule{
		Name: d.Get("name").(string),
	}

	_, err := c.DeleteCustomRule(customRuleToDelete)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
