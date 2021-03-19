---
page_title: "custom_rule Resource - terraform-provider-datadome"
subcategory: ""
description: |-
  The custom_rule resource allows you to configure a DataDome custom rule.
---

# Resource `datadome_custom_rule`

Creates a custom rule on DataDome dashboard

## Example Usage

```terraform
resource "datadome_custom_rule" "new" {
  name          = "my-custom-rule"
  query         = "ip: 192.168.1.1"
  response      = "whitelist"
  endpoint_type = "web"
  priority      = "normal"
}

```

## Argument Reference

- `name` - (Required) Name of your custom rule. The name is used as ID
- `query` - (Required) Your query refer to the DataDome [documentation]() !! PUT LINK
- `response` - (Required) The response behavior, must be one of `whitelist`, `captcha`, `block`
- `endpoint_type` - (Optionnal) The endpoint on which you want your custom rule to be applied. If no endpoint type is specified, the custom rule will be applied to all endpoint types.
- `priority` - (Optionnal) Your rule priority, must be one of `high`, `low`, `normal`


## Attributes Reference

In addition to all the arguments above, the following attributes are exported.
