---
page_title: "custom_rule Resource - terraform-provider-datadome"
subcategory: ""
description: |-
  The custom_rule resource allows you to configure a DataDome custom rule.
---

# Resource `datadome_custom_rule`

Creates a custom rule on DataDome dashboard

## Example Usage

### Usage with enabled field

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

### Usage with activation and expiration date

```terraform
resource "datadome_custom_rule" "new" {
  name          = "my-custom-rule"
  query         = "ip: 192.168.1.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "normal"
  activated_at  = "2030-01-31 23:59:59"
  expired_at    = "2050-01-31 23:59:59"
}
```

## Argument Reference

- `name` - (Required) Name of your custom rule. You cannot have multiple rules with the same name.
- `query` - (Required) Your query, for more information refer to the DataDome [documentation](https://docs.datadome.co/docs/syntax-guidelines)
- `response` - (Required) The response behavior, must be one of `allow`, `captcha`, `block`
- `endpoint_type` - (Optional) The endpoint on which you want your custom rule to be applied. If no endpoint type is specified, the custom rule will be applied to all endpoint types.
- `priority` - (Optional) Your rule priority, must be one of `high`, `low`, `normal`. Defaults to `high`.
- `enabled` - (Optional) Determines whether rule is enabled. If its value is set, it will override the value of `activated_at` and `expired_at` fields.
- `activated_at` - (Optional) Defines the date where the rule will be activated (Format Y-m-d H:i:s UTC+0).
- `expired_at` - (Optional) Defines the date where the rule will be deactivated (Format Y-m-d H:i:s UTC+0).

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.
