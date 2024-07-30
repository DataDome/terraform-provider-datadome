---
page_title: "Provider: DataDome"
subcategory: ""
description: |-
  Terraform provider for interacting with DataDome customer API.
---

# DataDome Provider

This provider can be used to create custom rules on your DataDome dashboard

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "datadome" {
  apikey = "exampleapikey123456"
}
```

## Schema

### Optional

- **apikey** (String, Optional) Management API key to authenticate to DataDome API. You can find it in [your dashboard](https://app.datadome.co/dashboard/management/integrations).
- **host** (String, Optional) Host of the DataDome custom rules API