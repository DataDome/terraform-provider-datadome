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

- **apikey** (String, Optional) API key to authenticate to DataDome API