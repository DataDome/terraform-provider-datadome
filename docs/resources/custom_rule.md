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
  response      = "allow"
  endpoint_type = "web"
  priority      = "normal"
  enabled       = true
}

```

## Argument Reference

- `name` - (Required) Name of your custom rule. You cannot have multiple rules with the same name.
- `query` - (Required) Your query, for more information refer to the DataDome [documentation](https://docs.datadome.co/docs/syntax-guidelines)
- `response` - (Required) The response behavior, must be one of `allow`, `captcha`, `block`
- `endpoint_type` - (Optional) The endpoint on which you want your custom rule to be applied. If no endpoint type is specified, the custom rule will be applied to all endpoint types.
- `priority` - (Optional) Your rule priority, must be one of `high`, `low`, `normal`. Defaults to `high`.
- `enabled` - (Optional) Determines whether rule is enabled. Defaults to `true`.


## Attributes Reference

In addition to all the arguments above, the following attributes are exported.
