# Deployment code for terraform_provider_turing_pi_bmc

TF provider for Turing PI BMC

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.4.0 |
| <a name="requirement_tfe"></a> [tfe](#requirement\_tfe) | ~>0.40 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_tfe"></a> [tfe](#provider\_tfe) | ~>0.40 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [tfe_registry_module.terraform-module](https://registry.terraform.io/providers/hashicorp/tfe/latest/docs/resources/registry_module) | resource |
| [tfe_oauth_client.gh_oauth_token](https://registry.terraform.io/providers/hashicorp/tfe/latest/docs/data-sources/oauth_client) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_tfe_token"></a> [tfe\_token](#input\_tfe\_token) | TFE access for agents | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

## License

Apache-2.0

## Author Information

- Justin Lecher <justin@jlec.de>
