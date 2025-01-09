---
page_title: "endpoint Resource - terraform-provider-datadome"
subcategory: ""
description: |-
  The endpoint resource allows you to configure a DataDome endpoint.
---

# Resource `datadome_endpoint`

Manage an endpoint on DataDome dashboard

## Example Usage

```terraform
resource "datadome_endpoint" "new" {
  name                 = "new_endpoint"
  description          = "An example of endpoint"
  position_before      = "Endpoint-ID"
  source               = "Web Browser"
  traffic_usage        = ""
  cookie_same_site     = "Lax"
  domain               = "example.org"
  path_inclusion       = ""
  path_exclusion       = ""
  user_agent_inclusion = "MYUSERAGENT"
  response_format      = "auto"
  detection_enabled    = true
  protection_enabled   = false
}
```

## Argument Reference

- `name` - (Required) The name of the endpoint resource.
- `description` - (Optional) The description of the endpoint resource.
- `position_before` - (Optional) The ID of the endpoint before which the new endpoint should be created. If this field is empty, it takes the ID of the default endpoint `WEB (default)`.
- `source` - (Required) Determine from where the traffic comes from. It only accepts `Api`, `Mobile App`, or `Web Browser`.
- `traffic_usage` - (Required) Determine for which purpose this endpoint is created. The value of this field depends on the `source` field:
  - For `Api`, it only accepts `General`.
  - For `Mobile App`, it accepts `General`, `Login`, `Payment`, `Cart`, `Forms`, or `Account Creation`.
  - For `Web Browser`, it accepts `Account Creation`, `Cart`, `Form`, `Forms`, `General`, `Login`, `Payment`, or `Rss`.
- `cookie_same_site` - (Optional) DataDome's cookie SameSite parameter for the endpoint. It only accepts `None`, `Lax`, or `Strict`. When not specified, it defaults to `Lax`.
- `domain` - (Optional) The domain for the endpoint when using the regex definition method.
- `path_inclusion` - (Optional) The path of inclusion for the endpoint when using the regex definition method.
- `path_exclusion` - (Optional) The path of exclusion for the endpoint when using the regex definition method.
- `user_agent_inclusion` - (Optional) The user agent inclusion for the endpoint when using the regex definition method.
- `response_format` - (Optional) The response format to use for challenged requests. It only accepts `auto`, `json`, or `html`. When not specified, it defaults to `auto`.
- `detection_enabled` - (Optional) Determine whether the detection is enabled. Defaults to `true`.
- `protection_enabled` - (Optional) Determing whether the protection is enabled. Defaults to `false`.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.
