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

### Usage with an overridden bot

```terraform
resource "datadome_custom_rule" "new" {
  name          = "my-custom-rule"
  query         = "ip: 192.168.1.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "normal"

  overridden_bot {
    uuid = "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Usage with a rate limit policy

```terraform
resource "datadome_custom_rule" "new" {
  name          = "my-custom-rule"
  query         = "ip: 192.168.1.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "normal"

  policy_options {
    rate_limit {
      applies_to               = "ip"
      threshold                = 100
      time_frame               = "15m"
      response_after_threshold = "block"
    }
  }
}
```

### Usage with a time box policy

```terraform
resource "datadome_custom_rule" "new" {
  name          = "my-custom-rule"
  query         = "ip: 192.168.1.1"
  response      = "allow"
  endpoint_type = "web"
  priority      = "normal"

  policy_options {
    time_box {
      authorized_hours_of_the_week = [9, 10, 11, 12, 13, 14, 15, 16, 17]
      response_outside_time_box    = "block"
    }
  }
}
```

## Argument Reference

- `name` - (Required) Name of your custom rule. You cannot have multiple rules with the same name.
- `query` - (Required) Your query, for more information refer to the DataDome [documentation](https://docs.datadome.co/docs/syntax-guidelines)
- `response` - (Required) The action applied to matching requests. Must be one of `allow`, `captcha`, `block`, `device_check`, `intent_based`, `monetize`. `device_check` triggers a device verification challenge. `intent_based` applies an intent-based evaluation. `monetize` triggers a monetization flow. `intent_based` and `monetize` are only valid when `overridden_bot` references an AI Agent. `policy_options` is only available for `allow` and `intent_based`.
- `endpoint_type` - (Optional) The endpoint on which you want your custom rule to be applied. If no endpoint type is specified, the custom rule will be applied to all endpoint types.
- `priority` - (Optional) Your rule priority, must be one of `high`, `low`, `normal`. Defaults to `high`.
- `enabled` - (Optional) Determines whether rule is enabled. If its value is set, it will override the value of `activated_at` and `expired_at` fields.
- `activated_at` - (Optional) Defines the date where the rule will be activated (Format Y-m-d H:i:s UTC+0).
- `expired_at` - (Optional) Defines the date where the rule will be deactivated (Format Y-m-d H:i:s UTC+0).
- `overridden_bot` - (Optional) The Verified Bot or AI Agent this rule applies to. Required when `response` is `intent_based` or `monetize`. When set, `policy_options.rate_limit.applies_to` must be `all_traffic`.
  - `uuid` - (Required) UUID of the Verified Bot or AI Agent.
  - `name` - (Computed) Name of the bot, populated from the API after creation.
- `policy_options` - (Optional) An optional policy block. Only one of `time_box` or `rate_limit` may be specified. Only available when `response` is `allow` or `intent_based`.
  - `rate_limit` - (Optional) Triggers an alternative response once a request threshold is exceeded within a time window. All sub-fields are required when this block is present.
    - `applies_to` - (Required) Scope over which the rate is counted. Must be one of `all_traffic`, `ip`, `session`. Use `all_traffic` when `overridden_bot` is set.
    - `threshold` - (Required) Maximum number of requests allowed within the defined time frame. Must be a positive integer.
    - `time_frame` - (Required) The time window over which the threshold is evaluated. Must be one of `1m`, `15m`, `1h`, `4h`, `1d`.
    - `response_after_threshold` - (Required) The action taken once the threshold is exceeded. Must be one of `block`, `captcha`, `device_check`.
  - `time_box` - (Optional) Restricts the rule to specific hours of the week, applying an alternative response outside the authorized window. All sub-fields are required when this block is present.
    - `authorized_hours_of_the_week` - (Required) List of authorized hour slots during the week. Each value represents an hour index from `0` (Monday 00:00) to `167` (Sunday 23:00).
    - `response_outside_time_box` - (Required) The action taken for requests received outside the authorized hours. Must be one of `block`, `captcha`, `device_check`.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.
