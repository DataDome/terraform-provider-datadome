package datadome

import (
	"context"

	"github.com/datadome/terraform-provider/common"
	"github.com/datadome/terraform-provider/datadome-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderConfig struct {
	ClientCustomRule common.API[datadome.CustomRule]
	ClientEndpoint   common.API[datadome.Endpoint]
}

// Provider of DataDome
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DATADOME_HOST", nil),
			},
			"apikey": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DATADOME_APIKEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"datadome_custom_rule": resourceCustomRule(),
			"datadome_endpoint":    resourceEndpoint(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure is used to configure the provider with the schema's variable
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apikey := d.Get("apikey").(string)

	var host *string

	hVal, ok := d.GetOk("host")
	if ok {
		tempHost := hVal.(string)
		host = &tempHost
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if apikey != "" {
		clientCustomRule, err := datadome.NewClientCustomRules(host, &apikey)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create DataDome client for custom rule management",
				Detail:   "Unable to authenticate user for authenticated DataDome client",
			})

			return nil, diags
		}

		clientEndpoint, err := datadome.NewClientEndpoint(host, &apikey)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create DataDome client for endpoint management",
				Detail:   "Unable to authenticate user for authenticated DataDome client",
			})

			return nil, diags
		}

		return &ProviderConfig{
			ClientCustomRule: clientCustomRule,
			ClientEndpoint:   clientEndpoint,
		}, diags
	}

	clientCustomRule, err := datadome.NewClientCustomRules(host, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create DataDome client for custom rule management",
			Detail:   "Unable to create anonymous DataDome client, please provide an api key",
		})
		return nil, diags
	}

	clientEndpoint, err := datadome.NewClientEndpoint(host, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create DataDome client for endpoint management",
			Detail:   "Unable to create anonymous DataDome client, please provide an api key",
		})

		return nil, diags
	}

	return &ProviderConfig{
		ClientCustomRule: clientCustomRule,
		ClientEndpoint:   clientEndpoint,
	}, diags
}
