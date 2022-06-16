module github.com/datadome/terraform-provider

go 1.16

require (
	github.com/datadome/terraform-provider/datadome-client-go v0.0.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.0-rc.2
)

replace github.com/datadome/terraform-provider/datadome-client-go => ./datadome-client-go
