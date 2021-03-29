terraform {
  required_providers {
    datadome = {
      version = "1.0.0"
      source  = "datadome/datadome"
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
