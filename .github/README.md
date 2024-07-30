# DataDome Terraform Provider

This Terraform Provider aims at creating custom rules using the [Management API](https://docs.datadome.co/reference/get_1-1-protection-custom-rules).

## Build the provider

Run the following command to build the provider

```shell
$ make build
```

## Test a sample configuration manually

1. Build the provider.

2. Install the provider.

```shell
$ make install
```

3. Navigate to the `examples` directory. 

```shell
$ cd examples
```

4. Inside `main.tf`, set your Management API Key that you can find in [your dashboard](https://app.datadome.co/dashboard/management/integrations). If you don't have one, you can generate it.

5. Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

6. Terraform will ask you if you want to perform these actions: enter yes.

7. Congrats! You created a new custom rule that you can see in your dashboard.

## Make a release

1. On the main branch, create a tag with the version number, starting with `v`.

2. Push it.

3. GHA will release the provider on the [Terraform registry](https://registry.terraform.io/providers/DataDome/datadome/latest).