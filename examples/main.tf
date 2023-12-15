terraform {
  required_version = ">= 0.13.0"
  required_providers {
    datadome = {
      version = "0.0.1"
      source  = "datadome.co/app/datadome"
    }
  }
}

provider "datadome" {
  apikey = "apikey"
}

resource "datadome_custom_rule" "new" {
  name          = "test-terraform"
  query         = "ip: 192.168.0.1"
  response      = "whitelist"
  endpoint_type = "web"
  priority      = "normal"
}
