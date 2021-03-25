terraform {
  required_providers {
    datadome = {
      version = "0.0.1"
      source = "datadome.co/app/datadome"
    }
  }
}

provider "datadome" {
  apikey = "apikey"
}

resource "datadome_custom_rule" "new" {
  name          = "TERRAFORMTEST"
  query         = "ip: 192.168.1.1a"
  response      = "whitelist"
  endpoint_type = "web"
  priority      = "normal"
}
